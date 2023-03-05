import io
import time

from confluent_kafka import Consumer
from datetime import datetime
from minio import Minio


def upload_to_minio(text):
    # minio config
    ACCESS_KEY = "minio123"
    SECRET_KEY = "minio123"
    MINIO_API_HOST = "localhost:9000"
    BUCKET_NAME = "test-buckets"
    MINIO_CLIENT = Minio(
        MINIO_API_HOST,
        access_key=ACCESS_KEY,
        secret_key=SECRET_KEY,
        secure=False
    )

    # time
    str_datetime = datetime.now().strftime("%d-%m-%Y")
    timestamp = int(round(time.time() * 1000))

    # file name
    FILE_NAME = f"mood_{str_datetime}_{timestamp}.json"

    # convert to bytes
    text_bytes = io.BytesIO(text)

    # upload to minio
    try:
        MINIO_CLIENT.put_object(BUCKET_NAME, FILE_NAME, text_bytes, len(text))
        print("Success upload to bucket")
    except:
        print("Error upload to bucket")


# kafka config
c = Consumer(
    {
        "bootstrap.servers": "localhost:9092",
        "group.id": "my-consumer-group",
        "auto.offset.reset": "earliest",
    }
)
c.subscribe(["topic-kafka-json"])

# pull data from kafka
count = 0
while count < 3:
    message = c.poll(1.0)  # pull messages from kafka every 1 second
    if message is None:
        print("message is None")
        count += 1
    else:
        text_encode = message.value()
        text = text_encode.decode("utf-8")
        print(f"Received message: {text}")
        upload_to_minio(text_encode)
