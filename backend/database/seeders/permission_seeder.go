package seeders

import (
	"go-starter-app/app/models"
	"go-starter-app/pkg/database/seeder"

	"gorm.io/gorm"
)

func init() {
	seeder.Register("permission_seeder", PermissionSeeder)
}

func PermissionSeeder(db *gorm.DB) error {
	permissions := []string{
		// User Management Permissions
		"user:read",           // View all users
		"user:read:own",       // View own profile
		"user:create",         // Create new users
		"user:update",         // Update any user
		"user:update:own",     // Update own profile
		"user:delete",         // Delete users
		"user:manage",         // Full user management (admin only)

		// ShortenLink Management Permissions
		"shortenlink:read",           // View all shorten links
		"shortenlink:read:own",       // View own shorten links
		"shortenlink:create",         // Create shorten links
		"shortenlink:update",         // Update any shorten link
		"shortenlink:update:own",     // Update own shorten links
		"shortenlink:delete",         // Delete any shorten link
		"shortenlink:delete:own",     // Delete own shorten links
		"shortenlink:manage",         // Full shorten link management

		// System Permissions
		"system:read",         // View system information
		"system:manage",       // System administration
	}

	for _, name := range permissions {
		if err := db.
			Where("name = ?", name).
			FirstOrCreate(&models.Permission{
				Name: name,
			}).Error; err != nil {
			return err
		}
	}

	return nil
}
