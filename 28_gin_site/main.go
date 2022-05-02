package main

import (
	"log"
	"net/http"

	//"fmt"
	"github.com/gin-gonic/gin"
	//"vendor/gin/gin"
)

func main() {
	router := gin.Default()
	router.GET("/user/:name", func(c *Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})

	router.POST("/user/", func(c *Context) {
		message := c.PostForm("message")
		jack := c.DefaultPostForm("jack", "admin")
		log.Println(message)
		if message == "" {
			c.JSON(200, H{
				"status":  "ok",
				"message": nil,
				"jack":    jack,
			})
		} else {
			c.JSON(200, H{
				"status":  "ok",
				"message": message,
				"jack":    jack,
			})
		}

	})

	// This handler will add a new router for /user/groups.
	// Exact routes are resolved before param routes, regardless of the order they were defined.
	// Routes starting with /user/groups are never interpreted as /user/:name/... routes
	//router.GET("/user/groups", func(c *gin.Context) {
	//	c.String(http.StatusOK, "The available groups are [...]", name)
	//})

	//router.PUT("/somePut", putting)
	//router.DELETE("/someDelete", deleting)
	//router.PATCH("/somePatch", patching)
	//router.HEAD("/someHead", head)
	//router.OPTIONS("/someOptions", options)

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

//func (p person) fullName() string {
//	return p.first + p.last
//}
//func posting(c *gin.Context) string {
// 	msg = c.JSON(200, gin.H{
//		"message": "postsome",
//	})
//	return msg
//}
