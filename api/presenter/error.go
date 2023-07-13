package presenter

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type ErrorMessage struct {
	Error   string      `json:"error"`
	Message string      `json:"message"`
	Details interface{} `json:"details"`
}

func ErrorResponse(c *fiber.Ctx, httpStatus int, err error) error {
	return c.Status(httpStatus).JSON(fiber.Map{"error": err.Error()})
}

func InternalServerErrorResponse(c *fiber.Ctx, err error) error {
	return ErrorResponse(c, http.StatusInternalServerError, err)
}

func BadRequestResponse(c *fiber.Ctx, err error) error {
	return ErrorResponse(c, http.StatusBadRequest, err)
}

func NotFoundResponse(c *fiber.Ctx, err error) error {
	return ErrorResponse(c, http.StatusNotFound, err)
}
