{{ config(materialized='view') }}

with moodtrack_employee as (
SELECT
    employee_id,
    avg(mood),
    max(mood),
    min(mood)
FROM default.moody_table
GROUP BY
    employee_id
ORDER BY employee_id
)

SELECT * from moodtrack_employee