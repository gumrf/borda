package api

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Errors []ErrorObject `json:"errors"`
}

type ErrorObject struct {
	Status string `json:"status"`
	Code   string `json:"code"`
	Title  string `json:"title,omitempty"`
	Detail string `json:"detail,omitempty"`
	// Source *SourceObject `json:"source,omitempty"`
}

// type SourceObject struct {
// 	Parameter Parameter `json:"parametr"`
// }

// type Parameter string

// TODO: error responses for common http errors

func NewAPIErrorResponse(errors ...ErrorObject) *ErrorResponse {
	return &ErrorResponse{
		Errors: errors,
	}
}

func APIErrorResponse(c *fiber.Ctx, errors ...ErrorObject) error {
	return c.Status(fiber.StatusNotFound).JSON(&ErrorResponse{
		Errors: errors,
	},

	// NewAPIErrorResponse(ErrorObject{
	// 	Status: strconv.Itoa(fiber.StatusNotFound),
	// 	Code:   http.StatusText(fiber.StatusNotFound),
	// }),
	)
}

func BadRequest(title, detail string) ErrorObject {
	return ErrorObject{
		Status: strconv.Itoa(fiber.StatusBadRequest),
		Code:   http.StatusText(fiber.StatusBadRequest),
		Title:  title,
		Detail: detail,
	}
}

func NotFound(title, detail string) ErrorObject {
	return ErrorObject{
		Status: strconv.Itoa(fiber.StatusNotFound),
		Code:   http.StatusText(fiber.StatusNotFound),
		Title:  title,
		Detail: detail,
	}
}
