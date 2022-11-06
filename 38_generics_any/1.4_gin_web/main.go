package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var mode = setMode()

const Version = "0.1022.1"

func setMode() bool {
	gin.SetMode(gin.DebugMode)
	return true
}
func main() {
	r := gin.Default()
	gin.SetMode(gin.DebugMode)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run("0.0.0.0:2998") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
