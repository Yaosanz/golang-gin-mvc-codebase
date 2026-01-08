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

	// Ambil semua permission
	var permissions []models.Permission
	if err := db.Find(&permissions).Error; err != nil {
		return err
	}

	// Ambil role admin
	var adminRole models.Role
	if err := db.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		return err
	}

	// Assign semua permission ke ADMIN
	for _, permission := range permissions {
		rp := models.RolePermission{
			RoleID:       adminRole.ID,
			PermissionID: permission.ID,
		}

		if err := db.
			Where(
				"role_id = ? AND permission_id = ?",
				rp.RoleID,
				rp.PermissionID,
			).
			FirstOrCreate(&rp).Error; err != nil {
			return err
		}
	}

	return nil
}
