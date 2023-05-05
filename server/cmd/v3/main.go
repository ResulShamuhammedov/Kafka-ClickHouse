package main

import (
	"log"

	"github.com/ResulShamuhammedov/Kafka-Clickhouse/server/handler"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gofiber/fiber/v2"
)

func main() {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		panic(err)
	}
	defer p.Close()

	handler := handler.NewHandlerV3(p)

	app := fiber.New()

	app.Post("/metric", handler.HandleGetMetrics)
	app.Post("/insert", handler.HandleInsert)

	err = app.Listen(":8080")
	if err != nil {
		log.Fatalf("could not listen: %s", err.Error())
	}
}
