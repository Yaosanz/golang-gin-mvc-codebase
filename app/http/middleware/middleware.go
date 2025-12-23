package middleware

import "github.com/gin-gonic/gin"

type Middleware struct {
	globalMiddleware []gin.HandlerFunc
	groupsMiddleware map[string][]gin.HandlerFunc
	routeMiddleware  map[string]gin.HandlerFunc
}

func NewMiddleware() *Middleware {
	return &Middleware{
		globalMiddleware: []gin.HandlerFunc{},
		groupsMiddleware: map[string][]gin.HandlerFunc{},
		routeMiddleware:  map[string]gin.HandlerFunc{},
	}
}

// Register all middleware here
func (m *Middleware) Register() {
	// Global middleware
	m.globalMiddleware = []gin.HandlerFunc{
		gin.Logger(),
		gin.Recovery(),
		CorsMiddleware(),
		// Add more global middleware here
	}

	// Predefine common groups
	m.groupsMiddleware["api"] = []gin.HandlerFunc{
		// Add your API group middleware here
	}
	m.groupsMiddleware["web"] = []gin.HandlerFunc{
		// Add your Web group middleware here
	}

	// Define route middleware
	m.routeMiddleware["auth"] = AuthMiddleware()
	m.routeMiddleware["throttle"] = ThrottleMiddleware()
}

func (m *Middleware) GetGlobalMiddleware() []gin.HandlerFunc {
	return m.globalMiddleware
}

func (m *Middleware) GetGroupMiddleware(name string) []gin.HandlerFunc {
	return m.groupsMiddleware[name]
}

func (m *Middleware) GetRouteMiddleware(name string) gin.HandlerFunc {
	return m.routeMiddleware[name]
}
