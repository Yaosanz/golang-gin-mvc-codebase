package seeders

import (
	"go-starter-app/app/models"
	"go-starter-app/pkg/database/seeder"

	"gorm.io/gorm"
)

func init() {
	seeder.Register("role_seeder", RoleSeeder)
}

func RoleSeeder(db *gorm.DB) error {
	roles := []models.Role{
		{Name: "admin"},
		{Name: "user"},
		{Name: "cms"},
	}

	for _, role := range roles {
		if err := db.
			Where("name = ?", role.Name).
			FirstOrCreate(&role).Error; err != nil {
			return err
		}
	}

	return nil
}
