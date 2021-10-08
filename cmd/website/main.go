package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"net/url"

	ts "github.com/digitorus/timestamp"
	"github.com/gin-gonic/gin"
)

func main() {
	tsaHost := flag.String("host", "http://localhost", "TSA service's hostname")
	tsaPort := flag.Int("port", 318, "TSA service's port number")
	flag.Parse()

	tsaURL, err := url.Parse(fmt.Sprintf("%s:%d", *tsaHost, *tsaPort))
	if err != nil {
		panic("Not a valid URL")
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
	r.GET("/certs/ca.pem", func(c *gin.Context) {
		resp, err := http.Get(fmt.Sprintf("%s/ca.pem", tsaURL.String()))
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to request CA certificate from service, err: %s", err.Error()), nil)
		}
		c.DataFromReader(http.StatusOK, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	})

	// Proxy download TSA certificate
	r.GET("/certs/tsa.pem", func(c *gin.Context) {
		resp, err := http.Get(fmt.Sprintf("%s/tsa_cert.pem", tsaURL.String()))
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to request TSA certificate from service, err: %s", err.Error()), nil)
		}
		c.DataFromReader(http.StatusOK, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	})

	// Sign a document
	r.POST("/sign", func(c *gin.Context) {
		file, header, err := c.Request.FormFile("doc")
		if err != nil {
			panic(err)
		}

		tsReq, err := ts.CreateRequest(file, nil)
		if err != nil {
			panic(err)
		}

		req, err := http.NewRequest(http.MethodGet, tsaURL.String(), bytes.NewReader(tsReq))
		if err != nil {
			panic(err)
		}
		req.Header.Set("Content-Type", "application/timestamp-query")

		tsResp, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}

		c.DataFromReader(http.StatusOK, tsResp.ContentLength, tsResp.Header.Get("Content-Type"), tsResp.Body, map[string]string{
			"Content-Disposition": "attachment; filename=\"" + header.Filename + ".tsr\"",
		})
	})

	r.Run(":8000")
}
