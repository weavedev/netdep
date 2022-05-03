package main

import "github.com/gin-gonic/gin"

func main() {
	// Gin is used as http server in "lab.weave.nl/nid/nid-core/pkg/utilities/httpserver"
	r := gin.Default()

	// @mark HTTP endpoint "/ping"
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
