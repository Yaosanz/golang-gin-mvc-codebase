package http

import (
	"context"
	"go-starter-app/app/http/middleware"
	"go-starter-app/app/http/routes"
	"go-starter-app/config"
	"go-starter-app/docs/httpdoc"
	"go-starter-app/interfaces"
	"go-starter-app/pkg/server"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Kernel struct {
	router     *gin.Engine
	http       *server.Http
	middleware *middleware.Middleware
	app        interfaces.KernelDependencies
}

// NewKernel initializes the HTTP Kernel with routes and middleware
func NewKernel(ctx context.Context, app interfaces.KernelDependencies) *Kernel {
	router := gin.New()
	router.MaxMultipartMemory = 100 << 20 // 100MB

	k := &Kernel{
		router: router,
		http:   server.NewHttp(ctx, app.GetConfig(), router), // Register server
		app:    app,
	}

	// setup middleware
	k.middleware = middleware.NewMiddleware()
	k.middleware.Register(k.app.(interfaces.IAppDependencies))
	k.router.Use(k.middleware.GetGlobalMiddleware()...)

	k.registerRoutes()                // register routes
	k.setGinMode(app.GetConfig())     // set gin mode
	k.setSwaggerInfo(app.GetConfig()) // set swagger info
	return k
}

// setGinMode gin mode switcher based on APP_ENV
func (k *Kernel) setGinMode(cfg *config.Config) {
	modeMap := map[string]string{
		"development": gin.DebugMode,
		"dev":         gin.DebugMode,
		"debug":       gin.DebugMode,
		"production":  gin.ReleaseMode,
		"prod":        gin.ReleaseMode,
		"test":        gin.TestMode,
	}

	mode, ok := modeMap[cfg.App().Env]
	if !ok {
		log.Printf(" - Invalid mode: %s. Falling back to default: %s", cfg.App().Env, gin.DebugMode)
		gin.SetMode(gin.DebugMode)
		return
	}
	gin.SetMode(mode)
}

// setSwaggerInfo setup swagger documentation info
func (k *Kernel) setSwaggerInfo(config *config.Config) {
	httpdoc.SwaggerInfo.Title = config.Swagger().Title
	httpdoc.SwaggerInfo.Description = config.Swagger().Description
	httpdoc.SwaggerInfo.Version = config.Swagger().Version
	httpdoc.SwaggerInfo.Host = config.Swagger().Host
	httpdoc.SwaggerInfo.BasePath = config.Swagger().BasePath
	schemas := strings.Split(config.Swagger().Schema, ",")
	httpdoc.SwaggerInfo.Schemes = schemas
}

// registerRoutes setup routes
func (k *Kernel) registerRoutes() {
	// swagger docs route
	k.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Redis health check (testing koneksi Redis / Docker Desktop)
	k.router.GET("/health/redis", func(c *gin.Context) {
		rdb := k.app.GetRedis()
		if rdb == nil {
			c.JSON(503, gin.H{"success": false, "message": "Redis tidak terkoneksi", "redis": "unavailable"})
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()
		if err := rdb.Ping(ctx).Err(); err != nil {
			c.JSON(503, gin.H{"success": false, "message": "Redis ping gagal", "error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"success": true, "message": "Redis OK", "redis": "connected"})
	})

	// register application http routes
	routes.Register(k.router, k.app, k.middleware)
}

// Server get http server instance
func (k *Kernel) Server() *server.Http {
	return k.http
}

// Run starts the HTTP server
func (k *Kernel) Run() {
	if err := k.http.Start(); err != nil {
		log.Fatalf("Failed to start HTTP server: \n%v", err)
	}
}

// Shutdown gracefully shuts down the HTTP server
func (k *Kernel) Shutdown() {
	if err := k.http.Shutdown(); err != nil {
		log.Printf("Failed to shutdown HTTP server: \n%v", err)
	}
}
