package dex

import (
	"errors"
	"net/http"
	// external packages
	"github.com/gin-gonic/gin"
	// project packages
	"github.com/ghilbut/polykube/pkg/auth"
)

func AddRoutes(g *gin.RouterGroup) {
	g.Use(middleware)
	g.Handle(http.MethodGet, "/clients", ListClients)
	g.Handle(http.MethodPost, "/clients", CreateClient)
	g.Handle(http.MethodGet, "/clients/:id", GetClient)
	g.Handle(http.MethodPut, "/clients/:id", UpdateClient)
	g.Handle(http.MethodDelete, "/clients/:id", DeleteClient)
}

func middleware(ctx *gin.Context) {
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

func isSessionAuthenticated(ctx *gin.Context) bool {
	user := auth.GetUser(ctx)
	return user != nil
}
