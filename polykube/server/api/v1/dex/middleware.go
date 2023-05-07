package dex

import (
	"errors"
	"net/http"

	// external packages
	"github.com/gin-gonic/gin"

	// project packages
	"github.com/ghilbut/polykube/pkg/auth"
)

const dexContextKey = "DEX"

func Middleware() gin.HandlerFunc {
	conn, err := newConn()
	if err != nil {
		panic(err)
	}

	return func(ctx *gin.Context) {
		setDexConnectionToContext(ctx, conn)

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

func setDexConnectionToContext(ctx *gin.Context, conn *Conn) {
	ctx.Set(dexContextKey, conn)
}

func getDexConnectionFromContext(ctx *gin.Context) *Conn {
	v := ctx.MustGet("DEX")
	return v.(*Conn)
}

func isSessionAuthenticated(ctx *gin.Context) bool {
	user := auth.GetUser(ctx)
	return user != nil
}
