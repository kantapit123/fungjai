{{ config(order_by='(timestamp)', engine='MergeTree()', materialized='table') }}

with stg_mood as (
SELECT *
FROM default.mood_responses
)

SELECT * from stg_mood