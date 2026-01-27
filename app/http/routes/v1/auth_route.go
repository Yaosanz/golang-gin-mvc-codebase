package v1

import (
	v1 "go-starter-app/app/http/controllers/v1"
	"go-starter-app/app/http/middleware"
	"go-starter-app/app/services"
	"go-starter-app/interfaces"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthRoute(
	router *gin.RouterGroup,
	app interfaces.IAppDependencies,
	mw *middleware.Middleware,
) {
	controller := v1.NewAuthController(app)
	secureAuth := middleware.NewSecureAuthMiddleware(app.GetService().GetAuthService().(*services.SecureAuthService), app.GetConfig().Jwt())

	// Rate limiters for auth endpoints
	loginLimiter := middleware.NewRateLimiter(app.GetRedis(), 5, time.Minute)   // 5 login attempts per minute
	refreshLimiter := middleware.NewRateLimiter(app.GetRedis(), 10, time.Minute) // 10 refresh per minute
	registerLimiter := middleware.NewRateLimiter(app.GetRedis(), 3, time.Minute) // 3 register per minute

	auth := router.Group("/auth")
	{
		auth.POST("/login", loginLimiter.Limit(), controller.Login)
		auth.POST("/refresh", refreshLimiter.Limit(), controller.Refresh)
		auth.POST("/register", registerLimiter.Limit(), controller.Register)
		auth.GET("/profile", secureAuth.AuthRequired(), controller.Profile)
	}
}