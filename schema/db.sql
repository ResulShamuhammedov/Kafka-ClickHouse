CREATE TABLE default.info_queue (
    name String,
    age int
) ENGINE = Kafka('kafka:9092', 'test-topic', 'test-group', 'JSONEachRow') settings kafka_thread_per_consumer = 0, kafka_num_consumers = 1;

CREATE TABLE default.info (
    name String,
    age int
) ENGINE = MergeTree ORDER BY name;