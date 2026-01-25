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
	}

	web := router.Group("/", mw.GetGroupMiddleware("web")...)
	{
		web.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "welcome to web"})
		})

		// Public redirect for short links
		web.GET("/mydigilearn/:code", func(c *gin.Context) {
			controller := v1Controller.NewShortenlinkController(app)
			controller.Redirect(c)
		})
	}
}
