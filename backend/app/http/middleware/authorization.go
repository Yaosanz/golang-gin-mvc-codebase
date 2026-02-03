package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"go-starter-app/app/models"
	"go-starter-app/interfaces"
)

func Authorization(app interfaces.IAppDependencies) gin.HandlerFunc {
	return func(c *gin.Context) {

		// user harus ada dari JWT middleware
		userInf, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Unauthorized",
			})
			return
		}

		user, ok := userInf.(*models.User)
		if !ok || user.ID == uuid.Nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Unauthorized",
			})
			return
		}

		// user harus punya role
		if user.RoleID == uuid.Nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Permission denied",
			})
			return
		}

		method := c.Request.Method
		path := c.FullPath()

		var count int64
		err := app.GetDB().
			Table("role_permissions rp").
			Joins("JOIN permissions p ON p.id = rp.permission_id").
			Where(
				"rp.role_id = ? AND p.method = ? AND p.path = ?",
				user.RoleID,
				method,
				path,
			).
			Count(&count).Error

		if err != nil || count == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Permission denied",
			})
			return
		}

		c.Next()
	}
}
