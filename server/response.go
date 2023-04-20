package main

import "github.com/gofiber/fiber/v2"

type errorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(errorResponse{Message: message})
}
