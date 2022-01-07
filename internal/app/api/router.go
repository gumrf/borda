package api

import (
	"github.com/gin-gonic/gin"
	"time"
	"fmt"
	"borda/internal/app/setup"
)

type ApiHandler struct {
    router *gin.Engine
}

func NewRoutes() *gin.Engine {
    r := ApiHandler{
        router: gin.Default(),
    }
	logger, err := setup.GetReadyLogger()
	if err != nil {
		fmt.Println("Logger is die :)")
	}
	logger.Println("HEAL CHECK HANDLER")
	r.router.GET("/", func (c *gin.Context)  {
		
		c.JSON(200, gin.H{
        	"message": "OK",
			"time" : time.Now().Format(time.UnixDate),
    	})
	})

    v1 := r.router.Group("/v1")
    r.addHealCheck(v1)
    return r.router
}
