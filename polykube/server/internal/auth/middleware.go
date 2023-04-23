package auth

import (
	"fmt"
	// project packages
	. "ghilbut.com/k8single/internal/auth/nextauth"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Middleware(ctx *gin.Context) {
	cookie, err := ctx.Cookie("next-auth.session-token")
	if err == nil {
		var user UserEntity
		query := fmt.Sprintf(
			"SELECT * FROM %s WHERE id = (SELECT user_id FROM %s WHERE session_token = ?)",
			UserTableName,
			SessionTableName)
		db := ctx.MustGet("DB").(*gorm.DB)
		result := db.Raw(query, cookie).Scan(&user)
		if result.Error != nil {
			panic(result.Error)
		}
		if result.RowsAffected != 1 {
			panic("no seesion")
		}
		ctx.Set("user", &user)
	}
	ctx.Next()
}
