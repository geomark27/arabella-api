package main

import (
	"arabella-api/internal/database"
	"arabella-api/internal/platform/config"
	"arabella-api/internal/platform/server"
	"log"

	"github.com/joho/godotenv"
)

// @title           DVRA API
// @version         1.2.0
// @description     API para sistema de reclutamiento y gestión de candidatos
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@dvra.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8001
// @BasePath  /api/v1
// @schemes   http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Cargar variables de entorno desde .env
	_ = godotenv.Load()

	// Cargar configuración
	cfg := config.Load()

	// Inicializar base de datos
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatal("Error conectando a la base de datos:", err)
	}
	defer func() {
		if err := database.CloseDB(); err != nil {
			log.Printf("Error cerrando la base de datos: %v", err)
		}
	}()

	// Crear servidor
	srv := server.New(cfg, db)

	// Mensaje de inicio
	log.Printf("🚀 Servidor %s iniciado en http://localhost:%s", "arabella-api", cfg.Port)
	log.Printf("✨ Proyecto generado con Loom")

	// Iniciar servidor
	if err := srv.Start(); err != nil {
		log.Fatal("Error iniciando servidor:", err)
	}
}
