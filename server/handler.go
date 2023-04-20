package main

import (
	"encoding/json"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/segmentio/kafka-go"
)

type Handler struct {
	KafkaWriter *kafka.Writer
}

func NewHandler(writer *kafka.Writer) *Handler {
	return &Handler{writer}
}

type RequestBody struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (h *Handler) handleGetMetrics(c *fiber.Ctx) error {
	body := new(RequestBody)
	if err := c.BodyParser(body); err != nil {
		return newErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	jsonData, err := json.Marshal(body)
	if err != nil {
		return newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	err = h.KafkaWriter.WriteMessages(c.UserContext(), kafka.Message{
		Value: jsonData,
	})
	if err != nil {
		return newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(http.StatusOK)
}
