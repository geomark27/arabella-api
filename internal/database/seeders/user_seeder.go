package seeders

import (
	"log"

	"arabella-api/internal/app/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserSeeder seeds initial users for Arabella Financial OS
type UserSeeder struct{}

// Run executes the user seeder
func (s *UserSeeder) Run(db *gorm.DB) error {

	// Crear usuario demo
	var demoUser models.User
	result := db.Where("email = ?", "demo@arabella.app").First(&demoUser)

	if result.Error == gorm.ErrRecordNotFound {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("demo123"), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("❌ Error hashing demo user password: %v", err)
			return err
		}

		demoUser = models.User{
			UserName:        "demo",
			Email:           "demo@arabella.app",
			PasswordHash:    string(hashedPassword),
			FirstName:       "Demo",
			LastName:        "User",
			DefaultCurrency: "USD",
			EmailVerified:   true,
			IsActive:        true,
			IsSuperAdmin:    false,
		}

		if err := db.Create(&demoUser).Error; err != nil {
			log.Printf("❌ Error creating demo user: %v", err)
			return err
		}
		
	} else if result.Error != nil {
		log.Printf("❌ Error checking demo user: %v", result.Error)
		return result.Error
	} else {
		log.Println("⏭️  Demo User already exists, skipping")
	}

	return nil
}
