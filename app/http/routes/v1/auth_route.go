package v1

import (
	v1 "go-starter-app/app/http/controllers/v1"
	"go-starter-app/app/http/middleware"
	"go-starter-app/app/services"
	"go-starter-app/interfaces"

	"github.com/gin-gonic/gin"
)

func AuthRoute(
	router *gin.RouterGroup,
	app interfaces.IAppDependencies,
	mw *middleware.Middleware,
) {
	controller := v1.NewAuthController(app)
	secureAuth := middleware.NewSecureAuthMiddleware(app.GetService().GetAuthService().(*services.SecureAuthService), app.GetConfig().Jwt())

	auth := router.Group("/auth")
	{
		auth.POST("/login", controller.Login)
		auth.POST("/register", controller.Register)
		auth.GET("/profile", secureAuth.AuthRequired(), controller.Profile)
	}
}