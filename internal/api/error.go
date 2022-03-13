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

// TODO: adopt NewErrorResponse to acept custom error compatible with
// 		 interface. It should replaces title and detail.
func NewErrorResponse(c *fiber.Ctx, status int, title string) error {
	return c.Status(status).JSON(
		ErrorResponse{
			Errors: []ErrorObject{
				ErrorObject{
					Status: strconv.Itoa(status),
					Code:   http.StatusText(status),
					Title:  title,
				},
			},
		})
}
