package terraform

import (
	"encoding/json"
	"net/http"
	// external packages
	"github.com/gin-gonic/gin"
)

func addSecretRoutes(g *gin.RouterGroup) {
	const path = "/secrets/:name"
	g.GET(path, HandleGetSecret)
	g.POST(path, HandleUpdateSecret)
	g.Handle(httpMethodLock, path, HandleSecretLock)
	g.Handle(httpMethodUnlock, path, HandleSecretUnlock)
}

func HandleGetSecret(ctx *gin.Context) {
	v := `{
  "version": 4,
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
  }
}`
	s := Secret{}
	json.Unmarshal([]byte(v), &s)

	ctx.JSON(200, s)
}

func HandleUpdateSecret(ctx *gin.Context) {
	ctx.Status(http.StatusForbidden)
}

func HandleSecretLock(ctx *gin.Context) {
	ctx.Status(http.StatusForbidden)
}

func HandleSecretUnlock(ctx *gin.Context) {
	ctx.Status(http.StatusForbidden)
}
