package main

import (
	"context"
	"fmt"
	"go-starter-app/app/http"
	"go-starter-app/app/jobs"
	"go-starter-app/bootstrap"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Overload(); err != nil {
		log.Println("⚠️ No .env file found, using system environment variables")
	}

	fmt.Println(" - App Timezone :", time.Now().Location().String())
	fmt.Println(" - App Version  :", os.Getenv("APP_VERSION"))

	// Initialize application logger
	appLogger := log.New(os.Stdout, "[APP] ", log.LstdFlags)

	// Bootstrap application (config, db, cache, etc)
	app, err := bootstrap.NewAppBootstrap()
	if err != nil {
		appLogger.Fatalf("Failed to initialize app:\n%v", err)
	}

	// Root context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start HTTP server
	kernel := http.NewKernel(ctx, app)
	kernel.Run()

	// Start cron jobs if enabled
	if app.GetConfig().App().EnableCron {
		if err := jobs.RunRegisteredJobs(app); err != nil {
			appLogger.Fatalf("Failed to register jobs:\n%v", err)
		}

		go app.GetScheduler().Start()
	}

	// Listen for OS shutdown signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Println("Shutting down application...")

	// Cancel context (notify goroutines)
	cancel()

	// Shutdown HTTP server
	kernel.Shutdown()

	// Stop scheduler if running
	if app.GetScheduler() != nil {
		app.GetScheduler().Stop()
	}

	// Close database connection
	if app.GetDB() != nil {
		sqlDB, err := app.GetDB().DB()
		if err != nil {
			appLogger.Fatalf("Failed to get database instance: %v", err)
		}
		if err := sqlDB.Close(); err != nil {
			appLogger.Fatalf("Failed to close database connection: %v", err)
		}
		appLogger.Println("Database connection closed.")
	}

	appLogger.Println("Application shutdown complete")
}
