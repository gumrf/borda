package api

import (
	"github.com/gin-gonic/gin"
	"time"
	"borda/internal/app/setup"
)

type ApiHandler struct {
    router *gin.Engine
}

func NewRoutes() *gin.Engine {
    r := ApiHandler{
        router: gin.Default(),
    }

	r.router.GET("/", func (c *gin.Context)  {
		logger := setup.GetLoggerInstance()
		logger.Log.Println("HEAL CHECK HANDLER")
		c.JSON(200, gin.H{
        	"message": "OK",
			"time" : time.Now().Format(time.UnixDate),
    	})
	})

    v1 := r.router.Group("/v1")
    r.addHealCheck(v1)

    return r.router
}
