package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/ResulShamuhammedov/Kafka-Clickhouse/server/handler"
	"github.com/ResulShamuhammedov/Kafka-Clickhouse/server/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/robfig/cron"
	"github.com/segmentio/kafka-go"

	"github.com/joho/godotenv"
)

const (
	createQueueTable = `CREATE TABLE IF NOT EXISTS default.info_queue (
		name String,
		age int,
		created_at DateTime
	) ENGINE = Kafka('kafka:9092', 'test-topic', 'test-group', 'JSONEachRow') settings kafka_thread_per_consumer = 0, kafka_num_consumers = 1;`

	createMergeTreeTable = `CREATE TABLE default.info (
		name String,
		age int,
		created_at DateTime
	) ENGINE = MergeTree ORDER BY created_at;`
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	batchSize, err := strconv.Atoi(os.Getenv("BATCH_SIZE"))
	if err != nil {
		log.Fatal(err)
	}

	conn, err := repository.NewClickhouseDB(repository.ClickhouseConfig{
		Host:     os.Getenv("CH_HOST"),
		Port:     os.Getenv("CH_PORT"),
		Username: os.Getenv("CH_USERNAME"),
		DBName:   os.Getenv("CH_NAME"),
		Password: os.Getenv("CH_PASSWORD"),
	})
	if err != nil {
		log.Fatal(err)
	}

	db := repository.NewDB(conn)

	/*
		// To create Kafka Engine Table there must be Kafka broker
		err = db.Conn.Exec(context.Background(), createQueueTable)
		if err != nil {
			log.Fatal(err)
		}
	*/
	// err = db.Conn.Exec(context.Background(), createMergeTreeTable)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"kafka:9092"},
		Topic:   "test-topic",
	})
	defer writer.Close()

	handler := handler.NewHandlerV1(writer)

	app := fiber.New()

	app.Post("/metric", handler.HandleGetMetrics)

	// Set up cron job
	cronJob := cron.New()
	cronJob.AddFunc("@every 5m", func() {
		err := db.BatchInsert(context.Background(), batchSize)
		if err != nil {
			log.Fatal(err)
		}
	})
	cronJob.Start()

	err = app.Listen(":8080")
	if err != nil {
		log.Fatalf("could not listen: %s", err.Error())
	}
}
