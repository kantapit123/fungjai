CREATE TABLE moody_table (
    employee_id String,
    datetime datetime,
    mood UInt32)
    ENGINE = S3(
        'http://minio:9000/moodybucket/moodtracker/*/*.csv',
        'minio',
        'minio123',
        'CSV')