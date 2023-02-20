package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	// external packages
	"github.com/gin-gonic/gin"
	"go.elastic.co/apm/module/apmgin/v2"
	"go.elastic.co/apm/module/apmhttp/v2"
	"go.elastic.co/apm/v2"

	apiv1 "ghilbut.com/k8single/api"
)

func init() {
	// NOTE(ghilbut): before running main
	// can register prometheus custom metrics
	//   * https://gabrieltanner.org/blog/collecting-prometheus-metrics-in-golang/
}

// @title          K8Single API
// @version        1.0
// @description    Manager for Applications on Kubernetes

// @contact.name   ghilbut
// @contact.email  ghilbut@gmail.com

// @license.name   MIT License
// @license.url    https://opensource.org/license/mit/

// @BasePath  /v1

func main() {
	r := gin.New()

	r.Use(
		apmgin.Middleware(r),
	)

	apiv1.RegHandler(r)

	r.NoRoute(ReverseProxy())
	r.SetTrustedProxies([]string{"localhost:3000"})
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
