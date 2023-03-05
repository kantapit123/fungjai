from datetime import datetime, timedelta
from airflow.models import DAG
from airflow.decorators import task
from airflow.operators.empty import EmptyOperator
from minio import Minio
import pandas as pd

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
    schedule="@once",
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
        host = "minio:9000"
        access_key = "minio123"
        secret_key = "minio123"
        bucket = "test-buckets"

        # Connect to minio
        minioClient = Minio(
            host, access_key=access_key, secret_key=secret_key, secure=False
        )
        key = "data.csv"

        # Get data from minio
        obj = minioClient.get_object(
            bucket,
            key,
        )
        df = pd.read_csv(obj, index_col=False)
        print("data", df)

        # Push to clickhouse
        import clickhouse_connect

        client = clickhouse_connect.get_client(
            host="click_server", port=8123, username="default"
        )
        client.insert_df("test1", df)
        query_result = client.query("SELECT * FROM test1")
        print(query_result.result_set)

    end = EmptyOperator(task_id="end")

    etl_minio_to_clickhouse_task = etl_minio_to_clickhouse()

    start >> etl_minio_to_clickhouse_task >> end
