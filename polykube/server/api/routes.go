package api

import (
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"path"

	// external packages
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	// project packages
	apiv1 "github.com/ghilbut/polykube/api/v1/rest"
	"github.com/ghilbut/polykube/api/v1/terraform"
)

func AddRoutes(r *gin.Engine) {

	v1 := r.Group("/v1")
	v1.GET("/hello", apiv1.Helloworld)

	t := v1.Group("/terraform")
	terraform.AddRoutes(t)

	r.GET("/metrics", GetMetricsHandler())
	r.GET("/ping", GetTestHandler())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusFound, path.Join(c.Request.RequestURI, "index.html"))
	})
}

func GetMetricsHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func GetTestHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	}
}
