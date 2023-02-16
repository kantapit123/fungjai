{{ config(materialized='view') }}

with moodtrack_avg as (
SELECT
    datetime,
    employee_id,
    avg(mood),
    max(mood),
    min(mood)
FROM default.moody_table
GROUP BY
    employee_id,
    datetime
ORDER BY datetime
)

SELECT * from moodtrack_avg