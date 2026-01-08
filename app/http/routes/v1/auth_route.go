package v1

import (
	"go-starter-app/app/http/controllers/v1"
	"go-starter-app/app/http/middleware"
	"go-starter-app/interfaces"

	"github.com/gin-gonic/gin"
)

func AuthRoute(
	router *gin.RouterGroup,
	app interfaces.IAppDependencies,
	mw *middleware.Middleware,
) {
	controller := v1.NewAuthController(app)

	auth := router.Group("/auth")
	{
		auth.POST("/login", controller.Login)
		auth.POST("/register", controller.Register)
	}
}
