from datetime import datetime, timedelta
from airflow.models import DAG
from airflow.decorators import task
from io import BytesIO
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

        # Prod
        # host = "fungjai-minio.dbiteam.com"
        # access_key = "rootuser"
        # secret_key = "rootpass123"
        # bucket_name = "fungjai"

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

        # # List objects information.

        # prefix = "responses/mood/"
        # file_save_dir = "path/to/save/the/objects/"

        # # Get the list of objects that match the pattern
        # objects = minio_client.list_objects(bucket_name
        # , prefix=prefix
        # , recursive=True)
        # # Iterate over the objects and create a DataFrame from each file
        # data_frames = []
        # # Iterate over the objects and download them
        # for obj in objects:
        #     object_name = obj.object_name
        #     if object_name.endswith('.json'):
        #         parts = object_name.split('/')
        #         folder_name = '/'.join(parts[:-1])
        #         # Extract the file name from the object's full path
        #         file_name = parts[-1]
        #         print("Folder:", folder_name)
        #         print("File:", file_name)
        #         file_path = f"{bucket_name}/{folder_name}/{file_name}"
        #         # Read the JSON file into a DataFrame
        #         data_frame = pd.read_json(obj)
        #         data_frames.append(data_frame)

        # # Concatenate all DataFrames into a single DataFrame
        # combined_data_frame = pd.concat(data_frames, ignore_index=True)
        # print(combined_data_frame)

        # # Get data from minio
        # obj = minio_client.get_object(
        #     bucket,
        #     key,
        # )

        # df = pd.read_csv(obj, index_col=False)
        # print("data", df)

        # Push to clickhouse
        import clickhouse_connect

        client = clickhouse_connect.get_client(
            host="click_server", port="8123", username="default"
        )

        # data = df.values.tolist()
        # table_name = 'test3'
        # client.insert("test2", data, column_names=["type"
        # , "postback"
        # , "webhookEventId"
        # , "deliveryContext"
        # , "timestamp"
        # , "source"
        # , "replyToken"
        # , "mode"])
        # client.insert_df(table=table_name,df=df.astype(str),database='default')
        # client.insert_df("YourTable", df)
        # DEV
        client.command(
            """INSERT INTO mood_responses
             SELECT *
              FROM s3('http://minio:9000/fungjai/responses/mood/*/*.json'
              , 'minio123'
              , 'minio123mak'
              , 'JSONEachRow')
              """
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
