package main

import (
	"github.com/gin-gonic/gin"
	controller "iad591server/api"
	"net/http"
)

func main() {
	router := gin.Default()
	api := router.Group("/api")
	{
		api.GET("/test", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"message": controller.Hello("Hieu Dep Trai"),
			})
		})

		api.GET("/post/:id", func(c *gin.Context) {
			id := c.Param("id")
			c.JSON(http.StatusOK, gin.H{
				"message": "This post has id is: " + id,
			})
		})
	}

	router.Run(":8080")
}
