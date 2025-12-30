package main

import (
	"arabella-api/internal/database"
	"arabella-api/internal/database/seeders"
	"arabella-api/internal/platform/config"
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

func main() {
	// Cargar variables de entorno desde .env
	_ = godotenv.Load()

	rootCmd := &cobra.Command{
		Use:   "console",
		Short: "Database console tools for <no value>",
		Long:  `Command-line tools for managing database migrations and seeders`,
	}

	// migrate command
	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
		Long:  `Automatically migrate database schema based on your models`,
		Run:   runMigrate,
	}
	migrateCmd.Flags().Bool("seed", false, "Run seeders after migration")
	migrateCmd.Flags().Bool("fresh", false, "Drop all tables and migrate from scratch")

	// seed command
	seedCmd := &cobra.Command{
		Use:   "seed",
		Short: "Run database seeders",
		Long:  `Populate database with initial or test data`,
		Run:   runSeed,
	}

	rootCmd.AddCommand(migrateCmd, seedCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func runMigrate(cmd *cobra.Command, args []string) {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("❌ Error connecting to database: %v", err)
	}
	defer database.CloseDB()

	// Check for fresh flag
	fresh, _ := cmd.Flags().GetBool("fresh")
	if fresh {
		log.Println("🗑️  Dropping all tables...")
		if err := db.Migrator().DropTable(database.AllModels...); err != nil {
			log.Printf("⚠️  Warning dropping tables: %v", err)
		}
	}

	// Run migrations
	log.Println("🔄 Running migrations...")
	if err := db.AutoMigrate(database.AllModels...); err != nil {
		log.Fatalf("❌ Migration error: %v", err)
	}
	log.Println("✅ Migrations completed successfully")

	// Check for seed flag
	withSeed, _ := cmd.Flags().GetBool("seed")
	if withSeed {
		runSeedLogic(db)
	}
}

func runSeed(cmd *cobra.Command, args []string) {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("❌ Error connecting to database: %v", err)
	}
	defer database.CloseDB()

	runSeedLogic(db)
}

func runSeedLogic(db *gorm.DB) {
	seeder := &seeders.DatabaseSeeder{}
	if err := seeder.Run(db); err != nil {
		log.Fatalf("❌ Seeder error: %v", err)
	}
}
