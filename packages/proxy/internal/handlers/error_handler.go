package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"time"
)

type DefaultError struct {
	Timestamp time.Time `json:"timestamp"`
	Status    int       `json:"status"`
	Error     string    `json:"error"`
	Exception string    `json:"exception"`
	Message   string    `json:"message"`
	Path      string    `json:"path"`
}

func NewDefaultErrorMessage(uri string, status int, message string, exception error) DefaultError {
	if status == 0 {
		status = http.StatusInternalServerError
	}
	return DefaultError{
		Timestamp: time.Now(),
		Status:    status,
		Error:     http.StatusText(status),
		Message:   message,
		Exception: exception.Error(),
		Path:      uri,
	}
}

func NewDefaultError(uri string, status int, exception error) DefaultError {
	if status == 0 {
		status = http.StatusInternalServerError
	}
	return DefaultError{
		Timestamp: time.Now(),
		Status:    status,
		Error:     http.StatusText(status),
		Exception: exception.Error(),
		Path:      uri,
	}
}

func DefaultErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}
	defErr := NewDefaultError(c.Path(), code, err)
	return c.Status(code).JSON(defErr)
}
