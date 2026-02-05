# Database Seeders Guide

## Overview

Database seeders populate the schema with sample/test data. After running migrations to create tables, run seeders to fill them with initial data.

**Seeder Tool:** `seeder.exe` (built from `cmd/seeder/main.go`)
**Purpose:** Initialize development/test databases with realistic data
**Location:** `cmd/seeder/` (Go seeder implementation)

## Quick Start

```bash
# Run all seeders
./seeder.exe

# Run specific seeder
./seeder.exe --seed=users

# Truncate tables before seeding (refresh)
./seeder.exe --fresh
```

## Available Seeders

### 1. RoleSeeder

**File:** `cmd/seeder/seeders/role_seeder.go`

Creates default system roles:

| Role    | Description          | Purpose                      |
| ------- | -------------------- | ---------------------------- |
| `admin` | System administrator | Full access to all features  |
| `cms`   | Content management   | Manage users and permissions |
| `user`  | Regular user         | Basic link shortening        |

```go
// Example data created:
{
  "id": "uuid",
  "name": "admin",
  "description": "System administrator",
  "created_at": "2024-01-01T00:00:00Z"
}
```

**Output:** 3 roles inserted

### 2. PermissionSeeder

**File:** `cmd/seeder/seeders/permission_seeder.go`

Creates default permission records:

| Permission            | Description              |
| --------------------- | ------------------------ |
| `create:users`        | Create new user accounts |
| `read:users`          | View user information    |
| `update:users`        | Modify user data         |
| `delete:users`        | Remove user accounts     |
| `create:shortenlinks` | Create shortened links   |
| `read:shortenlinks`   | View shortened links     |
| `update:shortenlinks` | Modify shortened links   |
| `delete:shortenlinks` | Remove shortened links   |
| `manage:roles`        | Manage user roles        |
| `manage:permissions`  | Manage permissions       |

**Output:** 10 permissions inserted

### 3. UserSeeder

**File:** `cmd/seeder/seeders/user_seeder.go`

Creates test user accounts:

| Username | Email             | Password | Role  | Purpose                    |
| -------- | ----------------- | -------- | ----- | -------------------------- |
| `admin`  | admin@example.com | admin123 | admin | Full access testing        |
| `cms`    | cms@example.com   | cms123   | cms   | Content management testing |
| `user1`  | user1@example.com | user123  | user  | Basic user testing         |
| `user2`  | user2@example.com | user123  | user  | Second test user           |

**Security Note:** Passwords are hashed using bcrypt in production seeders.

**Output:** 4 users created with appropriate roles

### 4. RolePermissionSeeder

**File:** `cmd/seeder/seeders/role_permission_seeder.go`

Associates permissions with roles:

| Role    | Permissions                                                 |
| ------- | ----------------------------------------------------------- |
| `admin` | All 10 permissions                                          |
| `cms`   | create:users, read:users, update:users, manage:roles        |
| `user`  | create:shortenlinks, read:shortenlinks, update:shortenlinks |

**Output:** 19 role-permission associations created

### 5. ShortenLinkSeeder

**File:** `cmd/seeder/seeders/shorten_link_seeder.go`

Creates sample shortened links:

| Original URL                | Short Code | Title           | Owner | Purpose       |
| --------------------------- | ---------- | --------------- | ----- | ------------- |
| `https://golang.org`        | `go`       | Golang Official | admin | Testing       |
| `https://github.com`        | `gh`       | GitHub          | user1 | Link tracking |
| `https://stackoverflow.com` | `so`       | Stack Overflow  | user1 | Analytics     |

**Features:**

- Tracks click counts (0-100 sample clicks)
- Optional expiration dates
- Device/country analytics

**Output:** 3+ shortened links with analytics data

## Execution Flow

```
1. Verify migrations completed
   └─ Check all 13 migrations applied

2. Run seeders
   └─ RoleSeeder
   └─ PermissionSeeder
   └─ UserSeeder
   └─ RolePermissionSeeder
   └─ ShortenLinkSeeder

3. Verify data
   └─ Query user count
   └─ Check role assignments
   └─ List shortened links
```

## Running Seeders

### Standard Execution

```bash
cd backend

# Build if needed
go build -o seeder.exe ./cmd/seeder

# Run all seeders
./seeder.exe
```

**Expected Output:**

```
🌱 Starting database seeders...

✅ RoleSeeder: 3 roles created
✅ PermissionSeeder: 10 permissions created
✅ UserSeeder: 4 users created
✅ RolePermissionSeeder: 19 associations created
✅ ShortenLinkSeeder: 3 links created

✅ All seeders completed successfully!
Total records: 39
```

### Fresh Database (Truncate First)

```bash
./seeder.exe --fresh
```

Truncates all tables, then re-seeds. Useful for:

- Cleaning up test data
- Resetting auto-increment IDs
- Starting fresh after failed tests

**⚠️ WARNING:** Deletes all existing data in seeded tables!

### Seed Specific Tables

```bash
./seeder.exe --seed=users
./seeder.exe --seed=shortenlinks
./seeder.exe --seed=roles,permissions
```

Runs only specified seeders without running others.

## Integration with Migrations

### Recommended Workflow

```bash
# 1. Drop existing schema (if starting fresh)
./migration.exe drop

# 2. Create schema
./migration.exe up

# 3. Populate with sample data
./seeder.exe
```

### Production Considerations

For production, typically:

1. Run migrations only (skip seeders)
2. Manually create admin user account
3. Load production data from external source

```bash
# Production flow (NO seeders)
./migration.exe up
# Manual admin user creation via API or manual SQL
```

## Sample Data for Testing

### Default Admin User

```json
{
  "email": "admin@example.com",
  "username": "admin",
  "password": "admin123",
  "role": "admin",
  "first_name": "Admin",
  "last_name": "User"
}
```

**Test Flow:**

1. Login with admin@example.com / admin123
2. Create new shortened links
3. View analytics
4. Manage users

### Test API Authentication

```bash
# Get JWT token
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin123"
  }'

# Response includes JWT token
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "uuid",
    "email": "admin@example.com",
    "roles": ["admin"]
  }
}

# Use token for authenticated requests
curl -H "Authorization: Bearer <token>" \
  http://localhost:8080/api/shortenlinks
```

### Testing User Roles

```bash
# Test as admin (full access)
Login: admin@example.com / admin123
Access: All endpoints, all resources

# Test as CMS user (limited access)
Login: cms@example.com / cms123
Access: User management, limited permissions

# Test as regular user (restricted access)
Login: user1@example.com / user123
Access: Own shortened links, basic features
```

## Creating Custom Seeders

To add a new seeder:

### 1. Create Seeder File

```go
// cmd/seeder/seeders/custom_seeder.go
package seeders

import "gorm.io/gorm"

type CustomSeeder struct{}

func (s *CustomSeeder) Run(db *gorm.DB) error {
	// Implement seeding logic
	// Insert data, validate, etc.
	return nil
}

func (s *CustomSeeder) Name() string {
	return "CustomSeeder"
}
```

### 2. Register in Main

```go
// cmd/seeder/main.go
seeders := []Seeder{
	&seeders.RoleSeeder{},
	&seeders.PermissionSeeder{},
	&seeders.UserSeeder{},
	&seeders.RolePermissionSeeder{},
	&seeders.ShortenLinkSeeder{},
	&seeders.CustomSeeder{},  // Add new seeder
}
```

### 3. Implement Logic

```go
func (s *CustomSeeder) Run(db *gorm.DB) error {
	data := []MyModel{
		{Name: "item1"},
		{Name: "item2"},
	}

	for _, item := range data {
		if err := db.Create(&item).Error; err != nil {
			return err
		}
	}

	return nil
}
```

### 4. Build and Test

```bash
go build -o seeder.exe ./cmd/seeder
./seeder.exe --seed=custom
```

## Testing with Seeded Data

### Verify Users Created

```bash
psql -U postgres -h db.twqpgihehuuzwghvhsft.supabase.co -d postgres -c \
  "SELECT email, username, is_active FROM users;"

# Expected output:
# email              | username | is_active
# admin@example.com  | admin    | t
# cms@example.com    | cms      | t
# user1@example.com  | user1    | t
# user2@example.com  | user2    | t
```

### Verify Roles Assigned

```bash
psql -U postgres -h db.twqpgihehuuzwghvhsft.supabase.co -d postgres -c \
  "SELECT u.email, r.name
   FROM users u
   JOIN user_roles ur ON u.id = ur.user_id
   JOIN roles r ON ur.role_id = r.id;"
```

### Verify Permissions Set

```bash
psql -U postgres -h db.twqpgihehuuzwghvhsft.supabase.co -d postgres -c \
  "SELECT r.name, p.name
   FROM roles r
   JOIN role_permissions rp ON r.id = rp.role_id
   JOIN permissions p ON rp.permission_id = p.id
   ORDER BY r.name;"
```

### Verify Links Created

```bash
psql -U postgres -h db.twqpgihehuuzwghvhsft.supabase.co -d postgres -c \
  "SELECT short_code, original_url, clicks FROM shorten_links;"
```

## Troubleshooting

### Seeder fails with "duplicate key value"

**Cause:** Seeder ran multiple times, trying to insert duplicate data

**Solutions:**

```bash
# Option 1: Use fresh flag
./seeder.exe --fresh

# Option 2: Truncate tables manually
psql -U postgres -h db.twqpgihehuuzwghvhsft.supabase.co -d postgres -c \
  "TRUNCATE users, roles, permissions, user_roles,
    role_permissions, shorten_links CASCADE;"
```

### "Table does not exist" error

**Cause:** Migrations not run before seeders

**Solution:**

```bash
./migration.exe up
./seeder.exe
```

### Foreign key constraint violation

**Cause:** Seeder running out of order or missing parent data

**Solution:**

- Ensure all 5 seeders run in correct order
- Don't manually delete parent records
- Use fresh flag to reset

### Password hashing failures

**Cause:** bcrypt library not available

**Solution:**

```bash
go get golang.org/x/crypto/bcrypt
./seeder.exe
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Seed Database

on: [push]

jobs:
  seed:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2

      - uses: actions/setup-go@v2
        with:
          go-version: 1.20

      - name: Run migrations
        env:
          DB_HOST: ${{ secrets.DB_HOST }}
          DB_USERNAME: ${{ secrets.DB_USERNAME }}
          DB_PASSWORD: ${{ secrets.DB_PASSWORD }}
        run: |
          cd backend
          go build -o migration.exe ./cmd/migration
          ./migration.exe up

      - name: Run seeders
        env:
          DB_HOST: ${{ secrets.DB_HOST }}
          DB_USERNAME: ${{ secrets.DB_USERNAME }}
          DB_PASSWORD: ${{ secrets.DB_PASSWORD }}
        run: |
          cd backend
          go build -o seeder.exe ./cmd/seeder
          ./seeder.exe
```

## Data Validation

### Post-Seeding Checklist

```bash
#!/bin/bash
# Verify all seeded data

echo "Checking users..."
USER_COUNT=$(psql -t -c "SELECT COUNT(*) FROM users;")
[[ $USER_COUNT -eq 4 ]] && echo "✅ Users: $USER_COUNT" || echo "❌ Users: $USER_COUNT"

echo "Checking roles..."
ROLE_COUNT=$(psql -t -c "SELECT COUNT(*) FROM roles;")
[[ $ROLE_COUNT -eq 3 ]] && echo "✅ Roles: $ROLE_COUNT" || echo "❌ Roles: $ROLE_COUNT"

echo "Checking permissions..."
PERM_COUNT=$(psql -t -c "SELECT COUNT(*) FROM permissions;")
[[ $PERM_COUNT -eq 10 ]] && echo "✅ Permissions: $PERM_COUNT" || echo "❌ Permissions: $PERM_COUNT"

echo "Checking role-permissions..."
RPPERM_COUNT=$(psql -t -c "SELECT COUNT(*) FROM role_permissions;")
[[ $RPPERM_COUNT -gt 0 ]] && echo "✅ Role-Permissions: $RPPERM_COUNT" || echo "❌ Missing"

echo "Checking links..."
LINK_COUNT=$(psql -t -c "SELECT COUNT(*) FROM shorten_links;")
[[ $LINK_COUNT -gt 0 ]] && echo "✅ Links: $LINK_COUNT" || echo "❌ Missing"
```

## Related Files

- **CLI:** `cmd/seeder/main.go` (seeder command implementation)
- **Seeders:** `cmd/seeder/seeders/` (individual seeder implementations)
- **Models:** `app/models/` (data model definitions)
- **Migrations:** `database/migrations/` (schema creation)
- **Config:** `.env` (database credentials)

## Next Steps

1. ✅ Review seeder list above
2. ✅ Verify seeder files exist in `cmd/seeder/seeders/`
3. ▶️ Run migrations first: `./migration.exe up`
4. ▶️ Run seeders: `./seeder.exe`
5. ▶️ Verify data with queries above
6. ▶️ Test API endpoints with seeded data

---

**Last Updated:** Seeder integration with Supabase completed  
**Status:** 5 seeders ready for execution  
**Next Task:** Execute seeders after migrations complete
