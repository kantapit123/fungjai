"""
### open api pipeline dbt Processing

"""
from datetime import timedelta

from airflow import DAG
from airflow.operators.empty import EmptyOperator
from airflow.utils import timezone
from airflow_dbt_python.operators.dbt import DbtRunOperator, DbtTestOperator
from airflow.operators.bash import BashOperator
from minio import Minio


# import slack_notification as slack

DOC_MD = """
## 💰 open api pipeline dbt
can make dbt into schedule

### Prerequisites

#### Connections
1. **s3_airflow_conn** [conn_type=`Amazon Web Services`, aws_ak, aws_sk]

#### Variables
2. **openapi_dbt_target** [type=`text`, ex=`ci-dev`]

"""

BUSINESS_DOMAIN = "clickhouse_dbt"
S3_BUCKET = "fungjai"
S3_PROJECT_DIR = f"s3://{S3_BUCKET}/dbt/{BUSINESS_DOMAIN}/dbt-project.zip"
S3_PROFILE_DIR = f"s3://{S3_BUCKET}/dbt/{BUSINESS_DOMAIN}/profiles.yml"
S3_CONN = "local_minio"

default_args = {
    "owner": "Data Analytics",
    "start_date": timezone.datetime(2023, 1, 23),
    # "on_failure_callback": slack.notify,
    "retries": 3,
    "retry_delay": timedelta(minutes=3),
    "pool": "fungjai_pool",
}

with DAG(
    "fungjai_dbt_processing",
    default_args=default_args,
    schedule="@daily",
    catchup=False,
    tags=["fungjai"],
    max_active_runs=1,
) as dag:
    start = EmptyOperator(task_id="start")

    dbt_run_staging = BashOperator(
        task_id="dbt_run_staging",
        bash_command="""
        cd /opt/airflow/dbt ;
        pwd ;
        dbt run
    
    """,
    )

    dbt_test = BashOperator(
        task_id="dbt_test",
        bash_command="""
        cd /opt/airflow/dbt ;
        pwd ;
        dbt test
    
    """,
    )

    # dbt_run_staging = DbtRunOperator(
    #     task_id="dbt_run_fungjai_pipeline",
    #     project_dir=S3_PROJECT_DIR,
    #     profiles_dir=S3_PROFILE_DIR,
    #     project_conn_id=S3_CONN,
    #     profiles_conn_id=S3_CONN,
    #     target="{{ var.value.fungjai_dbt_target }}",
    #     vars={"uploaded_datetime": "{{ ds }}"},
    #     fail_fast=True,
    # )

    # dbt_test = DbtTestOperator(
    #     task_id="dbt_test",
    #     project_dir=S3_PROJECT_DIR,
    #     profiles_dir=S3_PROFILE_DIR,
    #     project_conn_id=S3_CONN,
    #     profiles_conn_id=S3_CONN,
    #     target="{{ var.value.fungjai_dbt_target }}",
    # )

    end = EmptyOperator(task_id="end")

    start >> dbt_run_staging >> dbt_test >> end
