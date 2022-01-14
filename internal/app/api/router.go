package api

import (
	"borda/internal/app/setup"
	"time"

	"github.com/gin-gonic/gin"
)

type ApiHandler struct {
	router *gin.Engine
}

func NewRoutes() *gin.Engine {
	r := ApiHandler{
		router: gin.Default(),
	}
	logger := setup.GetLogger()
	logger.Debug("Route init")

	r.router.GET("/", func(c *gin.Context) {

		c.JSON(200, gin.H{
			"message": "OK",
			"time":    time.Now().Format(time.UnixDate),
		})
	})

	v1 := r.router.Group("/v1")
	r.addHealCheck(v1)
	return r.router
}
