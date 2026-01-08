# TODO: Fix Migration 000009 and Adjust Related Files

## Completed

- [x] Update migration 000009 up: Drop role column, update unique constraint to (role_id, permission_id)
- [x] Update migration 000009 down: Add back role column, update unique to (role, permission_id)
- [x] Create Role model
- [x] Create role_seeder.go
- [x] Update RolePermission model: Change Role to RoleID UUID, add FK to roles
- [x] Update user_seeder: Set RoleID by querying roles
- [x] Update auth_service Login: Pass user.RoleID to GenerateToken
- [x] Update JWT claims: Add RoleID UUID
- [x] Update GenerateToken: Accept roleID UUID
- [x] Update auth middleware: Set "role_id" to claims.RoleID (UUID)
- [x] Update permission middleware: Use roleID as UUID
- [x] Update permission repo: Use rp.role_id in query, change interface to accept uuid.UUID
- [x] Update permission service: Change to accept uuid.UUID
- [x] Update role_permission_seeder: Set RoleID by querying roles

## Remaining

- [x] Update auth_service Register: Set RoleID for new users
- [x] Run migrations up
- [x] Run seeders (role_seeder first, then others)
- [x] Test authentication and permission checks
