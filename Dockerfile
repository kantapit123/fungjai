FROM --platform=linux/amd64 apache/airflow:2.5.0-python3.10

LABEL team="SUT"

USER root
RUN apt -y update && apt install -y p7zip-full

USER airflow
RUN pip install \
    airflow-dbt==0.4.0 \
    airflow-dbt-python==1.0.5 \
    # dbt-core==1.4.6 \
    dbt-core==1.3.1 \
    minio==7.1.13 \
    # dbt-clickhouse==1.4.0 \
    dbt-clickhouse==1.3.3 \
    clickhouse-connect==0.6.2