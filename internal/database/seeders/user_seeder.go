package seeders

import (
	"log"

	"arabella-api/internal/app/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserSeeder seeds initial users
type UserSeeder struct{}

// Run executes the user seeder
func (s *UserSeeder) Run(db *gorm.DB) error {

	// 1. Crear SuperAdmin (modo dios - sin empresa)
	var superAdmin models.User
	result := db.Where("email = ?", "superadmin@dvra.com").First(&superAdmin)

	if result.Error == gorm.ErrRecordNotFound {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("SuperAdmin123!"), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("❌ Error hashing superadmin password: %v", err)
			return err
		}

		superAdmin = models.User{
			FirstName:    "Super",
			LastName:     "Admin",
			Email:        "superadmin@dvra.com",
			PasswordHash: string(hashedPassword),
			IsActive:     true,
		}

		if err := db.Create(&superAdmin).Error; err != nil {
			log.Printf("❌ Error creating superadmin user: %v", err)
			return err
		}

		log.Println("✅ SuperAdmin created successfully (email: superadmin@dvra.com, password: SuperAdmin123!)")
	} else if result.Error != nil {
		log.Printf("❌ Error checking superadmin: %v", result.Error)
		return result.Error
	} else {
		log.Println("⏭️  SuperAdmin already exists, skipping")
	}

	// 2. Crear Admin de empresa (para Azentic Sys)
	var companyAdmin models.User
	result = db.Where("email = ?", "admin@azentic.com").First(&companyAdmin)

	if result.Error == gorm.ErrRecordNotFound {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("Admin123!"), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("❌ Error hashing admin password: %v", err)
			return err
		}

		companyAdmin = models.User{
			FirstName:    "Azentic",
			LastName:     "Systems",
			Email:        "admin@azentic.com",
			PasswordHash: string(hashedPassword),
			IsActive:     true,
		}

		if err := db.Create(&companyAdmin).Error; err != nil {
			log.Printf("❌ Error creating company admin user: %v", err)
			return err
		}

		log.Println("✅ Company Admin created successfully (email: admin@azentic.com, password: Admin123!)")
	} else if result.Error != nil {
		log.Printf("❌ Error checking company admin: %v", result.Error)
		return result.Error
	} else {
		log.Println("⏭️  Company Admin already exists, skipping")
	}

	return nil
}
