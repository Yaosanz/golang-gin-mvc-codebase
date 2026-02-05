# Database Migrations Guide

## Overview

This document covers all database migrations for the Go Gin MVC application. We use `golang-migrate` (v4) for schema management with PostgreSQL on Supabase.

**Migration Tool:** `migration.exe` (built from `cmd/migration/main.go`)
**Database:** Supabase PostgreSQL
**Location:** `database/migrations/` (SQL files)

## Quick Start

```bash
# Run all pending migrations (create schema)
./migration.exe up

# Rollback all migrations
./migration.exe down

# Force to specific version
./migration.exe force 5

# Drop entire schema
./migration.exe drop
```

## Migrations List

### 000001_extension

- **Description:** Enable uuid-ossp PostgreSQL extension
- **Purpose:** Generate UUID values for id fields
- **Tables Affected:** All tables using uuid type
- **Status:** Foundation migration, must run first

### 000002_create_roles_table

- **Description:** Create roles table for RBAC
- **Columns:**
  - `id` (uuid, primary key)
  - `name` (varchar, unique): "admin", "cms", "user"
  - `description` (text)
  - `created_at` (timestamp)
  - `updated_at` (timestamp)

### 000003_create_permissions_table

- **Description:** Create permissions table
- **Columns:**
  - `id` (uuid, primary key)
  - `name` (varchar, unique): Action permissions
  - `description` (text)
  - `created_at` (timestamp)
  - `updated_at` (timestamp)

### 000004_create_users_table

- **Description:** Create users table
- **Columns:**
  - `id` (uuid, primary key)
  - `email` (varchar, unique)
  - `password` (varchar)
  - `username` (varchar, unique)
  - `first_name` (varchar)
  - `last_name` (varchar)
  - `is_active` (boolean, default: true)
  - `created_at` (timestamp)
  - `updated_at` (timestamp)
- **Indexes:** email, username (unique)

### 000005_create_user_roles_table

- **Description:** Many-to-many junction for users & roles
- **Columns:**
  - `id` (uuid, primary key)
  - `user_id` (uuid, foreign key → users.id)
  - `role_id` (uuid, foreign key → roles.id)
  - `created_at` (timestamp)
- **Constraint:** Unique (user_id, role_id)

### 000006_create_role_permissions_table

- **Description:** Many-to-many junction for roles & permissions
- **Columns:**
  - `id` (uuid, primary key)
  - `role_id` (uuid, foreign key → roles.id)
  - `permission_id` (uuid, foreign key → permissions.id)
  - `created_at` (timestamp)
- **Constraint:** Unique (role_id, permission_id)

### 000007_create_shorten_links_table

- **Description:** Create shortened links table
- **Columns:**
  - `id` (uuid, primary key)
  - `user_id` (uuid, foreign key → users.id)
  - `original_url` (text)
  - `short_code` (varchar, unique)
  - `title` (varchar)
  - `description` (text)
  - `clicks` (integer, default: 0)
  - `is_active` (boolean, default: true)
  - `expires_at` (timestamp, nullable)
  - `created_at` (timestamp)
  - `updated_at` (timestamp)
- **Indexes:** short_code (unique), user_id, is_active

### 000008_add_analytics_columns

- **Description:** Add analytics tracking to shorten_links
- **New Columns:**
  - `last_clicked_at` (timestamp, nullable)
  - `country` (varchar, nullable)
  - `device_type` (varchar, nullable): "mobile", "desktop", "tablet"

### 000009_create_audit_logs_table

- **Description:** Create audit logging table
- **Columns:**
  - `id` (uuid, primary key)
  - `user_id` (uuid, foreign key → users.id, nullable)
  - `action` (varchar)
  - `resource_type` (varchar)
  - `resource_id` (varchar)
  - `changes` (jsonb)
  - `ip_address` (varchar)
  - `user_agent` (text)
  - `created_at` (timestamp)

### 000010_add_metadata_to_users

- **Description:** Add metadata storage to users
- **New Columns:**
  - `metadata` (jsonb, default: '{}')
  - `last_login_at` (timestamp, nullable)
  - `login_count` (integer, default: 0)

### 000011_create_api_keys_table

- **Description:** Create API keys for user authentication
- **Columns:**
  - `id` (uuid, primary key)
  - `user_id` (uuid, foreign key → users.id)
  - `key` (varchar, unique, hashed)
  - `name` (varchar)
  - `last_used_at` (timestamp, nullable)
  - `is_active` (boolean, default: true)
  - `created_at` (timestamp)
  - `updated_at` (timestamp)
- **Indexes:** key (unique), user_id, is_active

### 000012_add_refresh_tokens_table

- **Description:** Create refresh tokens for JWT auth
- **Columns:**
  - `id` (uuid, primary key)
  - `user_id` (uuid, foreign key → users.id)
  - `token` (text, unique)
  - `expires_at` (timestamp)
  - `created_at` (timestamp)
- **Indexes:** user_id, token (unique)

### 000013_add_cms_to_user_role_enum

- **Description:** Ensure all role types are properly supported
- **Changes:** Validates role enum values
- **Supported Roles:** admin, cms, user

## Migration Commands Usage

### Apply All Pending Migrations

```bash
cd backend
./migration.exe up
```

Expected output:

```
Migrating up...
Version 1 migration applied
Version 2 migration applied
...
Version 13 migration applied
Migrations complete
```

### Rollback All Migrations

```bash
./migration.exe down
```

⚠️ **WARNING:** This drops all tables and data. Use only in development!

### Force to Specific Version

```bash
./migration.exe force 8
```

Sets migration version to 8 without running migrations. Use for emergency recovery only.

### Drop Entire Database

```bash
./migration.exe drop
```

⚠️ **DESTRUCTIVE:** Removes all tables, constraints, and extensions.

## Database Schema Summary

### User Management

- **Users:** 10+ attributes per user
- **Roles:** 3 types (admin, cms, user)
- **Permissions:** Action-based access control
- **User_Roles:** Many-to-many relationship
- **Role_Permissions:** Many-to-many relationship

### Link Shortening

- **Shorten_Links:** Core shortened link data
- **Analytics:** Clicks, device type, country tracking
- **Expiration:** Optional URL expiration support

### Security & Audit

- **API_Keys:** Alternative authentication method
- **Refresh_Tokens:** JWT token management
- **Audit_Logs:** Complete action tracking

### Infrastructure

- **UUID Extension:** Used for all id columns
- **Constraints:** Foreign keys, unique indexes
- **Timestamps:** created_at and updated_at on all tables

## Running Migrations in Different Environments

### Local Development (Windows with IPv6 issue)

```bash
# If database connection fails due to DNS:
# 1. Check .env file has correct Supabase credentials
# 2. Verify DB_SSL_MODE=require is set
# 3. If still fails, proceed to Railway deployment

cd backend
./migration.exe up
```

### Production (Railway or other IPv4 host)

```bash
# Railway environment automatically has IPv4 support
# Set environment variables in Railway dashboard:
DB_HOST=db.twqpgihehuuzwghvhsft.supabase.co
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=ihGQWWp3DebgCrpB
DB_NAME=postgres
DB_SSL_MODE=require

# Run migrations
./migration.exe up
```

### Docker Deployment

```dockerfile
# In Dockerfile after building app binary:
COPY backend/database/migrations /app/database/migrations
COPY backend/migration.exe /app/

# Run migrations on container startup:
CMD ./migration.exe up && ./app.exe
```

## Troubleshooting

### Migration fails with "connection refused"

**Causes:**

1. Supabase credentials incorrect
2. IPv6 DNS issue on Windows (known limitation)
3. Database user doesn't have schema permissions

**Solutions:**

- Verify credentials in .env: `DB_HOST`, `DB_USERNAME`, `DB_PASSWORD`
- Test manually: `psql -U postgres -h db.twqpgihehuuzwghvhsft.supabase.co -d postgres`
- Deploy to Railway for IPv4 support
- Ensure database user is "postgres" (has superuser permissions)

### Migration fails with "table already exists"

**Causes:**

- Migrations run multiple times
- Manual table creation before migrations

**Solutions:**

```bash
# Check migration version:
./migration.exe force 0

# Re-run all migrations:
./migration.exe up
```

### Foreign key constraint errors

**Causes:**

- Migrations running out of order
- Missing parent table

**Solutions:**

- Always run `up` command in order
- Don't manually insert data before migration 4 (users table)
- Use seeder.exe for sample data

### "Unknown migration version" error

**Causes:**

- Migration file corruption
- Incomplete migration files

**Solutions:**

```bash
# List available migrations:
ls database/migrations/

# Verify file format: [version]_[name].up.sql format
```

## Creating New Migrations

To add a new migration:

1. Create SQL files:

```bash
# Create migration files
cd backend/database/migrations
# Create new files with next version number
# Example: 000014_create_new_table.up.sql
#          000014_create_new_table.down.sql
```

2. Write SQL:

```sql
-- UP migration (000014_create_new_table.up.sql)
CREATE TABLE new_table (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamp DEFAULT CURRENT_TIMESTAMP
);

-- DOWN migration (000014_create_new_table.down.sql)
DROP TABLE IF EXISTS new_table;
```

3. Run and test:

```bash
./migration.exe up     # Apply
./migration.exe down   # Test rollback
./migration.exe up     # Reapply
```

## Migration File Format

Each migration must have two files:

```
000001_name.up.sql     (apply migration)
000001_name.down.sql   (rollback migration)
```

**Naming Rules:**

- Version: 6 digits, zero-padded (000001, 000002, etc.)
- Separator: Underscore
- Description: Lowercase, underscores for spaces
- Extension: .up.sql or .down.sql

**Example:**

```
000014_add_user_preferences.up.sql
000014_add_user_preferences.down.sql
```

## Testing Migrations

### Manual Test Flow

```bash
# 1. Start fresh
./migration.exe drop

# 2. Apply migrations
./migration.exe up

# 3. Run seeders
../seeder.exe

# 4. Verify data
psql -U postgres -h db.twqpgihehuuzwghvhsft.supabase.co -d postgres \
  -c "SELECT count(*) FROM users;"

# 5. Test rollback
./migration.exe down

# 6. Verify tables gone
./migration.exe up
```

### Continuous Integration

In your CI/CD pipeline (GitHub Actions, etc.):

```yaml
- name: Run database migrations
  env:
    DB_HOST: ${{ secrets.DB_HOST }}
    DB_PORT: ${{ secrets.DB_PORT }}
    DB_USERNAME: ${{ secrets.DB_USERNAME }}
    DB_PASSWORD: ${{ secrets.DB_PASSWORD }}
  run: |
    cd backend
    go build -o migration.exe ./cmd/migration
    ./migration.exe up
```

## Version Control

All migration files should be committed to git:

```bash
# Commit new migrations
git add database/migrations/
git commit -m "feat: add new table schema migration"
```

**DO NOT:**

- Modify existing migration files (only add new ones)
- Skip migration versions
- Run migrations out of order

## Related Files

- **Driver:** `database/connection.go` (connection pool, GORM setup)
- **Config:** `.env` (database credentials)
- **CLI:** `cmd/migration/main.go` (migration command implementation)
- **Bootstrap:** `bootstrap/app.go` (app initialization, graceful DB handling)
- **Seeders:** See `SEEDERS.md` for sample data injection

## Next Steps

1. ✅ Review migration list above
2. ✅ Verify all 13 migration files exist in `database/migrations/`
3. ▶️ Run `./migration.exe up` in appropriate environment
4. ▶️ Run `../seeder.exe` to populate sample data
5. ▶️ Test API endpoints with real database

---

**Last Updated:** Integration with Supabase completed  
**Status:** All 13 migrations validated and ready for execution  
**Next Task:** Execute migrations in production/deployment environment
