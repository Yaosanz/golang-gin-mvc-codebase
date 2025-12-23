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
)

func main() {
	fmt.Println(" - App Timezone: ", time.Now().Location().String())
	fmt.Println(" - App Version: ", os.Getenv("APP_VERSION"))

	// initialize application logger
	// todo: move to app/bootstrap.go (centralize all initializations)
	appLogger := log.New(os.Stdout, "[APP] ", log.LstdFlags)

	// initialize app bootstrap
	app, err := bootstrap.NewAppBootstrap()
	if err != nil {
		appLogger.Fatalf("Failed to initialize app: \n%v", err)
	}

	// Create a context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start http server
	kernel := http.NewKernel(ctx, app)
	kernel.Run()

	// Start cron jobs if enabled
	// todo: move to app/bootstrap.go (centralize all initializations)
	if app.GetConfig().App().EnableCron {
		// register all available jobs
		err = jobs.RunRegisteredJobs(app)
		if err != nil {
			appLogger.Fatalf("Failed to register jobs: \n%v", err)
		}

		// Start the scheduler in a separate goroutine
		go app.GetScheduler().Start()
	}

	// Channel to capture OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal
	<-quit
	appLogger.Println("Shutting down application...")

	// Shutdown HTTP server
	kernel.Shutdown()

	// stop scheduler if running
	// todo: move to app/bootstrap.go (centralize all initializations)
	if app.GetScheduler() != nil {
		app.GetScheduler().Stop()
	}

	// Close database connection
	// todo: move to app/bootstrap.go (centralize all initializations)
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

	appLogger.Println("Scheduler stopping...")
	// No need to explicitly stop the scheduler as it stops when the context is canceled.

	appLogger.Println("Application shutdown complete")
}
