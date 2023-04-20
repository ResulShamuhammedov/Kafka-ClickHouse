CREATE TABLE default.info_queue (
    name String,
    age int
) ENGINE = Kafka('localhost:9092', 'test-topic', 'test-group', 'JSONEachRow');

CREATE TABLE default.info (
    name String,
    age int
) ENGINE = MergeTree ORDER BY name;