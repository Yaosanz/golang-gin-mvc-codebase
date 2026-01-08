package main

import (
	"fmt"
	"go-starter-app/config"
	"go-starter-app/pkg/database"
	"log"

	"github.com/spf13/cobra"
	"gorm.io/gorm"

	// important: this will trigger init() in all seeder files to register them
	_ "go-starter-app/database/seeders"
	"go-starter-app/pkg/database/seeder"
)

var db *gorm.DB

func initDB() {
	if db == nil {
		// Load configuration
		cfg, err := config.NewConfig()
		if err != nil {
			log.Fatalf("Failed to load configuration: %v", err)
		}

		// Initialize database connection
		gormDB, dbErr := database.NewPostgres(cfg)
		if dbErr != nil {
			log.Fatalf("Failed to initialize database: %v", dbErr)
		}

		db = gormDB
	}
}

var rootCmd = &cobra.Command{
	Use:   "seeder",
	Short: "Database seeder CLI",
	Long:  "A CLI tool for running database seeders in Go projects",
}

// Command to run all seeders
var allCmd = &cobra.Command{
	Use:   "run:all",
	Short: "Run all seeders",
	Run: func(cmd *cobra.Command, args []string) {
		initDB()
		seeder.SetOrder([]string{"permission_seeder", "role_seeder", "role_permissions_seeder", "user_seeder"})
		s := seeder.NewSeeder(db)
		if err := s.RunAll(); err != nil {
			seeder.LogFatalf("Seeder failed: %v", err)
		}
	},
}

// Command to run a specific seeder by name
var specificCmd = &cobra.Command{
	Use:   "run:one [name]",
	Short: "Run a specific seeder",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		initDB()
		s := seeder.NewSeeder(db)
		if err := s.RunOne(args[0]); err != nil {
			seeder.LogFatalf("Seeder failed: %v", err)
		}
	},
}

// Command to run all seeders inside a single transaction
var allTxCmd = &cobra.Command{
	Use:   "run:all-tx",
	Short: "Run all seeders inside a single transaction",
	Run: func(cmd *cobra.Command, args []string) {
		initDB()
		s := seeder.NewSeeder(db)
		if err := s.RunAllWithTransaction(); err != nil {
			seeder.LogFatalf("Seeder failed (transaction rolled back): %v", err)
		}
	},
}

// Command to run a specific seeder inside a transaction
var runTxCmd = &cobra.Command{
	Use:   "run:one-tx [name]",
	Short: "Run a specific seeder inside a transaction",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		initDB()
		s := seeder.NewSeeder(db)
		if err := s.RunOneWithTransaction(args[0]); err != nil {
			seeder.LogFatalf("Seeder failed (transaction rolled back): %v", err)
		}
	},
}

// Command to list all registered seeders
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all registered seeders",
	Run: func(cmd *cobra.Command, args []string) {
		names := seeder.ListSeeders()
		if len(names) == 0 {
			fmt.Println("No seeders registered.")
			return
		}

		fmt.Println("Registered Seeders:")
		for i, name := range names {
			fmt.Printf("%d. %s\n", i+1, name)
		}
	},
}

func init() {
	rootCmd.AddCommand(allCmd, specificCmd, allTxCmd, runTxCmd, listCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}
