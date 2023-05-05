package api

import (
	"github.com/ghilbut/polykube/api/v1/dex"
	"github.com/ghilbut/polykube/pkg/auth"
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
	r.GET("/healthz", getHealthCheckHandler())
	r.GET("/swagger/*any", swaggerHandlerFunc)
	r.GET("/swagger", swaggerRedirectHandlerFunc)

	// API v1 Group
	v1 := r.Group("/v1", auth.Middleware)
	// Dex IDP handlers
	v1dex := v1.Group("/dex", dex.Middleware())
	v1dex.GET("/clients", dex.ListClients)
	v1dex.POST("/clients", dex.CreateClient)
	v1dex.GET("/clients/:id", dex.GetClient)
	v1dex.PUT("/clients/:id", dex.UpdateClient)
	v1dex.DELETE("/clients/:id", dex.DeleteClient)
	// Terraform secrets handlers
	v1terraform := v1.Group("/terraform")
	v1terraform.GET("/secrets/:name", terraform.HandleGetSecret)
	v1terraform.POST("/secrets/:name", terraform.HandleCreateSecret)
	v1terraform.DELETE("/secrets/:name", terraform.HandleDeleteSecret)
	v1terraform.POST("/secrets/:name/values", terraform.HandleCreateSecretValue)
	v1terraform.PUT("/secrets/:name/values/:key", terraform.HandleUpdateSecretValue)
	v1terraform.DELETE("/secrets/:name/values/:key", terraform.HandleDeleteSecretValue)
}

// GetMetrics godoc
// @Summary      Prometheus metrics
// @Description  get metrics for prometheus exporter
// @Tags         system
// @Produce      plain
// @Success      200  {object}  string
// @Failure      500  {object}  string
// @Router       /metrics [get]
func getMetricsHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(ctx *gin.Context) {
		h.ServeHTTP(ctx.Writer, ctx.Request)
	}
}

// HealthCheck godoc
// @Summary      Health check
// @Description  health check for process controller and load balancer
// @Tags         system
// @Produce      plain
// @Success      200  {path}  string
// @Failure      500  {object}  string
// @Router       /healthz [get]
func getHealthCheckHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "OK")
	}
}

var swaggerHandlerFunc gin.HandlerFunc = ginSwagger.WrapHandler(swaggerfiles.Handler)

var swaggerRedirectHandlerFunc gin.HandlerFunc = func(ctx *gin.Context) {
	location := path.Join(ctx.Request.RequestURI, "index.html")
	ctx.Redirect(http.StatusFound, location)
}
