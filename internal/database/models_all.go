package database

import (
	"arabella-api/internal/app/models"
)

// AllModels contains all models for dynamic migration
// Add your models here to include them in auto-migration
var AllModels = []interface{}{
	&models.User{},
	&models.SystemValue{},
	&models.Currency{},
	&models.Category{},
	&models.Transaction{},
	&models.JournalEntry{},
	&models.Account{},
}
