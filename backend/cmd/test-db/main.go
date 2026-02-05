package main

import (
	"context"
	"log"
	"net"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env from root directory
	envPath := "C:/Users/sandy/golang-gin-mvc-codebase/.env"
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("⚠️  Warning: Could not load .env from %s: %v\n", envPath, err)
		log.Println("   Trying fallback paths...")
		_ = godotenv.Load("../../../.env")
		_ = godotenv.Load(".env")
	}
	
	log.Println("📂 Checking environment variables...")

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Println("❌ DATABASE_URL not set in .env or environment")
		log.Println("Expected format: postgresql://user:password@host:port/database?sslmode=require")
		os.Exit(1)
	}

	log.Println("============================================================")
	log.Println("🧪 SUPABASE CONNECTION TEST")
	log.Println("============================================================")

	// Display connection info (without password)
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbSSL := os.Getenv("DB_SSL_MODE")

	log.Printf("📍 Host: %s\n", dbHost)
	log.Printf("🔌 Port: %s\n", dbPort)
	log.Printf("📦 Database: %s\n", dbName)
	log.Printf("👤 User: %s\n", dbUser)
	log.Printf("🔒 SSL Mode: %s\n", dbSSL)
	log.Println()

	// Attempt connection
	log.Println("🔄 Attempting connection to Supabase...")

	// Configure custom DNS resolver if needed
	ctx := context.Background()
	
	// Try to resolve hostname first
	resolver := net.DefaultResolver
	ips, dnsErr := resolver.LookupIP(ctx, "ip", dbHost)
	if dnsErr != nil {
		log.Printf("⚠️  Warning: DNS lookup failed for %s: %v\n", dbHost, dnsErr)
		log.Println("   Retrying with system resolver...")
	} else {
		log.Printf("✓ DNS resolved: %s -> %v\n", dbHost, ips)
	}

	conn, err := pgx.Connect(ctx, databaseURL)
	if err != nil {
		log.Printf("❌ Connection failed: %v\n", err)
		
		// Check if this is a Tenant/User authentication error
		if strings.Contains(err.Error(), "Tenant or user not found") {
			log.Println("\n⚠️  AUTHENTICATION ERROR - Tenant or user not found")
			log.Println("\n🔍 TROUBLESHOOTING STEPS:")
			log.Println("   1. Verify Supabase project still exists:")
			log.Println("      → Go to https://supabase.com/dashboard")
			log.Println("      → Check if 'my-rest-api-db' project is listed")
			log.Println()
			log.Println("   2. If project exists, get fresh connection credentials:")
			log.Println("      → Open Project Dashboard")
			log.Println("      → Settings > Database > Connection Pooling")
			log.Println("      → Copy the 'Connection string' for Session Pooler")
			log.Println()
			log.Println("   3. Update .env with new DATABASE_URL")
			log.Println()
			log.Println("   4. Verify password hasn't been reset:")
			log.Println("      → Settings > Database > Reset database password")
			log.Println("      → If reset, you'll need new connection string")
		}
		log.Println("\n❓ Other troubleshooting:")
		log.Println("  1. Check DATABASE_URL in .env is correct")
		log.Println("  2. Verify password in Supabase dashboard")
		log.Println("  3. Verify Supabase project is running and accessible")
		log.Println("  4. Check internet connection")
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	// Test query
	log.Println("🔄 Running test query: SELECT version()")

	var version string
	err = conn.QueryRow(context.Background(), "SELECT version()").Scan(&version)
	if err != nil {
		log.Printf("❌ Query failed: %v\n", err)
		os.Exit(1)
	}

	log.Println()
	log.Println("============================================================")
	log.Println("✅ SUCCESS! Connected to Supabase")
	log.Println("============================================================")
	log.Printf("📊 Database Version: %s\n", version)
	log.Println()
	log.Println("✨ Your Supabase connection is working perfectly!")
	log.Println("   Ready for backend API deployment")
}
