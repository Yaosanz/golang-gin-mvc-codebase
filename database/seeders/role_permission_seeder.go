package seeders

import (
	"go-starter-app/app/models"
	"go-starter-app/pkg/database/seeder"

	"gorm.io/gorm"
)

func init() {
	seeder.Register("role_permissions_seeder", RolePermissionSeeder)
}

func RolePermissionSeeder(db *gorm.DB) error {

	// =====================================================
	// ADMIN ROLE - Full System Access
	// =====================================================
	var adminRole models.Role
	if err := db.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		return err
	}

	// Admin gets ALL permissions
	adminPermissions := []string{
		// User Management - Full Access
		"user:read", "user:read:own", "user:create", "user:update",
		"user:update:own", "user:delete", "user:manage",

		// ShortenLink Management - Full Access
		"shortenlink:read", "shortenlink:read:own", "shortenlink:create",
		"shortenlink:update", "shortenlink:update:own", "shortenlink:delete",
		"shortenlink:delete:own", "shortenlink:manage",

		// System Management - Full Access
		"system:read", "system:manage",
	}

	for _, permName := range adminPermissions {
		var perm models.Permission
		if err := db.Where("name = ?", permName).First(&perm).Error; err != nil {
			continue // Skip if permission doesn't exist
		}

		rp := models.RolePermission{
			RoleID:       adminRole.ID,
			PermissionID: perm.ID,
		}

		if err := db.Where("role_id = ? AND permission_id = ?", rp.RoleID, rp.PermissionID).
			FirstOrCreate(&rp).Error; err != nil {
			return err
		}
	}

	// =====================================================
	// USER ROLE - Limited Self-Service Access
	// =====================================================
	var userRole models.Role
	if err := db.Where("name = ?", "user").First(&userRole).Error; err != nil {
		return err
	}

	// Regular users get limited permissions focused on self-service
	userPermissions := []string{
		// User Management - Own Profile Only
		"user:read:own", "user:update:own",

		// ShortenLink Management - Own Resources
		"shortenlink:read:own", "shortenlink:create",
		"shortenlink:update:own", "shortenlink:delete:own",

		// System - Basic Read Access
		"system:read",
	}

	for _, permName := range userPermissions {
		var perm models.Permission
		if err := db.Where("name = ?", permName).First(&perm).Error; err != nil {
			continue // Skip if permission doesn't exist
		}

		rp := models.RolePermission{
			RoleID:       userRole.ID,
			PermissionID: perm.ID,
		}

		if err := db.Where("role_id = ? AND permission_id = ?", rp.RoleID, rp.PermissionID).
			FirstOrCreate(&rp).Error; err != nil {
			return err
		}
	}

	// =====================================================
	// CMS ROLE - Content Management Focus
	// =====================================================
	var cmsRole models.Role
	if err := db.Where("name = ?", "cms").First(&cmsRole).Error; err != nil {
		return err
	}

	// CMS users focus on content management with some user visibility
	cmsPermissions := []string{
		// User Management - Read-Only Access (for content attribution)
		"user:read", "user:read:own", "user:update:own",

		// ShortenLink Management - Full Content Control
		"shortenlink:read", "shortenlink:read:own", "shortenlink:create",
		"shortenlink:update", "shortenlink:update:own", "shortenlink:delete",
		"shortenlink:delete:own", "shortenlink:manage",

		// System - Basic Read Access
		"system:read",
	}

	for _, permName := range cmsPermissions {
		var perm models.Permission
		if err := db.Where("name = ?", permName).First(&perm).Error; err != nil {
			continue // Skip if permission doesn't exist
		}

		rp := models.RolePermission{
			RoleID:       cmsRole.ID,
			PermissionID: perm.ID,
		}

		if err := db.Where("role_id = ? AND permission_id = ?", rp.RoleID, rp.PermissionID).
			FirstOrCreate(&rp).Error; err != nil {
			return err
		}
	}

	return nil
}
