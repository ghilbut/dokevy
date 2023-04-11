package terraform

import "github.com/gin-gonic/gin"

const (
	httpMethodLock   = "LOCK"
	httpMethodUnlock = "UNLOCK"
	path             = "/:name"
)

func RegTerraformHandler(g *gin.RouterGroup) {
	states := g.Group("/states")
	{
		states.GET(path, HandleState)
		states.POST(path, HandleStateUpdate)
		states.Handle(httpMethodLock, path, HandleStateLock)
		states.Handle(httpMethodUnlock, path, HandleStateUnlock)
	}
	secrets := g.Group("/secrets")
	{
		secrets.GET(path, HandleSecret)
		secrets.POST(path, HandleSecretUpdate)
		secrets.Handle(httpMethodLock, path, HandleSecretLock)
		secrets.Handle(httpMethodUnlock, path, HandleSecretUnlock)
	}
}
