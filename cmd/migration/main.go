package main

import (
	"database/sql"
	"go-starter-app/config"
	"go-starter-app/pkg/database"
	"log"

	//"os"
	"errors"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

func main() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}

var RootCmd = &cobra.Command{
	Use:   "migration",
	Short: "Database migration tool",
	Long:  "A tool for managing database migrations.",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help() // Show help message when no subcommand is provided
	},
}

func init() {
	RootCmd.AddCommand(upCmd, downCmd, forceCmd, dropCmd)
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Run all up migrations",
	Run: func(cmd *cobra.Command, args []string) {
		runMigrations("up", 0)
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Run all down migrations",
	Run: func(cmd *cobra.Command, args []string) {
		runMigrations("down", 0)
	},
}

var forceCmd = &cobra.Command{
	Use:   "force [version]",
	Short: "Force migration to a specific version",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		version, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalf("Invalid version: %v", err)
		}
		runMigrations("force", version)
	},
}

var dropCmd = &cobra.Command{
	Use:   "drop",
	Short: "Drop the entire database",
	Run: func(cmd *cobra.Command, args []string) {
		runMigrations("drop", 0)
	},
}

func runMigrations(command string, version int) {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Define PostgresSQL connection details
	dsn := database.PostgresDSN(cfg)

	// Open the database connection
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Failed to close the database: %v", err)
		}
	}(db)

	// Verify connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Database connection test failed: %v", err)
	}

	// Create migration driver
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed to create migration driver: %v", err)
	}

	// Initialize migration instance
	migrationPath := "file://" + cfg.MigrationPath
	m, err := migrate.NewWithDatabaseInstance(
		migrationPath,
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("Failed to initialize golang-migrate: %v", err)
	}

	// Handle migration commands
	switch command {
	case "up":
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("Failed to run up migrations: %v", err)
		}
		log.Println("Migrations applied successfully")
	case "down":
		if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("Failed to run down migrations: %v", err)
		}
		log.Println("Migrations rolled back successfully")
	case "force":
		if err := m.Force(version); err != nil {
			log.Fatalf("Failed to force migrations to version %d: %v", version, err)
		}
		log.Printf("Migrations forced to version %d", version)
	case "drop":
		if err := m.Drop(); err != nil {
			log.Fatalf("Failed to drop database: %v", err)
		}
		log.Println("Database dropped successfully")
	default:
		log.Fatalf("Unknown command: %s", command)
	}
}
