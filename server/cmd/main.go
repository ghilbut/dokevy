package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"

	// external packages
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.elastic.co/apm/module/apmhttp/v2"
	"go.elastic.co/apm/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	// project packages
	"github.com/ghilbut/polykube/api"
)

func init() {
	// NOTE(ghilbut): before running main
	// can register prometheus custom metrics
	//   * https://gabrieltanner.org/blog/collecting-prometheus-metrics-in-golang/
}

// @title          PolyKube API
// @version        1.0
// @description    Manager for Applications on Kubernetes
//
// @contact.name   ghilbut
// @contact.email  ghilbut@gmail.com
//
// @license.name   The GNU General Public License v3.0
// @license.url    https://www.gnu.org/licenses/gpl-3.0.en.html
//
// @BasePath       /
//
// @securityDefinitions.apikey ApiKey "sentence with space"
// @in header
// @name Authorization

func main() {

	// logrus
	log.SetLevel(log.TraceLevel)

	// gorm
	dsn := "host=localhost user=postgres password=postgrespw dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Seoul"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	//db.AutoMigrate()

	chdir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	configPath := filepath.Join(chdir, "config.yaml")
	fmt.Println(configPath)

	// gin-gonic
	r := gin.New()

	r.Use(
		gin.Recovery(),
		//apmgin.Middleware(r),
		func(ctx *gin.Context) {
			ctx.Set("DB", db)
			ctx.Next()
		},
		//auth.Middleware,
	)

	api.AddRoutes(r)

	// r.NoRoute(ReverseProxy())
	// r.SetTrustedProxies([]string{
	// 	"localhost:3000",
	// })
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
