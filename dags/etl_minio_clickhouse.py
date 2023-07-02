from datetime import datetime, timedelta
from airflow.models import DAG
from airflow.decorators import task
from io import BytesIO
from airflow.operators.empty import EmptyOperator
from minio import Minio
import pandas as pd
import clickhouse_connect


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
        # Config
        # host = Variable.get("host")
        # access_key = Variable.get("access_key")
        # secret_key = Variable.get("secret_key")
        # bucket = Variable.get("bucket")

        # DEV
        host = "minio:9000"
        access_key = "minio123"
        secret_key = "minio123mak"
        bucket_name = "fungjai"
        prefix = "responses/mood/"

        # Connect to minio
        minio_client = Minio(
            host, access_key=access_key, secret_key=secret_key, secure=False
        )
        # key = "responses"
        objects = minio_client.list_objects(bucket_name, prefix=prefix, recursive=True)
        df = pd.DataFrame()
        for obj in objects:
            print(obj.object_name)
            if obj.object_name.endswith(
                ".json"
            ):  # Ensure the object has the .json extension
                # Download the file from MinIO
                data = minio_client.get_object(bucket_name, obj.object_name)
                # Read the file content into a BytesIO object
                file_data = BytesIO(data.read())
                # Convert the data to a pandas DataFrame
                file_df = pd.read_json(file_data)
                # Append the file_df to the main DataFrame
                df = df.append(file_df, ignore_index=True)

        print(df.to_string())
        print(df.dtypes)

        client = clickhouse_connect.get_client(
            host="click_server", port="8123", username="default"
        )
        client.command(
            "INSERT INTO mood_responses SELECT * FROM s3('http://minio:9000/fungjai/responses/mood/*/*.json', 'minio123', 'minio123mak', 'JSONEachRow')"
        )
        # Prod
        # client.command("INSERT INTO mood_responses SELECT * FROM s3(
        # 'https://fungjai-minio.dbiteam.com/fungjai/responses/mood/*/*.json',
        # 'rootuser', 'rootpass123', '
        # JSONEachRow')")

        # print(query_result)

    end = EmptyOperator(task_id="end")

    etl_minio_to_clickhouse_task = etl_minio_to_clickhouse()

    start >> etl_minio_to_clickhouse_task >> end
