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
	var adminRole, userRole models.Role
	if err := db.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		return err
	}
	if err := db.Where("name = ?", "user").First(&userRole).Error; err != nil {
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
		},
		{
			Name:     "Sandy Budi Wirawan",
			Username: "sandy",
			Email:    "sandy@gmail.com",
			Phone:    nil,
			Password: hashedPassword,
			IsActive: active,
			Role:     models.RoleUser,
			RoleID:   userRole.ID,
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

		log.Printf(
			"[SEEDER][USER]: User %s already exists (role=%s), skipping",
			existing.Email,
			existing.Role,
		)
	}

	return nil
}
