package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("web/templates/*")
	r.Static("/assets", "web/static")

	r.GET("/", func(c *gin.Context) {
		if c.GetHeader("Content-Type") == "application/timestamp-query" {
			fmt.Println("Timestamp request")
		}
		c.HTML(http.StatusOK, "index.go.html", gin.H{})
	})

	r.POST("/sign", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "Vale crack",
		})
	})

	r.Run(":8000")
}
