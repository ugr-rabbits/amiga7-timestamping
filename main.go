package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

const TSA_SERVICE = "http://localhost:2020"

func main() {
	remote, err := url.Parse(TSA_SERVICE)
	if err != nil {
		panic(err)
	}

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
	r.GET("/certs/ca", func(c *gin.Context) {
		proxy := httputil.NewSingleHostReverseProxy(remote)
		proxy.Director = func(req *http.Request) {
			req.Header = c.Request.Header
			req.Host = remote.Host
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			req.URL.Path = "/ca.pem"
		}
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	// Proxy download TSA certificate
	r.GET("/certs/tsa", func(c *gin.Context) {
		proxy := httputil.NewSingleHostReverseProxy(remote)
		proxy.Director = func(req *http.Request) {
			req.Header = c.Request.Header
			req.Host = remote.Host
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			req.URL.Path = "/tsa_cert.pem"
		}
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	// Proxy sign a document
	r.POST("/sign", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "Vale crack",
		})
	})

	r.Run(":8000")
}
