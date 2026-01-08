package v1

import (
	controllerV1 "go-starter-app/app/http/controllers/v1"
	"go-starter-app/app/http/middleware"
	"go-starter-app/interfaces"

	"github.com/gin-gonic/gin"
)

func UserRoute(
	router *gin.RouterGroup,
	app interfaces.KernelDependencies,
	mw *middleware.Middleware,
) {
	controller := controllerV1.NewUserController(app)

	users := router.Group(
		"/users",
		mw.GetRouteMiddleware("jwt"),
	)

	{
		users.GET(
			"",
			middleware.RequirePermission("user:read"),
			controller.FindAll,
		)

		users.GET(
			"/:id",
			middleware.RequirePermission("user:read"),
			controller.FindByID, // ⬅️ WAJIB ADA
		)

		users.POST(
			"",
			middleware.RequirePermission("user:create"),
			controller.Create,
		)

		users.PATCH(
			"/:id",
			middleware.RequirePermission("user:update"),
			controller.Update,
		)

		users.DELETE(
			"/:id",
			middleware.RequirePermission("user:delete"),
			controller.Delete,
		)
	}
}
