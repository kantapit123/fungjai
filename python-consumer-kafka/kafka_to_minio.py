import io
import time

from confluent_kafka import Consumer
from datetime import datetime
from minio import Minio
import json
from dotenv import load_dotenv
import os

load_dotenv()


def upload_to_minio(data_json, timestamp, hash_userId):
    # minio config
    MINIO_ACCESS_KEY_ID = os.getenv("MINIO_ACCESS_KEY_ID")  # "minio123"
    MINIO_SECRET_ACCESS_KEY = os.getenv("MINIO_SECRET_ACCESS_KEY")  # "minio123mak"
    MINIO_ENDPOINT = os.getenv("MINIO_ENDPOINT")  # "localhost:9000"
    MINIO_BUCKET_NAME = os.getenv("MINIO_BUCKET_NAME")  # "test-buckets"
    MINIO_CLIENT = Minio(
        MINIO_ENDPOINT,
        access_key=MINIO_ACCESS_KEY_ID,
        secret_key=MINIO_SECRET_ACCESS_KEY,
        secure=False,
    )

    # time
    str_datetime = datetime.utcnow().strftime("%Y-%m-%d")
    timestamp = int(round(time.time() * 1000))

    # file name
    FILE_NAME = f"responses/mood/{str_datetime}/{timestamp}-{hash_userId}.json"

    # convert to bytes
    data_bytes = io.BytesIO(data_json)

    # upload to minio
    try:
        MINIO_CLIENT.put_object(MINIO_BUCKET_NAME, FILE_NAME, data_bytes, len(text))
        print("Success upload to bucket")
    except:
        print("Error upload to bucket")


# kafka config
BOOTSTRAP_SERVERS = os.getenv("BOOTSTRAP_SERVERS")
KAFKA_GROUP = os.getenv("KAFKA_GROUP")

consumer = Consumer(
    {
        "bootstrap.servers": "kafka:9092",
        "group.id": "my-group",
        "auto.offset.reset": "earliest",
    }
)
consumer.subscribe(["mood_responses"])

# pull data from kafka
while True:
    message = consumer.poll(1.0)  # pull messages from kafka every 1 second
    if message is None:
        pass
    else:
        text_encode = message.value()
        text = text_encode.decode("utf-8")
        print(f"Received message: {text}")
        str_2_json = json.loads(text)
        timestamp = str_2_json["timestamp"]
        user_hash_id = str_2_json["hash_userId"]
        print(timestamp)
        print(user_hash_id)
        upload_to_minio(text_encode, timestamp, user_hash_id)

c.close()
