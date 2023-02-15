package main

import (
	"net/http"
	"net/http/httputil"
	// external packages
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.NoRoute(ReverseProxy())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}

func ReverseProxy() gin.HandlerFunc {

	target := "localhost:3000"

	return func(c *gin.Context) {
		director := func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = target
			req.Header["my-header"] = []string{"Ghilbut"}
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
