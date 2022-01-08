package api

import (
	"github.com/gin-gonic/gin"
)

func (r ApiHandler) addHealCheck(rg *gin.RouterGroup) {
    ping := rg.Group("/ping")
    ping.GET("/", pongFunction)
}

func pongFunction(c *gin.Context) {
    c.JSON(200, gin.H{
        "message": "pong",
    })
}