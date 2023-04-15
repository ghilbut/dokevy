package terraform

import (
	"errors"
	"net/http"
	"net/http/httputil"

	// external packages
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	httpMethodLock   = "LOCK"
	httpMethodUnlock = "UNLOCK"
)

func AddRoutes(g *gin.RouterGroup) {
	const (
		userpath  = "/users/:user"
		grouppath = "/groups/:group"
	)

	users := g.Group(userpath, middleware)
	addStateRoutes(users)
	addSecretRoutes(users)

	groups := g.Group(grouppath, middleware)
	addStateRoutes(groups)
	addSecretRoutes(groups)
}

func middleware(ctx *gin.Context) {
	w := ctx.Writer

	dump, _ := httputil.DumpRequest(ctx.Request, true)
	log.Trace(string(dump))

	username, password, ok := ctx.Request.BasicAuth()
	if !ok {
		const (
			k = "WWW-Authenticate"
			v = `Basic realm="Give username and password", charset="UTF-8"`
			m = "No basic auth present"
		)
		ctx.Writer.Header().Set(k, v)
		ctx.AbortWithError(http.StatusUnauthorized, errors.New(m))
		return
	}

	if !isAuthorised(ctx, username, password) {
		const (
			k = "WWW-Authenticate"
			v = `Basic realm="Give username and password", charset="UTF-8"`
			m = "Invalid username or password"
		)
		w.Header().Set(k, v)
		ctx.AbortWithError(http.StatusUnauthorized, errors.New(m))
		return
	}

	ctx.Next()
}

func isAuthorised(ctx *gin.Context, username, password string) bool {
	// TODO(ghilbut): check username and password
	return true
}
