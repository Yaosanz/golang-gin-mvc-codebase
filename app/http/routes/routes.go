package routes

import (
	"go-starter-app/app/http/middleware"
	v1Route "go-starter-app/app/http/routes/v1"
	"go-starter-app/interfaces"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.Engine, app interfaces.KernelDependencies, middleware *middleware.Middleware) {
	// api routes group (use /api prefix and api middleware)
	apiGroup := router.Group("/api", middleware.GetGroupMiddleware("api")...)
	{
		// v1 routes group
		v1RouteGroup := apiGroup.Group("/v1")
		v1Route.UserRoute(v1RouteGroup, app, middleware)
	}

	// web routes group (use web middleware)
	web := router.Group("/", middleware.GetGroupMiddleware("web")...)
	{
		web.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "welcome to web"})
		})
	}
}
