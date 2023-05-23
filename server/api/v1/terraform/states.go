package terraform

import (
	"net/http"
	// external packages
	"github.com/gin-gonic/gin"
)

const (
	httpMethodLock   = "LOCK"
	httpMethodUnlock = "UNLOCK"
)

func addStateRoutes(g *gin.RouterGroup) {
	const path = "/:name"
	g.GET(path, HandleGetState)
	g.POST(path, HandleUpdateState)
	g.Handle(httpMethodLock, path, HandleLockState)
	g.Handle(httpMethodUnlock, path, HandleUnlockState)
}

func HandleGetState(ctx *gin.Context) {
	v := `{
  "version": 4,
  "terraform_version": "1.4.2",
  "serial": 1,
  "lineage": "145ec231-90be-4a48-e96d-00b782f31f82",
  "outputs": {
    "a": {
      "value": "a",
      "type": "string"
    },
    "b": {
      "value": "b",
      "type": "string",
      "sensitive": true
    }
  },
  "resources": [],
  "check_results": null
}`

	ctx.String(http.StatusOK, v)
}

func HandleUpdateState(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}

func HandleLockState(ctx *gin.Context) {
	//ctx.Status(http.StatusNotImplemented)
	ctx.Status(http.StatusOK)
}

func HandleUnlockState(ctx *gin.Context) {
	//ctx.Status(http.StatusNotImplemented)
	ctx.Status(http.StatusOK)
}
