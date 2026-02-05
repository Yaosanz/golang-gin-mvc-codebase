#!/bin/bash
# Database Migration Guide Script
# Quick reference for migration commands and workflows

set -e

MIGRATIONS_DIR="database/migrations"
MIGRATION_CLI="./migration.exe"

echo "================================"
echo "Database Migration Guide"
echo "================================"
echo ""

# Check if migration CLI exists
if [[ ! -f "$MIGRATION_CLI" ]]; then
    echo "⚙️  Building migration CLI..."
    go build -o migration.exe ./cmd/migration
    echo "✅ Migration CLI built"
    echo ""
fi

# Display menu
show_menu() {
    echo "Quick Migration Commands:"
    echo "1. View available migrations"
    echo "2. Apply all pending migrations"
    echo "3. Rollback all migrations"
    echo "4. Force to specific version"
    echo "5. Drop entire database"
    echo "6. Show migration status"
    echo "7. Complete workflow (migrate + seed)"
    echo ""
}

# View migrations
view_migrations() {
    echo "Available Migrations:"
    echo "===================="
    ls -1 "$MIGRATIONS_DIR" | grep -E "\.(up|down)\.sql$" | sed 's/\.up\.sql//' | sed 's/\.down\.sql//' | sort -u | nl
    echo ""
    echo "Location: $MIGRATIONS_DIR/"
    echo ""
}

# Apply migrations
apply_migrations() {
    echo "Applying pending migrations..."
    echo "Command: $MIGRATION_CLI up"
    echo ""
    
    if $MIGRATION_CLI up; then
        echo "✅ Migrations applied successfully"
    else
        echo "❌ Migration failed. Check database connection."
        return 1
    fi
    echo ""
}

# Rollback migrations
rollback_migrations() {
    echo "⚠️  WARNING: Rolling back all migrations will DROP all tables!"
    read -p "Continue? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "Rolling back all migrations..."
        echo "Command: $MIGRATION_CLI down"
        echo ""
        
        if $MIGRATION_CLI down; then
            echo "✅ Rollback successful"
        else
            echo "❌ Rollback failed"
            return 1
        fi
    else
        echo "Rollback cancelled"
    fi
    echo ""
}

# Force version
force_version() {
    read -p "Enter target version: " version
    echo "Forcing migration to version $version..."
    echo "Command: $MIGRATION_CLI force $version"
    echo ""
    
    if $MIGRATION_CLI force "$version"; then
        echo "✅ Forced to version $version"
    else
        echo "❌ Force failed"
        return 1
    fi
    echo ""
}

# Drop database
drop_database() {
    echo "⚠️  WARNING: This will DROP EVERYTHING!"
    echo "All tables, constraints, and data will be deleted."
    read -p "Continue? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "Dropping database..."
        echo "Command: $MIGRATION_CLI drop"
        echo ""
        
        if $MIGRATION_CLI drop; then
            echo "✅ Database dropped"
        else
            echo "❌ Drop failed"
            return 1
        fi
    else
        echo "Drop cancelled"
    fi
    echo ""
}

# Show status
show_status() {
    echo "Migration Status"
    echo "================"
    echo ""
    
    # Count migration files
    UP_COUNT=$(find "$MIGRATIONS_DIR" -name "*.up.sql" | wc -l)
    DOWN_COUNT=$(find "$MIGRATIONS_DIR" -name "*.down.sql" | wc -l)
    
    echo "Migration files found:"
    echo "  UP migrations:   $UP_COUNT"
    echo "  DOWN migrations: $DOWN_COUNT"
    
    if [[ $UP_COUNT -eq $DOWN_COUNT ]] && [[ $UP_COUNT -gt 0 ]]; then
        echo "  Status: ✅ All migrations have paired up/down files"
    else
        echo "  Status: ❌ Mismatch between up and down migrations"
    fi
    echo ""
    
    # Check database connection
    echo "Database Connection:"
    if [[ -f ".env" ]]; then
        echo "  ✅ .env file found"
        # Extract DB info without showing password
        DB_HOST=$(grep "^DB_HOST=" .env | cut -d= -f2)
        DB_NAME=$(grep "^DB_NAME=" .env | cut -d= -f2)
        echo "  Host: $DB_HOST"
        echo "  Database: $DB_NAME"
    else
        echo "  ❌ .env file not found"
    fi
    echo ""
}

# Complete workflow
complete_workflow() {
    echo "Complete Migration Workflow"
    echo "==========================="
    echo ""
    echo "This will:"
    echo "1. Drop existing database (if confirmed)"
    echo "2. Apply all migrations"
    echo "3. Run seeders"
    echo ""
    
    read -p "Drop existing database first? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "Dropping database..."
        if ! $MIGRATION_CLI drop; then
            echo "❌ Drop failed"
            return 1
        fi
    fi
    
    echo ""
    echo "Applying migrations..."
    if ! $MIGRATION_CLI up; then
        echo "❌ Migration failed"
        return 1
    fi
    
    echo ""
    echo "Running seeders..."
    if [[ -f "seeder.exe" ]]; then
        if ./seeder.exe; then
            echo "✅ Seeders completed"
        else
            echo "⚠️  Seeders had issues (migrations still applied)"
        fi
    else
        echo "⚠️  seeder.exe not found. Build it with:"
        echo "   go build -o seeder.exe ./cmd/seeder"
        echo "   Then run: ./seeder.exe"
    fi
    echo ""
    echo "✅ Workflow complete!"
    echo ""
}

# Main menu loop
while true; do
    show_menu
    read -p "Select option (1-7) or 'q' to quit: " choice
    
    case $choice in
        1) view_migrations ;;
        2) apply_migrations ;;
        3) rollback_migrations ;;
        4) force_version ;;
        5) drop_database ;;
        6) show_status ;;
        7) complete_workflow ;;
        q|Q) 
            echo "Exiting migration guide"
            exit 0
            ;;
        *)
            echo "Invalid option. Please select 1-7 or q."
            echo ""
            ;;
    esac
done
