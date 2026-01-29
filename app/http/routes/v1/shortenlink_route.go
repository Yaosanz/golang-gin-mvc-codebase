package v1

import (
	controllerV1 "go-starter-app/app/http/controllers/v1"
	"go-starter-app/app/http/middleware"
	"go-starter-app/app/services"
	"go-starter-app/interfaces"

	"github.com/gin-gonic/gin"
)

func ShortenlinkRoute(
	router *gin.RouterGroup,
	app interfaces.KernelDependencies,
	mw *middleware.Middleware,
) {
	controller := controllerV1.NewShortenlinkController(app)

	// Get secure auth middleware
	secureAuth := middleware.NewSecureAuthMiddleware(app.GetService().GetAuthService().(*services.SecureAuthService), app.GetConfig().Jwt())

	// PROTECTED CRUD (JWT REQUIRED)
	api := router.Group(
		"/shorten-links",
		secureAuth.AuthRequired(),
	)
	{
		// READ operations - users can read own, admin/cms can read all
		api.GET("/", secureAuth.PermissionRequired("shortenlink:read:own"), controller.FindAll)           // List own (users) or all (admin/cms)
		api.GET("/:id", secureAuth.PermissionRequired("shortenlink:read:own"), controller.FindByID)        // Get by ID if own (users) or any (admin/cms)

		// CREATE operation
		api.POST("", secureAuth.PermissionRequired("shortenlink:create"), controller.Create)         // Create new

		// UPDATE operations - users can update own, admin/cms can update any
		api.PATCH("/:id", secureAuth.PermissionRequired("shortenlink:update:own"), controller.Update)    // Update own (users) or any (admin/cms)

		// DELETE operations - users can delete own, admin/cms can delete any
		api.DELETE("/:id", secureAuth.PermissionRequired("shortenlink:delete:own"), controller.Delete)   // Delete own (users) or any (admin/cms)
	}

	// Compatibility routes for Postman collection (no dash, code-based operations)
	compat := router.Group("/shortenlinks")
	{
		compat.GET("/:code", controller.GetByCode)
		compat.POST("", secureAuth.AuthRequired(), controller.CreateByCode)
		compat.PUT("/:code", secureAuth.AuthRequired(), controller.UpdateByCode)
	}

}
