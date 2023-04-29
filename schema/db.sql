CREATE TABLE IF NOT EXISTS test-db.info_queue (
    name String,
    age int,
    created_at DateTime
) ENGINE = Kafka('kafka:9092', 'test-topic', 'test-group', 'JSONEachRow') settings kafka_thread_per_consumer = 0, kafka_num_consumers = 1;

CREATE TABLE  IF NOT EXISTS test-db.info (
    name String,
    age int,
    created_at DateTime
) ENGINE = MergeTree ORDER BY created_at;