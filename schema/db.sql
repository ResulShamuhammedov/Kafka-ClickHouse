CREATE TABLE IF NOT EXISTS default.info_queue (
    name String,
    age int,
    created_at DateTime
) ENGINE = Kafka('kafka:9092', 'test-topic', 'test-group', 'JSONEachRow') settings kafka_thread_per_consumer = 0, kafka_num_consumers = 1, kafka_max_block_size = 10;

CREATE TABLE  IF NOT EXISTS default.info (
    name String,
    age int,
    created_at DateTime
) ENGINE = MergeTree ORDER BY created_at;

-- Create materialized view
CREATE MATERIALIZED VIEW IF NOT EXISTS default.mv_queue_info
TO default.info
AS
SELECT *
FROM (
  SELECT *, count() OVER () as cnt
  FROM default.info_queue
) 
WHERE cnt >= 10;

CREATE MATERIALIZED VIEW IF NOT EXISTS default.mv_info 
TO default.info AS SELECT * FROM default.info_queue;