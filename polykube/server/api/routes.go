package api

import (
	"github.com/ghilbut/polykube/api/v1/dex"
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

	v1 := r.Group("/v1")
	// Dex IDP
	v1.GET("/dex/clients", dex.ListClients)
	v1.POST("/dex/clients", dex.CreateClient)
	v1.GET("/dex/clients/:id", dex.GetClient)
	v1.PUT("/dex/clients/:id", dex.UpdateClient)
	v1.DELETE("/dex/clients/:id", dex.DeleteClient)
	// Terraform secrets
	v1.GET("/terraform/secrets/:name", terraform.HandleGetSecret)
	v1.POST("/terraform/secrets/:name", terraform.HandleCreateSecret)
	v1.DELETE("/terraform/secrets/:name", terraform.HandleDeleteSecret)
	v1.POST("/terraform/secrets/:name/values", terraform.HandleCreateSecretValue)
	v1.PUT("/terraform/secrets/:name/values/:key", terraform.HandleUpdateSecretValue)
	v1.DELETE("/terraform/secrets/:name/values/:key", terraform.HandleDeleteSecretValue)
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
