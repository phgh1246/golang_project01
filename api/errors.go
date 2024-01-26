package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	if apiError, ok := err.(APIError); ok {
		return c.Status(apiError.Code).JSON(apiError)
	}
	apiError := NewError(http.StatusInternalServerError, err.Error())
	return c.Status(apiError.Code).JSON(apiError)
}

type APIError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// implementing Error interface
func (e APIError) Error() string {
	return e.Msg
}

func NewError(code int, msg string) APIError {
	return APIError{
		Code: code,
		Msg:  msg,
	}
}

func ErrorUnauthorized() APIError {
	return APIError{
		Code: http.StatusUnauthorized,
		Msg:  "unauthorized request",
	}
}

func ErrorInvalidID() APIError {
	return APIError{
		Code: http.StatusBadRequest,
		Msg:  "invalid ID given",
	}
}

func ErrorBadRequest() APIError {
	return APIError{
		Code: http.StatusBadRequest,
		Msg:  "invalid JSON request",
	}
}

func ErrorResourceNotFound(res string) APIError {
	return APIError{
		Code: http.StatusNotFound,
		Msg:  res + " resource not found",
	}
}
