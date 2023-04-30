package api

import (
	"net/http"
	"path"

	// external packages
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// project packages
	_ "github.com/ghilbut/polykube/api/docs"
	"github.com/ghilbut/polykube/api/v1/terraform"
)

func AddRoutes(r *gin.Engine) {
	r.GET("/metrics", getMetricsHandler())
	r.GET("/healthz", healthCheckHandlerFunc)
	r.GET("/swagger/*any", swaggerHandlerFunc)
	r.GET("/swagger", swaggerRedirectHandlerFunc)

	v1 := r.Group("/v1")
	t := v1.Group("/terraform")
	terraform.AddRoutes(t)
}

func getMetricsHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(ctx *gin.Context) {
		h.ServeHTTP(ctx.Writer, ctx.Request)
	}
}

var healthCheckHandlerFunc gin.HandlerFunc = func(ctx *gin.Context) {
	ctx.String(http.StatusOK, "OK")
}

var swaggerHandlerFunc gin.HandlerFunc = ginSwagger.WrapHandler(swaggerfiles.Handler)

var swaggerRedirectHandlerFunc gin.HandlerFunc = func(ctx *gin.Context) {
	location := path.Join(ctx.Request.RequestURI, "index.html")
	ctx.Redirect(http.StatusFound, location)
}
