package seeders

import (
	"log"

	"go-starter-app/app/models"
	"go-starter-app/helpers"
	"go-starter-app/pkg/database/seeder"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const defaultPassword = "secret123"

func init() {
	seeder.Register("user_seeder", userSeeder)
}

func userSeeder(db *gorm.DB) error {
	active := true

	hashedPassword, err := helpers.HashPassword(defaultPassword, 12)
	if err != nil {
		return err
	}

	// Get role IDs
	var adminRole, userRole, cmsRole models.Role
	if err := db.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		return err
	}
	if err := db.Where("name = ?", "user").First(&userRole).Error; err != nil {
		return err
	}
	if err := db.Where("name = ?", "cms").First(&cmsRole).Error; err != nil {
		return err
	}

	users := []models.User{
		{
			Name:     "Super Admin",
			Username: "admin",
			Email:    "admin@mydigilearn.com",
			Phone:    nil,
			Password: hashedPassword,
			IsActive: active,
			Role:     models.RoleAdmin,
			RoleID:   adminRole.ID,
			Roles:    []models.Role{adminRole},
		},
		{
			Name:     "Regular User",
			Username: "user",
			Email:    "user@mydigilearn.com",
			Phone:    nil,
			Password: hashedPassword,
			IsActive: active,
			Role:     models.RoleUser,
			RoleID:   userRole.ID,
			Roles:    []models.Role{userRole},
		},
		{
			Name:     "CMS Manager",
			Username: "cms",
			Email:    "cms@mydigilearn.com",
			Phone:    nil,
			Password: hashedPassword,
			IsActive: active,
			Role:     models.RoleCMS,
			RoleID:   cmsRole.ID,
			Roles:    []models.Role{cmsRole},
		},
		{
			Name:     "Multi Role User",
			Username: "multi",
			Email:    "multi@mydigilearn.com",
			Phone:    nil,
			Password: hashedPassword,
			IsActive: active,
			Role:     models.RoleUser, 
			RoleID:   userRole.ID,
			Roles:    []models.Role{userRole, cmsRole}, // Multiple roles
		},
	}

	for _, item := range users {
		var existing models.User

		err := db.
			Where("email = ?", item.Email).
			First(&existing).
			Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := db.Create(&item).Error; err != nil {
					log.Printf(
						"[SEEDER][USER]: Failed to seed %s (%s): %v",
						item.Name,
						item.Email,
						err,
					)
					continue
				}

				log.Printf(
					"[SEEDER][USER]: Seeded %s (%s) role=%s | password=%s",
					item.Name,
					item.Email,
					item.Role,
					defaultPassword,
				)
				continue
			}

			log.Printf(
				"[SEEDER][USER]: Error checking user %s: %v",
				item.Email,
				err,
			)
			continue
		}

		// Update existing user to ensure active
		if !existing.IsActive {
			if err := db.Model(&existing).Update("is_active", true).Error; err != nil {
				log.Printf(
					"[SEEDER][USER]: Failed to activate user %s: %v",
					existing.Email,
					err,
				)
			} else {
				log.Printf(
					"[SEEDER][USER]: Activated existing user %s (role=%s)",
					existing.Email,
					existing.Role,
				)
			}
		} else {
			log.Printf(
				"[SEEDER][USER]: User %s already exists and active (role=%s), skipping",
				existing.Email,
				existing.Role,
			)
		}
	}

	return nil
}
