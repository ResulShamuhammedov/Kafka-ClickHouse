# Kafka-ClickHouse
Simple repo to test integration of Kafka and ClickHouse

### Documentation

1. Clone the repo:

```bash
git clone git@github.com:ResulShamuhammedov/Kafka-ClickHouse.git
```

1. Change directory to Kafka-Clickhouse and run `docker compose up` command
2. Then create new Topic in Kafka container. To do that open the terminal and run:
    
    ```bash
    docker exec -it kafka /opt/bitnami/kafka/bin/kafka-topics.sh --create --bootstrap-server kafka:9092 --replication-factor 1 --partitions 1 --topic test-topic
    ```
    
3. Enter the Clickhouse container in interactive mode to create required tables and monitor what happens:

```bash
docker exec -it clickhouse bash

cd bin

clickhouse client
```

if after that in terminal you see :) it means you entered clickhouse shell.

1. Then run these commands to create tables
    
    ```bash
    CREATE TABLE default.info_queue (
        name String,
        age int
    ) ENGINE = Kafka('kafka:9092', 'test-topic', 'test-group', 'JSONEachRow') settings kafka_thread_per_consumer = 0, kafka_num_consumers = 1;
    
    CREATE TABLE default.info (
        name String,
        age int
    ) ENGINE = MergeTree ORDER BY name;
    ```
    
2. Now, there is service running on port 8080
API endpoint: `POST : http://localhost:8080/metric`
    
    Body:
    
    ```json
    {
        "name": "Resul",
        "age": 21
    }
    ```
    
    ---
    
    when you send this request and it returns OK, it means backend has successfully published message to Kafka, and Clickhouses’s Kafka engine table read it.
    
    If you want to read from `info_queue` table, you should first:
    
    ```sql
    SET stream_like_engine_allow_direct_select=1;
    ```