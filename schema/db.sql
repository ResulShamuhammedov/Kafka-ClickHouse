CREATE TABLE default.info_queue (
    name String,
    age int
) ENGINE = Kafka('kafka:9092', 'test-topic', 'test-group', 'JSONEachRow') SETTINGS kafka_thread_per_consumer = 0, kafka_num_consumers = 1;

CREATE TABLE default.info (
    name String,
    age int
) ENGINE = MergeTree ORDER BY name;

-- SET stream_poll_timeout_ms=20000;

CREATE MATERIALIZED VIEW IF NOT EXISTS default.mv_info 
TO default.info AS SELECT * FROM default.info_queue;