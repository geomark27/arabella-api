package seeders

import (
	"arabella-api/internal/app/models"
	"log"

	"gorm.io/gorm"
)

// SystemValueSeeder seeds system catalog values
type SystemValueSeeder struct{}

// Run executes the system value seeder
func (s *SystemValueSeeder) Run(db *gorm.DB) error {
	log.Println("🔄 Seeding system values...")

	systemValues := []models.SystemValue{
		// ACCOUNT TYPES
		{CatalogType: "ACCOUNT_TYPE", Value: "BANK", Label: "Cuenta Bancaria", Description: strPtr("Cuenta en banco o institución financiera"), DisplayOrder: 1, IsActive: true},
		{CatalogType: "ACCOUNT_TYPE", Value: "CASH", Label: "Efectivo", Description: strPtr("Dinero en efectivo"), DisplayOrder: 2, IsActive: true},
		{CatalogType: "ACCOUNT_TYPE", Value: "CREDIT_CARD", Label: "Tarjeta de Crédito", Description: strPtr("Línea de crédito"), DisplayOrder: 3, IsActive: true},
		{CatalogType: "ACCOUNT_TYPE", Value: "SAVINGS", Label: "Ahorro", Description: strPtr("Cuenta de ahorros"), DisplayOrder: 4, IsActive: true},
		{CatalogType: "ACCOUNT_TYPE", Value: "INVESTMENT", Label: "Inversión", Description: strPtr("Cuenta de inversión"), DisplayOrder: 5, IsActive: true},
		{CatalogType: "ACCOUNT_TYPE", Value: "CATEGORY", Label: "Categoría", Description: strPtr("Cuenta nominal para categorización"), DisplayOrder: 6, IsActive: true},

		// ACCOUNT CLASSIFICATION
		{CatalogType: "ACCOUNT_CLASSIFICATION", Value: "ASSET", Label: "Activo", Description: strPtr("Recursos que generan valor"), DisplayOrder: 1, IsActive: true},
		{CatalogType: "ACCOUNT_CLASSIFICATION", Value: "LIABILITY", Label: "Pasivo", Description: strPtr("Obligaciones y deudas"), DisplayOrder: 2, IsActive: true},
		{CatalogType: "ACCOUNT_CLASSIFICATION", Value: "EQUITY", Label: "Patrimonio", Description: strPtr("Capital propio"), DisplayOrder: 3, IsActive: true},
		{CatalogType: "ACCOUNT_CLASSIFICATION", Value: "INCOME", Label: "Ingreso", Description: strPtr("Entrada de dinero"), DisplayOrder: 4, IsActive: true},
		{CatalogType: "ACCOUNT_CLASSIFICATION", Value: "EXPENSE", Label: "Gasto", Description: strPtr("Salida de dinero"), DisplayOrder: 5, IsActive: true},

		// TRANSACTION TYPES
		{CatalogType: "TRANSACTION_TYPE", Value: "INCOME", Label: "Ingreso", Description: strPtr("Entrada de dinero"), DisplayOrder: 1, IsActive: true},
		{CatalogType: "TRANSACTION_TYPE", Value: "EXPENSE", Label: "Gasto", Description: strPtr("Salida de dinero"), DisplayOrder: 2, IsActive: true},
		{CatalogType: "TRANSACTION_TYPE", Value: "TRANSFER", Label: "Transferencia", Description: strPtr("Movimiento entre cuentas"), DisplayOrder: 3, IsActive: true},
		{CatalogType: "TRANSACTION_TYPE", Value: "DEBT_PAYMENT", Label: "Pago de Deuda", Description: strPtr("Pago de tarjeta de crédito o préstamo"), DisplayOrder: 4, IsActive: true},

		// CATEGORY TYPES
		{CatalogType: "CATEGORY_TYPE", Value: "INCOME", Label: "Ingreso", Description: strPtr("Categoría de ingreso"), DisplayOrder: 1, IsActive: true},
		{CatalogType: "CATEGORY_TYPE", Value: "EXPENSE", Label: "Gasto", Description: strPtr("Categoría de gasto"), DisplayOrder: 2, IsActive: true},

		// JOURNAL ENTRY TYPES
		{CatalogType: "JOURNAL_ENTRY_TYPE", Value: "DEBIT", Label: "Débito", Description: strPtr("Debe"), DisplayOrder: 1, IsActive: true},
		{CatalogType: "JOURNAL_ENTRY_TYPE", Value: "CREDIT", Label: "Crédito", Description: strPtr("Haber"), DisplayOrder: 2, IsActive: true},
	}

	for _, sv := range systemValues {
		var existing models.SystemValue
		result := db.Where("catalog_type = ? AND value = ?", sv.CatalogType, sv.Value).First(&existing)

		if result.Error == gorm.ErrRecordNotFound {
			if err := db.Create(&sv).Error; err != nil {
				log.Printf("❌ Error creating system value %s.%s: %v", sv.CatalogType, sv.Value, err)
				return err
			}
			log.Printf("   ✅ Created: %s.%s = %s", sv.CatalogType, sv.Value, sv.Label)
		} else if result.Error != nil {
			log.Printf("❌ Error checking system value %s.%s: %v", sv.CatalogType, sv.Value, result.Error)
			return result.Error
		}
		// If exists, skip silently
	}

	log.Println("✅ System values seeded successfully")
	return nil
}

// Helper function to create string pointers
func strPtr(s string) *string {
	return &s
}
