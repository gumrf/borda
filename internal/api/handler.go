package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct{}

func NewHandlers() *Handler {
	return &Handler{}
}

func (h *Handler) Init() *gin.Engine {
	// Init gin router
	router := gin.New()

	// Test route
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Group api routes with /api/v1 prefix
	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.GET("/health", h.health)
		}
	}

	return router
}

func (h *Handler) health(c *gin.Context) {
	type healthResponse struct {
		Message string
	}

	c.JSON(http.StatusOK, healthResponse{
		Message: "health",
	})
}
