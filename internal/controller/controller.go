package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewController(router *gin.Engine) {
	// health probe
	router.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })
}
