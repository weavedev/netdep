//nolint
package main

import (
	"github.com/gin-gonic/gin"
	"lab.weave.nl/internships/tud-2022/static-analysis-project/test/sample/servicecalls"
)

func main() {
	var testService servicecalls.TESTService
	testService.FirstMethod(1, 2, 3)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
