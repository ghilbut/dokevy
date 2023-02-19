package main

import (
	"fmt"
	"go.elastic.co/apm/module/apmhttp/v2"
	"net/http"
	"net/http/httputil"
	// external packages
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.elastic.co/apm/module/apmgin/v2"
	"go.elastic.co/apm/v2"
)

func init() {
	// NOTE(ghilbut): before running main
	// can register prometheus custom metrics
	//   * https://gabrieltanner.org/blog/collecting-prometheus-metrics-in-golang/
}

func main() {
	r := gin.Default()
	r.NoRoute(ReverseProxy())
	r.SetTrustedProxies([]string{"localhost:3000"})

	r.Use(
		apmgin.Middleware(r),
	)

	r.GET("/metrics", func() gin.HandlerFunc {
		h := promhttp.Handler()
		return func(c *gin.Context) {
			h.ServeHTTP(c.Writer, c.Request)
		}
	}())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Run()
}

func ReverseProxy() gin.HandlerFunc {

	scheme := "http"
	target := "localhost:3000"

	return func(c *gin.Context) {
		director := func(req *http.Request) {
			req.URL.Scheme = scheme
			req.URL.Host = target
			if tx := apm.TransactionFromContext(c.Request.Context()); tx != nil {
				tx.Name = fmt.Sprintf("%s %s", req.Method, req.URL.Path)
				apmhttp.SetHeaders(req, tx.TraceContext(), false)
			}
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
