package routes

import (
	v1Controller "go-starter-app/app/http/controllers/v1"
	"go-starter-app/app/http/middleware"
	v1Route "go-starter-app/app/http/routes/v1"
	"go-starter-app/interfaces"

	"github.com/gin-gonic/gin"
)

func Register(
	router *gin.Engine,
	app interfaces.KernelDependencies,
	mw *middleware.Middleware,
) {
	api := router.Group("/api", mw.GetGroupMiddleware("api")...)
	{
		v1 := api.Group("/v1")
		v1Route.AuthRoute(v1, app.(interfaces.IAppDependencies), mw)
		v1Route.UserRoute(v1, app, mw)
		v1Route.ShortenlinkRoute(v1, app, mw)

		// Compatibility routes without version prefix to match Postman collection
		v0 := api.Group("")
		v1Route.AuthRoute(v0, app.(interfaces.IAppDependencies), mw)
		v1Route.UserRoute(v0, app, mw)
		v1Route.ShortenlinkRoute(v0, app, mw)
	}

	web := router.Group("/", mw.GetGroupMiddleware("web")...)
	{
		web.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "welcome to web"})
		})

		// Public redirect for short links - Multiple URL formats for flexibility
		controller := v1Controller.NewShortenlinkController(app)
		web.GET("/r/:code", controller.Redirect)           // /r/{code}
		web.GET("/mydigilearn/:code", controller.Redirect) // /mydigilearn/{code} (legacy)
		web.GET("/s/:code", controller.Redirect)           // /s/{code} (short alias)
	}

	// Handle common browser requests to avoid 404 spam
	router.GET("/favicon.ico", func(c *gin.Context) {
		c.Status(204) // No Content
	})

	// Ignore hot-reload files (webpack dev server)
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		// Silently ignore webpack hot-update files
		if len(path) > 15 && path[len(path)-15:] == "hot-update.json" {
			c.Status(204)
			return
		}
		c.JSON(404, gin.H{"error": "route not found"})
	})
}
