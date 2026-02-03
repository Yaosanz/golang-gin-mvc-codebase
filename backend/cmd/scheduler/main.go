package main

import (
	"context"
	"fmt"
	"go-starter-app/app/jobs"
	"go-starter-app/bootstrap"
	"go-starter-app/config"
	"go-starter-app/pkg/database"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println(" - App Timezone: ", time.Now().Location().String())
	fmt.Println(" - App Version: ", os.Getenv("APP_VERSION"))

	// load configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: \n%v", err)
	}

	// initialize database connection
	_, dbErr := database.NewPostgres(cfg)
	if dbErr != nil {
		log.Fatalf("Failed to initialize database: \n%v", dbErr)
	}

	appBootstrap, err := bootstrap.NewAppBootstrap()
	if err != nil {
		log.Fatalf("Failed to initialize app: \n%v", err)
	}

	// Create a context for graceful shutdown
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Channel to capture OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// register all available jobs
	// register all available jobs
	err = jobs.RunRegisteredJobs(appBootstrap)
	if err != nil {
		log.Fatalf("Failed to register jobs: \n%v", err)
	}

	// Start the scheduler in a separate goroutine
	go appBootstrap.GetScheduler().Start()

	// Wait for shutdown signal
	<-quit
	log.Println("Shutting down application...")

	log.Println("Scheduler stopping...")
	// No need to explicitly stop the scheduler as it stops when the context is canceled.

	log.Println("Application shutdown complete")
}
