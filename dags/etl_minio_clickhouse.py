from datetime import datetime, timedelta
from airflow.models import DAG, Variable
from airflow.decorators import task
from io import BytesIO
from airflow.operators.empty import EmptyOperator
from minio import Minio
import pandas as pd
import clickhouse_connect

DOC_MD = """

### Prerequisites

#### Variables
    Import this variables.json into airflow
    {
        "clickhouse_host_name": "fungjai-clickhouse-dev.dbiteam.com",
        "clickhouse_interface": "https",
        "clickhouse_password": "🔑",
        "clickhouse_port": 443,
        "clickhouse_username": "admin",
        "minio_access_key": "🔑",
        "minio_bucket_name": "fungjai",
        "minio_host_name": "fungjai-minio-dev.dbiteam.com",
        "minio_prefix": "responses/mood/",
        "minio_secret_key": "🔑"
    }
"""

default_args = {
    "owner": "airflow",
    "depends_on_past": False,
    "start_date": datetime(2023, 2, 15),
    "email_on_failure": False,
    "email_on_retry": False,
    "retries": 1,
    "retry_delay": timedelta(minutes=5),
}

with DAG(
    dag_id="A_etl_minio_clickhouse",
    default_args=default_args,
    schedule="@daily",
    catchup=False,
    tags=["A_upload_file_to_minio"],
    description="DAG to upload a file to MinIO",
    max_active_runs=1,
):
    start = EmptyOperator(task_id="start")

    @task
    def etl_minio_to_clickhouse():

        minio_host_name = Variable.get("minio_host_name")
        minio_access_key = Variable.get("minio_access_key")
        minio_secret_key = Variable.get("minio_secret_key")
        minio_bucket_name = Variable.get("minio_bucket_name")
        minio_prefix = Variable.get("minio_prefix")

        # Connect to minio
        minio_client = Minio(
            minio_host_name,
            access_key=minio_access_key,
            secret_key=minio_secret_key,
            secure=False,
        )
        # key = "responses"
        objects = minio_client.list_objects(minio_bucket_name, prefix=minio_prefix, recursive=True)
        df = pd.DataFrame()
        for obj in objects:
            print(obj.object_name)
            if obj.object_name.endswith(".json"):  # Ensure the object has the .json extension
                # Download the file from MinIO
                data = minio_client.get_object(minio_bucket_name, obj.object_name)
                # Read the file content into a BytesIO object
                file_data = BytesIO(data.read())
                # Convert the data to a pandas DataFrame
                file_df = pd.read_json(file_data)
                # Append the file_df to the main DataFrame
                df = df.append(file_df, ignore_index=True)

        print(df.to_string())
        print(df.dtypes)
        clickhouse_interface = Variable.get("clickhouse_interface")
        clickhouse_host = Variable.get("clickhouse_host_name")
        clickhouse_port = Variable.get("clickhouse_port")
        clickhouse_username = Variable.get("clickhouse_username")
        clickhouse_password = Variable.get("clickhouse_password")

        clickhouse_client = clickhouse_connect.get_client(
            interface=clickhouse_interface,
            host=clickhouse_host,
            port=clickhouse_port,
            username=clickhouse_username,
            password=clickhouse_password,
        )

        clickhouse_client.command(
            f"""INSERT INTO mood_responses SELECT * FROM s3(
                '{clickhouse_interface}://{minio_host_name}/{minio_bucket_name}/{minio_prefix}*/*.json',
                '{minio_access_key}',
                '{minio_secret_key}',
                'JSONEachRow'
                )
            """
        )

    end = EmptyOperator(task_id="end")

    etl_minio_to_clickhouse_task = etl_minio_to_clickhouse()

    start >> etl_minio_to_clickhouse_task >> end
