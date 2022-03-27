package api

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ErrorsResponse struct {
	Errors []Error `json:"errors"`
}

type Error struct {
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

func NewErrorResponse(c *fiber.Ctx, status int, title string, detail ...string) error {
	return c.Status(status).JSON(ErrorsResponse{
		Errors: []Error{
			Error{
				Status: strconv.Itoa(status),
				Code:   http.StatusText(status),
				Title:  title,
				Detail: func(detail []string) string {
					if len(detail) > 0 {
						return detail[0]
					}
					return ""
				}(detail),
			},
		},
	})
}

type BadRequestResponse ErrorsResponse
