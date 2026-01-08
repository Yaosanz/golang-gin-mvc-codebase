package middleware

import (
	"net/http"

	"go-starter-app/app/http/utils"
	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString("role") != "ADMIN" {
			utils.SendError(c, http.StatusForbidden, "Admin only", nil)
			c.Abort()
			return
		}
		c.Next()
	}
}
