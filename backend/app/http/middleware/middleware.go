package middleware

import (
	"go-starter-app/app/services"
	"go-starter-app/interfaces"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
	globalMiddleware map[string][]gin.HandlerFunc
	groupMiddleware  map[string][]gin.HandlerFunc
	routeMiddleware  map[string]gin.HandlerFunc
}

func NewMiddleware() *Middleware {
	return &Middleware{
		globalMiddleware: make(map[string][]gin.HandlerFunc),
		groupMiddleware:  make(map[string][]gin.HandlerFunc),
		routeMiddleware:  make(map[string]gin.HandlerFunc),
	}
}

// Register ALL middleware here
func (m *Middleware) Register(app interfaces.IAppDependencies) {
	// Global middleware
	m.globalMiddleware["global"] = []gin.HandlerFunc{
		gin.Logger(),
		gin.Recovery(),
		CorsMiddleware(),
	}

	// Group middleware 
	m.groupMiddleware["api"] = []gin.HandlerFunc{}
	m.groupMiddleware["web"] = []gin.HandlerFunc{}

	// Route middleware
	secureAuth := NewSecureAuthMiddleware(app.GetService().GetAuthService().(*services.SecureAuthService), app.GetConfig().Jwt())
	m.routeMiddleware["jwt"] = secureAuth.AuthRequired()
	m.routeMiddleware["secure_auth"] = secureAuth.AuthRequired()
}



// ===== GETTERS =====

func (m *Middleware) GetGlobalMiddleware() []gin.HandlerFunc {
	return m.globalMiddleware["global"]
}

func (m *Middleware) GetGroupMiddleware(name string) []gin.HandlerFunc {
	return m.groupMiddleware[name]
}

func (m *Middleware) GetRouteMiddleware(name string) gin.HandlerFunc {
	return m.routeMiddleware[name]
}
