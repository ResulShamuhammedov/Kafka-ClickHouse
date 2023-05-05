package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gofiber/fiber/v2"
)

type HandlerV3 struct {
	Producer *kafka.Producer
}

func NewHandlerV3(producer *kafka.Producer) *HandlerV3 {
	return &HandlerV3{producer}
}

func (h *HandlerV3) HandleGetMetrics(c *fiber.Ctx) error {
	body := new(RequestBody)
	if err := c.BodyParser(body); err != nil {
		return newErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	jsonData, err := json.Marshal(body)
	if err != nil {
		return newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	topic := "test-topic"
	err = h.Producer.Produce(&kafka.Message{
		Value:          jsonData,
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
	}, nil)
	if err != nil {
		return newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(http.StatusOK)
}

func (h *HandlerV3) HandleInsert(c *fiber.Ctx) error {
	for i := 1; i <= 1000000; i++ {
		jsonData, _ := json.Marshal(RequestBody{
			Name: strconv.Itoa(i),
			Age:  i,
		})
		topic := "test-topic"
		err := h.Producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          jsonData,
		}, nil)
		if err != nil {
			return newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
	}
	return c.SendStatus(http.StatusOK)
}
