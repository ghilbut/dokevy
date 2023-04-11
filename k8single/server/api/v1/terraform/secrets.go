package terraform

import (
	"net/http"
	// external packages
	"github.com/gin-gonic/gin"
)

func HandleSecret(ctx *gin.Context) {

}

func HandleSecretUpdate(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusBadRequest)
}

func HandleSecretLock(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusBadRequest)
}

func HandleSecretUnlock(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusBadRequest)
}
