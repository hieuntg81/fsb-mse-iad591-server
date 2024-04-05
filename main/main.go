package main

import (
	controller "fsb-mse-iad591-server/handler"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()
	api := router.Group("/handler")
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

	err := router.Run(":8080")
	if err != nil {
		return
	}
}
