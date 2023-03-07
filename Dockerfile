FROM --platform=linux/amd64 apache/airflow:2.5.0-python3.10

LABEL team="SUT"

USER root
RUN apt -y update && apt install -y p7zip-full

USER airflow
RUN pip install \
    minio==7.1.13 \
    clickhouse-connect==0.5.11
