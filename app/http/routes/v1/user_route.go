package v1

import (
	v1 "go-starter-app/app/http/controllers/v1"
	"go-starter-app/app/http/middleware"
	"go-starter-app/interfaces"

	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.RouterGroup, app interfaces.KernelDependencies, middleware *middleware.Middleware) {
	controller := v1.NewUserController(app)
	route := router.Group("/users")
	{
		route.GET("/", controller.FindAll)
		route.GET("/:id", controller.FindByID)
		route.POST("/", controller.Create)
		route.PATCH("/:id", controller.Update)
		route.DELETE("/:id", controller.Delete)
	}
}
