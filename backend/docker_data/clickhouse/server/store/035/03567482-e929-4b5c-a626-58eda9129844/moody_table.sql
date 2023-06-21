ATTACH TABLE _ UUID 'e9756b15-8294-450b-86a9-a388ae9b4e24'
(
    `employee_id` String,
    `datetime` DateTime,
    `mood` UInt32
)
ENGINE = S3('http://minio:9000/moodybucket/moodtracker/*/*.csv', 'minio', 'minio123', 'CSV')
