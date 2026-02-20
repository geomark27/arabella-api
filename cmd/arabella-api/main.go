package main

import (
	// Importar los docs generados por swag para registrarlos en el servidor
	_ "arabella-api/docs/swagger"
	"arabella-api/internal/database"
	"arabella-api/internal/platform/config"
	"arabella-api/internal/platform/server"
	"log"

	"github.com/joho/godotenv"
)

// @title           Arabella Financial OS API
// @version         1.0.0
// @description     API REST para el sistema de gestión financiera personal Arabella. Implementa contabilidad de doble partida, cálculo de Runway y soporte multi-moneda.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Arabella Support
// @contact.email  support@arabella.app

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1
// @schemes   http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Escribe "Bearer" seguido de un espacio y tu JWT de acceso. Ejemplo: "Bearer eyJhbGci..."

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
	log.Printf("📚 Documentación API: http://localhost:%s/swagger/index.html", cfg.Port)
	log.Printf("❤️  Health check:      http://localhost:%s/api/v1/health", cfg.Port)
	log.Printf("📊 Dashboard:         http://localhost:%s/api/v1/dashboard", cfg.Port)
	log.Printf("🔐 Auth:              http://localhost:%s/api/v1/auth/login", cfg.Port)
	log.Printf("🌍 Entorno:           %s", cfg.Environment)
	log.Printf("✨ Proyecto generado con Loom")

	// Iniciar servidor
	if err := srv.Start(); err != nil {
		log.Fatal("Error iniciando servidor:", err)
	}
}
