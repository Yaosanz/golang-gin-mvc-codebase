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
		"user:read",
		"user:create",
		"user:update",
		"user:delete",
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
