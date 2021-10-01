package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const TSA_SERVICE = "http://localhost:2020"

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("web/templates/*")
	r.Static("/assets", "web/static")

	// Serve index page
	r.GET("/", func(c *gin.Context) {
		if c.GetHeader("Content-Type") == "application/timestamp-query" {
			fmt.Println("Timestamp request")
		}
		c.HTML(http.StatusOK, "index.go.html", gin.H{})
	})

	// Proxy download CA certificate
	r.GET("/certs/ca.pem", func(c *gin.Context) {
		resp, err := http.Get(fmt.Sprintf("%s/ca.pem", TSA_SERVICE))
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to request CA certificate from service, err: %s", err.Error()), nil)
		}
		c.DataFromReader(http.StatusOK, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	})

	// Proxy download TSA certificate
	r.GET("/certs/tsa.pem", func(c *gin.Context) {
		resp, err := http.Get(fmt.Sprintf("%s/tsa_cert.pem", TSA_SERVICE))
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to request TSA certificate from service, err: %s", err.Error()), nil)
		}
		c.DataFromReader(http.StatusOK, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	})

	// Sign a document
	r.POST("/sign", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "Vale crack",
		})
	})

	r.Run(":8000")
}
