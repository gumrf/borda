package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ApiHandler struct {
	logger *zap.Logger
}

func NewApiHandler(logger *zap.Logger) *ApiHandler {
	return &ApiHandler{
		logger: logger,
	}
}

func (h *ApiHandler) Init() *gin.Engine {
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

func (h *ApiHandler) health(c *gin.Context) {
	healthResponse := struct{ Message string }{Message: "health"}

	c.JSON(http.StatusOK, healthResponse)
}
