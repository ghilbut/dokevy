package terraform

import (
	"net/http"
	// external packages
	"github.com/gin-gonic/gin"
)

func HandleState(ctx *gin.Context) {

}

func HandleStateUpdate(ctx *gin.Context) {
	if ctx.Request.Method != http.MethodPost {
		ctx.AbortWithStatus(http.StatusMethodNotAllowed)
		return
	}
}

func HandleStateLock(ctx *gin.Context) {
	if ctx.Request.Method != "LOCK" {
		ctx.AbortWithStatus(http.StatusMethodNotAllowed)
		return
	}
}

func HandleStateUnlock(ctx *gin.Context) {
	if ctx.Request.Method != "UNLOCK" {
		ctx.AbortWithStatus(http.StatusMethodNotAllowed)
		return
	}
}
