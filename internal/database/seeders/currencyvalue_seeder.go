package seeders

import (
	"arabella-api/internal/app/models"
	"log"

	"gorm.io/gorm"
)

// CurrencySeeder seeds the currencies catalog
type CurrencySeeder struct{}

// Run executes the currency seeder
func (s *CurrencySeeder) Run(db *gorm.DB) error {
	log.Println("🔄 Seeding currencies...")

	currencies := []models.Currency{
		// ── Principales globales ──────────────────────────────────────────
		{Code: "USD", Name: "US Dollar", Symbol: "$", IsActive: true},
		{Code: "EUR", Name: "Euro", Symbol: "€", IsActive: true},
		{Code: "GBP", Name: "British Pound", Symbol: "£", IsActive: true},
		{Code: "JPY", Name: "Japanese Yen", Symbol: "¥", IsActive: true},
		{Code: "CHF", Name: "Swiss Franc", Symbol: "CHF", IsActive: true},
		{Code: "CAD", Name: "Canadian Dollar", Symbol: "$", IsActive: true},
		{Code: "AUD", Name: "Australian Dollar", Symbol: "$", IsActive: true},
		{Code: "CNY", Name: "Chinese Yuan", Symbol: "¥", IsActive: true},

		// ── Latinoamérica ─────────────────────────────────────────────────
		{Code: "MXN", Name: "Mexican Peso", Symbol: "$", IsActive: true},
		{Code: "COP", Name: "Colombian Peso", Symbol: "$", IsActive: true},
		{Code: "ARS", Name: "Argentine Peso", Symbol: "$", IsActive: true},
		{Code: "BRL", Name: "Brazilian Real", Symbol: "R$", IsActive: true},
		{Code: "CLP", Name: "Chilean Peso", Symbol: "$", IsActive: true},
		{Code: "PEN", Name: "Peruvian Sol", Symbol: "S/", IsActive: true},
		{Code: "UYU", Name: "Uruguayan Peso", Symbol: "$", IsActive: true},
		{Code: "PYG", Name: "Paraguayan Guaraní", Symbol: "₲", IsActive: true},
		{Code: "BOB", Name: "Bolivian Boliviano", Symbol: "Bs.", IsActive: true},
		{Code: "VES", Name: "Venezuelan Bolívar", Symbol: "Bs.S", IsActive: true},
		{Code: "GTQ", Name: "Guatemalan Quetzal", Symbol: "Q", IsActive: true},
		{Code: "HNL", Name: "Honduran Lempira", Symbol: "L", IsActive: true},
		{Code: "NIO", Name: "Nicaraguan Córdoba", Symbol: "C$", IsActive: true},
		{Code: "CRC", Name: "Costa Rican Colón", Symbol: "₡", IsActive: true},
		{Code: "PAB", Name: "Panamanian Balboa", Symbol: "B/.", IsActive: true},
		{Code: "DOP", Name: "Dominican Peso", Symbol: "RD$", IsActive: true},
		{Code: "CUP", Name: "Cuban Peso", Symbol: "$", IsActive: true},
		{Code: "TTD", Name: "Trinidad and Tobago Dollar", Symbol: "TT$", IsActive: true},
		{Code: "JMD", Name: "Jamaican Dollar", Symbol: "J$", IsActive: true},
		{Code: "ECU", Name: "Ecuadorian (uses USD)", Symbol: "$", IsActive: false}, // Ecuador usa USD
	}

	created := 0
	skipped := 0

	for _, c := range currencies {
		var existing models.Currency
		result := db.Where("code = ?", c.Code).First(&existing)

		if result.Error == gorm.ErrRecordNotFound {
			if err := db.Create(&c).Error; err != nil {
				log.Printf("❌ Error creating currency %s: %v", c.Code, err)
				return err
			}
			log.Printf("   ✅ Created: %s — %s (%s)", c.Code, c.Name, c.Symbol)
			created++
		} else if result.Error != nil {
			log.Printf("❌ Error checking currency %s: %v", c.Code, result.Error)
			return result.Error
		} else {
			skipped++
		}
	}

	log.Printf("✅ Currencies seeded — created: %d, skipped (already exist): %d", created, skipped)
	return nil
}
