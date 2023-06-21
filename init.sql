CREATE TABLE IF NOT EXISTS mood_responses (
    type String,
    postback Tuple (
        data String
    ),
    webhookEventId String,
    deliveryContext Tuple (
        isRedelivery UInt8
    ),
    timestamp DateTime64(9),
    source Tuple (
        type String,
        userId String
    ),
    replyToken String,
    mode String
) ENGINE = MergeTree()
-- ENGINE = S3(
--         'http://minio:9000/fungjai/responses/mood/*/*.json',
--         'minio123',
--         'minio123mak',
--         'JSONEachRow')
ORDER BY timestamp;