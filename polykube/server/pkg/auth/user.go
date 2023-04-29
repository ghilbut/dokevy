package auth

import "github.com/gin-gonic/gin"

type User struct {
}

func GetUser(ctx *gin.Context) *User {
	v, exists := ctx.Get("user")
	if !exists {
		return nil
	}
	return v.(*User)
}
