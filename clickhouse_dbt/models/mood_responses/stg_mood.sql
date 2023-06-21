{{ config(materialized='view') }}

with stg_mood_responses as (
SELECT timestamp, source.userId as uid, postback.data as response, multiIf(response = 'question=1&value=1', 1, response = 'question=1&value=2', 2, response = 'question=1&value=3', 3, response = 'question=1&value=4', 4, null) as mood
FROM clickhouse_dbt.source_mood
)

SELECT * from stg_mood_responses