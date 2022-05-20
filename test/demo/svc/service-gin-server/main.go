package main

import "github.com/gin-gonic/gin"

func main() {
	// Gin is used as http server in "lab.weave.nl/nid/nid-core/pkg/utilities/httpserver"
	// Used by a number of services
	app := gin.Default()

	// define GET endpoint
	app.GET("/endpoint/get", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Hello world",
		})
	})

	// define POST endpoint
	app.POST("/endpoint/post", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Hello world",
		})
	})

	// run the app
	app.Run()
}
