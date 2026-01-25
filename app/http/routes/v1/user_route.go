package v1

import (
	controllerV1 "go-starter-app/app/http/controllers/v1"
	"go-starter-app/app/http/middleware"
	"go-starter-app/app/services"
	"go-starter-app/interfaces"

	"github.com/gin-gonic/gin"
)

func UserRoute(
	router *gin.RouterGroup,
	app interfaces.KernelDependencies,
	mw *middleware.Middleware,
) {
	controller := controllerV1.NewUserController(app)

	// Get secure auth middleware
	secureAuth := middleware.NewSecureAuthMiddleware(app.GetService().GetAuthService().(*services.SecureAuthService), app.GetConfig().Jwt())

	// Authenticated users can update their own profile without specifying an ID
	selfUsers := router.Group(
		"/users",
		secureAuth.AuthRequired(),
	)
	{
		selfUsers.PUT("", controller.UpdateSelf)
		selfUsers.PUT("/", controller.UpdateSelf)
	}

	// ADMIN ONLY - User Management
	adminUsers := router.Group(
		"/users",
		secureAuth.AuthRequired(),
		secureAuth.AdminRequired(), // Only admin can manage users
	)
	{
		adminUsers.GET("", controller.FindAll)       // Admin can list all users
		adminUsers.GET("/:id", controller.FindByID)  // Admin can view any user
		adminUsers.POST("", controller.Create)       // Admin can create users
		adminUsers.PATCH("/:id", controller.Update)  // Admin can update any user
		adminUsers.PUT("/:id", controller.Update)    // Admin can update any user via PUT
		adminUsers.DELETE("/:id", controller.Delete) // Admin can delete users
	}
}
