package seeders

import (
	"go-starter-app/app/models"
	"go-starter-app/helpers"
	"go-starter-app/pkg/database/seeder"
	"log"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func init() {
	seeder.Register("user_seeder", user)
}

func user(db *gorm.DB) error {
	// Define the seed data
	active := true
	hashedPassword, err := helpers.HashPassword("password", 12)
	if err != nil {
		return err
	}

	data := []models.User{
		{
			Name:     "Dewi Sartika",
			Username: "dewisartika",
			Email:    "dewisartika@gmail.com",
			Phone:    nil,
			Password: hashedPassword, // use the hashed password
			IsActive: &active,
		},
		{
			Name:     "Agus Sutanto",
			Username: "agussutanto",
			Email:    "agussutanto@gmail.com",
			Phone:    nil,
			Password: hashedPassword, // use the hashed password
			IsActive: &active,
		},
	}

	for _, item := range data {
		var existing models.User
		err := db.Where("email = ?", item.Email).First(&existing).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Not found, proceed to create
				err = db.Create(&item).Error
				if err != nil {
					log.Printf("[SEEDER]: Failed to seed user %v: %v", item.Name, err)
					continue // continue to next item
				}

				log.Printf("[SEEDER]: Successfully seeded user %v", item.Name)
				continue // continue to next item
			}

			// Unexpected error
			log.Printf("[SEEDER]: Failed to check existing user %v: %v", item.Name, err)
			continue // continue to next item
		}

		// Already exists, skip
		log.Printf("[SEEDER]: user %v already exists (ID %d), skipping", existing.Name, existing.ID)
	}

	return nil
}
