{{ config(order_by='(timestamp)', engine='MergeTree()', materialized='incremental', unique_key='timestamp') }}

with mood as (
SELECT *
FROM default.mood
)

SELECT * from mood

{% if is_incremental() %}

-- this filter will only be applied on an incremental run
where timestamp = timestamp

{% endif %}