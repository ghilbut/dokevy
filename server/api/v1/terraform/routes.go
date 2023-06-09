package terraform

import (
	"crypto/sha512"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	// external packages
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	// project packages
	"github.com/ghilbut/polykube/pkg/auth"
	"github.com/ghilbut/polykube/pkg/terraform"
)

func AddRoutes(g *gin.RouterGroup) {
	const (
		secrets_path = "/secrets"
	)
	secrets := g.Group(secrets_path, middleware)
	// Get values by terraform_remote_state
	secrets.GET("/:group", HandleGetSecret)

	secrets.GET("/groups/:name", HandleGetSecret)
	secrets.POST("/groups/:name", HandleCreateSecret)
	secrets.DELETE("/groups/:name", HandleDeleteSecret)
	secrets.POST("/groups/:name/values", HandleCreateSecretValue)
	secrets.PUT("/groups/:name/values/:key", HandleUpdateSecretValue)
	secrets.DELETE("/groups/:name/values/:key", HandleDeleteSecretValue)
}

func middleware(ctx *gin.Context) {
	if !isSessionAuthenticated(ctx) {
		return
	}

	if isBasicAuthenticated(ctx) {
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

func isBasicAuthenticated(ctx *gin.Context) bool {
	if username, password, ok := ctx.Request.BasicAuth(); ok {
		return isBasicAuthorised(ctx, username, password)
	}
	return false
}

func isBasicAuthorised(ctx *gin.Context, username, password string) bool {
	var key string
	query := fmt.Sprintf("SELECT password FROM %s WHERE username = ? ;", terraform.BotTableName)
	db := ctx.MustGet("DB").(*gorm.DB)
	if err := db.Raw(query, username).Row().Scan(&key); err != nil {
		log.Error(err)
		return false
	}
	sum := sha512.Sum512([]byte(password))
	sha := fmt.Sprintf("%x", sum)
	return key == sha
}
