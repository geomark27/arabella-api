# arabella-api - Makefile

# Cargar variables de entorno desde .env
ifneq (,$(wildcard .env))
	include .env
	export
endif

.PHONY: build run test clean fmt vet deps help

# Variables
APP_NAME=arabella-api
BUILD_DIR=build
CMD_DIR=cmd/$(APP_NAME)
BRANCH := $(shell git branch --show-current)
COMPOSE_DEV=docker-compose -f docker-compose.dev.yml
COMPOSE_PROD=docker-compose -f docker-compose.yml

# Comandos principales
help: ## Muestra esta ayuda
	@echo "📋 Comandos disponibles:"
	@echo ""
	@echo "  🔨 Build & Run:"
	@echo "    build | run | dev"
	@echo ""
	@echo "  🐳 Desarrollo (solo DB):"
	@echo "    db-up [fresh=1] [location=1] - Inicia PostgreSQL + opciones"
	@echo "    db-fresh                     - db-up + reset automático"
	@echo "    db-fresh-full                - db-up + reset + locations"
	@echo "    db-down | db-restart | db-logs | db-clean | db-shell"
	@echo ""
	@echo "  🐳 Producción (API + DB):"
	@echo "    up | down | restart | logs | logs-api | rebuild"
	@echo ""
	@echo "  🗃️  Database:"
	@echo "    db-migrate | db-seed | fresh"
	@echo ""
	@echo "  🧪 Testing:"
	@echo "    test | test-coverage | fmt | vet | lint"
	@echo ""
	@echo "  📦 Git ($(BRANCH)):"
	@echo "    push m='msg' | pull | status | sync m='msg'"
	@echo ""
	@echo "  🧹 Utils:"
	@echo "    clean | deps | install-tools"
	@echo ""
	@echo "  💡 Ejemplos:"
	@echo "    make db-up                    → Solo DB"
	@echo "    make db-up fresh=1            → DB + reset"
	@echo "    make db-up fresh=1 location=1 → DB + reset + data"
	@echo "    make db-fresh                 → Atajo rápido"
	@echo ""

build: ## Compila la aplicación
	@echo "🔨 Compilando $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)/main.go
	@echo "✅ Compilación exitosa: $(BUILD_DIR)/$(APP_NAME)"

run: ## Ejecuta la aplicación
	@echo "🚀 Ejecutando $(APP_NAME)..."
	@go run $(CMD_DIR)/main.go

test: ## Ejecuta los tests
	@echo "🧪 Ejecutando tests..."
	@go test -v ./...

test-coverage: ## Ejecuta tests con cobertura
	@echo "🧪 Ejecutando tests con cobertura..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Reporte de cobertura generado: coverage.html"

fmt: ## Formatea el código
	@echo "🎨 Formateando código..."
	@go fmt ./...

vet: ## Ejecuta go vet
	@echo "🔍 Analizando código..."
	@go vet ./...

lint: ## Ejecuta golangci-lint (requiere instalación)
	@echo "🔍 Ejecutando linter..."
	@golangci-lint run

deps: ## Descarga las dependencias
	@echo "📦 Descargando dependencias..."
	@go mod download
	@go mod tidy

clean: ## Limpia archivos generados
	@echo "🧹 Limpiando archivos generados..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@go clean

dev: ## Modo desarrollo (requiere air para hot reload)
	@echo "🔥 Iniciando en modo desarrollo..."
	@air
	
dev-full: ## Setup completo desarrollo (DB + migrate + seed + run)
	@echo "🚀 Starting development environment..."
	@$(MAKE) db-up
	@echo "⏳ Waiting for PostgreSQL..."
	@sleep 3
	@loom db:migrate --seed
	@echo "✅ Ready! Starting API..."
	@go run $(CMD_DIR)/main.go

fresh: ## Reset completo (clean DB + migrate + seed)
	@echo "🔄 Fresh install..."
	@$(MAKE) db-clean
	@$(MAKE) db-up
	@echo "⏳ Waiting for PostgreSQL..."
	@sleep 3
	@loom db:fresh --seed
	@if [ "$(location)" = "1" ]; then \
		echo "🌍 Poblando ubicaciones..."; \
		$(MAKE) db-location; \
	fi
	@echo "✅ Database fresh and seeded!"

install-tools: ## Instala herramientas de desarrollo
	@echo "🛠️  Instalando herramientas de desarrollo..."
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# ============================================
# COMANDOS DOCKER DESARROLLO (solo PostgreSQL)
# ============================================

db-up: ## Inicia PostgreSQL [usa fresh=1 para reset completo]
	@echo "🐳 Starting PostgreSQL (DEV MODE)..."
	@$(COMPOSE_DEV) up -d postgres
	@echo "✅ PostgreSQL running on localhost:$(DB_PORT)"
	@if [ "$(fresh)" = "1" ]; then \
		echo ""; \
		echo "🔄 Fresh flag detected! Running database reset..."; \
		echo "⏳ Waiting for PostgreSQL to be ready..."; \
		sleep 3; \
		loom db:fresh --seed; \
		if [ "$(location)" = "1" ]; then \
			echo "🌍 Poblando ubicaciones..."; \
			$(MAKE) db-location; \
		fi; \
		echo "✅ Database fresh and seeded!"; \
	else \
		echo "💡 TIP: Corre tu API con 'make run' o 'make dev'"; \
		echo "💡 Para reset completo usa: make db-up fresh=1"; \
	fi

db-down: ## Detiene PostgreSQL
	@echo "🛑 Stopping PostgreSQL..."
	@$(COMPOSE_DEV) stop postgres

db-restart: ## Reinicia PostgreSQL
	@echo "🔄 Restarting PostgreSQL..."
	@$(COMPOSE_DEV) restart postgres

db-logs: ## Muestra logs de PostgreSQL
	@$(COMPOSE_DEV) logs -f postgres

db-clean: ## Elimina PostgreSQL y volúmenes
	@echo "🧹 Cleaning database..."
	@$(COMPOSE_DEV) down -v
	@echo "✅ Database cleaned"

db-shell: ## Accede a psql en el contenedor
	@$(COMPOSE_DEV) exec postgres psql -U $(DB_USER) -d $(DB_NAME)

db-fresh: ## Alias: db-up con fresh automático
	@$(MAKE) db-up fresh=1

db-fresh-full: ## Alias: db-up + fresh + locations
	@$(MAKE) db-up fresh=1 location=1

# ============================================
# COMANDOS DOCKER PRODUCCIÓN (API + DB)
# ============================================

up: ## Levanta toda la aplicación (API + PostgreSQL) - PRODUCCIÓN
	@echo "🚀 Starting Arabella API (PRODUCTION MODE)..."
	@$(COMPOSE_PROD) up -d
	@echo "✅ API running on http://localhost:$(PORT)"

down: ## Detiene toda la aplicación
	@echo "🛑 Stopping Arabella..."
	@$(COMPOSE_PROD) down

restart: ## Reinicia toda la aplicación
	@echo "🔄 Restarting Arabella..."
	@$(COMPOSE_PROD) restart

logs: ## Muestra logs de todos los servicios
	@$(COMPOSE_PROD) logs -f

logs-api: ## Muestra logs solo de la API
	@$(COMPOSE_PROD) logs -f app

rebuild: ## Reconstruye y levanta la API
	@echo "🔨 Rebuilding Arabella API..."
	@$(COMPOSE_PROD) build --no-cache app
	@$(COMPOSE_PROD) up -d app
	@echo "✅ API rebuilt and running!"

# ============================================
# COMANDOS GIT
# ============================================

push:
	@if [ -z "$(m)" ]; then \
		echo "❌ Error: Debes proporcionar un mensaje"; \
		echo "   Uso: make push m='tu mensaje de commit'"; \
		exit 1; \
	fi
	@echo "📦 Agregando archivos..."
	@git add .
	@echo "✏️  Commiteando: $(m)"
	@git commit -m "$(m)"
	@echo "🚀 Pusheando a origin/$(BRANCH)..."
	@git push origin $(BRANCH)
	@echo "✅ Push completado exitosamente!"

pull:
	@echo "⬇️  Pulling desde origin/$(BRANCH)..."
	@git fetch origin
	@git pull origin $(BRANCH)
	@echo "✅ Pull completado!"

status:
	@echo "📊 Estado de Git (rama: $(BRANCH)):"
	@echo ""
	@git status

sync:
	@if [ -z "$(m)" ]; then \
		echo "❌ Error: Debes proporcionar un mensaje"; \
		echo "   Uso: make sync m='tu mensaje de commit'"; \
		exit 1; \
	fi
	@echo "⬇️  Pulling cambios..."
	@git pull origin $(BRANCH)
	@echo "📦 Agregando archivos..."
	@git add .
	@echo "✏️  Commiteando: $(m)"
	@git commit -m "$(m)"
	@echo "🚀 Pusheando a origin/$(BRANCH)..."
	@git push origin $(BRANCH)
	@echo "✅ Sincronización completada!"

db-migrate: ## Ejecuta migraciones con LOOM
	@echo "🗃️  Running migrations..."
	@loom db:migrate

db-seed: ## Ejecuta seeders con LOOM
	@echo "🌱 Running seeders..."
	@loom db:seed