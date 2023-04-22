package main

import (
	"log"
	"time"

	"github.com/ResulShamuhammedov/Kafka-Clickhouse/server/handler"
	"github.com/Shopify/sarama"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 500 * time.Millisecond
	config.Producer.Return.Successes = true

	brokerList := []string{"localhost:9092"}

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			panic(err)
		}
	}()

	handler := handler.NewHandlerV2(producer)

	app := fiber.New()

	app.Post("/metric", handler.HandleGetMetrics)

	err = app.Listen(":8080")
	if err != nil {
		log.Fatalf("could not listen: %s", err.Error())
	}
}
