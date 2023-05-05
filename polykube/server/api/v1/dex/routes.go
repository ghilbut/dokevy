package dex

import (
	"errors"
	"github.com/ghilbut/polykube/pkg/dex"
	"net/http"

	// external packages
	"github.com/gin-gonic/gin"

	// project packages
	"github.com/ghilbut/polykube/pkg/auth"
)

func AddRoutes(g *gin.RouterGroup) {
	g.Use(Middleware())
	g.GET("/clients", ListClients)
	g.POST("/clients", CreateClient)
	g.GET("/clients/:id", GetClient)
	g.PUT("/clients/:id", UpdateClient)
	g.DELETE("/clients/:id", DeleteClient)
}

func Middleware() gin.HandlerFunc {
	client, err := dex.NewConn()
	if err != nil {
		panic(err)
	}

	return func(ctx *gin.Context) {
		ctx.Set("DEX", client)

		if !isSessionAuthenticated(ctx) {
			return
		}

		const (
			k = "WWW-Authenticate"
			v = `Basic realm="Give username and password", charset="UTF-8"`
			m = "No basic auth present"
		)
		ctx.Writer.Header().Set(k, v)
		ctx.AbortWithError(http.StatusUnauthorized, errors.New(m))
	}
}

func isSessionAuthenticated(ctx *gin.Context) bool {
	user := auth.GetUser(ctx)
	return user != nil
}
