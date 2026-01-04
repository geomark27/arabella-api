# arabella-api - Makefile generado por Loom

.PHONY: build run test clean fmt vet deps help

# Variables
APP_NAME=arabella-api
BUILD_DIR=build
CMD_DIR=cmd/$(APP_NAME)
BRANCH := $(shell git branch --show-current)

# Comandos principales
help: ## Muestra esta ayuda
	@echo "📋 Comandos disponibles:"
	@echo ""
	@echo "  🔨 Compilación y Ejecución:"
	@echo "    make build        - Compila la aplicación"
	@echo "    make run          - Ejecuta la aplicación"
	@echo "    make dev          - Modo desarrollo con hot reload (requiere air)"
	@echo ""
	@echo "  🧪 Testing y Calidad:"
	@echo "    make test         - Ejecuta los tests"
	@echo "    make test-coverage - Ejecuta tests con cobertura"
	@echo "    make fmt          - Formatea el código"
	@echo "    make vet          - Ejecuta go vet"
	@echo "    make lint         - Ejecuta golangci-lint"
	@echo ""
	@echo "  📦 Git (rama actual: $(BRANCH)):"
	@echo "    make push m='mensaje' - Add + Commit + Push a $(BRANCH)"
	@echo "    make pull             - Pull desde origin/$(BRANCH)"
	@echo "    make status           - Ver estado de git"
	@echo "    make sync m='mensaje' - Pull + Push (sincronizar)"
	@echo ""
	@echo "  🧹 Utilidades:"
	@echo "    make clean        - Limpia archivos generados"
	@echo "    make deps         - Descarga las dependencias"
	@echo "    make install-tools - Instala herramientas de desarrollo"
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

install-tools: ## Instala herramientas de desarrollo
	@echo "🛠️  Instalando herramientas de desarrollo..."
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# ============================================
# COMANDOS GIT
# ============================================

# Push rápido: make push m="tu mensaje"
push:
	@if [ -z "$(m)" ]; then \
		echo "❌ Error: Debes proporcionar un mensaje"; \
		echo "   Uso: make push m='tu mensaje de commit'"; \
		exit 1; \
	fi
	@echo "📦 Agregando archivos..."
	@git add .
	@echo "✍️  Commiteando: $(m)"
	@git commit -m "$(m)"
	@echo "🚀 Pusheando a origin/$(BRANCH)..."
	@git push origin $(BRANCH)
	@echo "✅ Push completado exitosamente!"

# Pull desde origin
pull:
	@echo "⬇️  Pulling desde origin/$(BRANCH)..."
	@git fetch origin
	@git pull origin $(BRANCH)
	@echo "✅ Pull completado!"

# Ver estado de git
status:
	@echo "📊 Estado de Git (rama: $(BRANCH)):"
	@echo ""
	@git status

# Sincronizar (pull + push)
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
	@echo "✍️  Commiteando: $(m)"
	@git commit -m "$(m)"
	@echo "🚀 Pusheando a origin/$(BRANCH)..."
	@git push origin $(BRANCH)
	@echo "✅ Sincronización completada!"

# Comandos de base de datos (para futuras implementaciones)
db-migrate: ## Ejecuta migraciones (cuando se implemente)
	@echo "🗃️  Migraciones de base de datos no implementadas aún"

db-seed: ## Ejecuta seeders (cuando se implemente)
	@echo "🌱 Seeders de base de datos no implementados aún"
