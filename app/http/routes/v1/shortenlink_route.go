package v1

import (
	controllerV1 "go-starter-app/app/http/controllers/v1"
	"go-starter-app/app/http/middleware"
	"go-starter-app/interfaces"

	"github.com/gin-gonic/gin"
)

func ShortenlinkRoute(
	router *gin.RouterGroup,
	app interfaces.KernelDependencies,
	middleware *middleware.Middleware,
) {
	controller := controllerV1.NewShortenlinkController(app)

	// PROTECTED CRUD (JWT REQUIRED)
	api := router.Group(
		"/shorten-links",
		middleware.GetRouteMiddleware("jwt"),
	)
	{
		api.GET("/", controller.FindAll)
		api.POST("", controller.Create)
		api.GET("/:id", controller.FindByID)
		api.PATCH("/:id", controller.Update)
		api.DELETE("/:id", controller.Delete)
	}

	// PUBLIC REDIRECT
	router.GET("/r/:code", controller.Redirect)
}
