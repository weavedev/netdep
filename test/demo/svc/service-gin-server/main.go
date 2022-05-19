package main

import "github.com/gin-gonic/gin"

func main() {
	// Gin is used as http server in "lab.weave.nl/nid/nid-core/pkg/utilities/httpserver"
	// Used by a number of services
	app := gin.Default()

	// @mark HTTP endpoint "/ping"
	app.GET("/endpoint/get", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Hello world",
		})
	})

	app.POST("/endpoint/post", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Hello world",
		})
	})

	app.Run()
}
