package main

import (
	"fmt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"sync"
	"time"

	// external packages
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.elastic.co/apm/module/apmgin/v2"
	"go.elastic.co/apm/module/apmhttp/v2"
	"go.elastic.co/apm/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"k8s.io/client-go/rest"
	"k8s.io/kubectl/pkg/proxy"

	// project packages
	apiv1 "github.com/ghilbut/polykube/api"
	"github.com/ghilbut/polykube/pkg/auth"
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

	// logrus
	log.SetLevel(log.TraceLevel)

	//ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

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
			apmgin.Middleware(r),
			func(ctx *gin.Context) {
				ctx.Set("DB", db)
				ctx.Next()
			},
			auth.Middleware,
		)

		apiv1.AddRoutes(r)
		r.Handle(http.MethodGet, "/pods", func(ctx *gin.Context) {
			userHomeDir, err := os.UserHomeDir()
			if err != nil {
				log.Fatalf("error getting user home dir: %v\n", err)
			}

			kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")
			// kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
			kubeConfig := clientcmd.GetConfigFromFileOrDie(kubeConfigPath)
			if err != nil {
				log.Fatalf("Error getting kubernetes config: %v\n", err)
			}

			// jhkim := kubeConfig.Contexts["jhkim"]
			// auth := kubeConfig.AuthInfos["jhkim"]
			// cluster := kubeConfig.Clusters["jhkim"]
			// restConfig := rest.Config{}

			// kubeConfig.Clusters["jhkim"].ProxyURL = "unix://tmp/kubectl.sock"
			// kubeConfig.Clusters["jhkim"].Server = "/tmp/kubectl.sock"
			// kubeConfig.Clusters["jhkim"].ProxyURL = "http://127.0.0.1:9090"
			kubeConfig.CurrentContext = "jhkim"
			data, err := clientcmd.Write(*kubeConfig)
			if err != nil {
				log.Fatal(err)
			}

			restConfig, err := clientcmd.RESTConfigFromKubeConfig(data)
			if err != nil {
				log.Fatal(err)
			}

			clientset, err := kubernetes.NewForConfig(restConfig)
			if err != nil {
				log.Fatalf("error getting kubernetes config: %v\n", err)
			}

			ops := v1.ListOptions{}
			pods := clientset.CoreV1().Pods("argo")
			list, err := pods.List(ctx.Request.Context(), ops)
			if err != nil {
				log.Fatal(err)
			}
			log.Info(list.Items)

			ctx.JSON(http.StatusOK, list.Items)
		})

		r.NoRoute(ReverseProxy())
		r.SetTrustedProxies([]string{"localhost:3000"})
		r.Run()
	}()

	go func() {
		defer wg.Done()

		const (
			staticDir    = "/tmp"
			apiPrefix    = "/"
			staticPrefix = "/static/"
			unixSocket   = "/tmp/kubectl.sock"
		)

		var keepalive time.Duration

		appendServerPath := false

		clientConfig := &rest.Config{}
		filter := &proxy.FilterServer{
			AcceptPaths: proxy.MakeRegexpArrayOrDie(proxy.DefaultPathAcceptRE),
			RejectPaths: proxy.MakeRegexpArrayOrDie(proxy.DefaultPathRejectRE),
			// AcceptHosts:   proxy.MakeRegexpArrayOrDie(proxy.DefaultHostAcceptRE),
			AcceptHosts:   proxy.MakeRegexpArrayOrDie("^(.*)$"),
			RejectMethods: proxy.MakeRegexpArrayOrDie(proxy.DefaultMethodRejectRE),
		}

		server, err := proxy.NewServer(staticDir, apiPrefix, staticPrefix, filter, clientConfig, keepalive, appendServerPath)
		if err != nil {
			log.Error(err)
			return
		}
		//l, err := server.ListenUnix(unixSocket)
		l, err := server.Listen("127.0.0.1", 9090)
		if err != nil {
			log.Error(err)
			return
		}
		server.ServeOnListener(l)
	}()

	wg.Wait()
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
