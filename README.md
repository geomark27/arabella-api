# Arabella Financial OS

> Sistema Operativo Financiero Personal con Contabilidad de Doble Partida

**Versión:** v1.0.0 - Phase 1 (95% Completo)  
**Estado:** ✅ Backend Funcional | ⚠️ Testing Pendiente | ❌ Frontend No Iniciado

---

## 🎯 ¿Qué es Arabella?

Arabella no es solo otra app de gastos. Es un **Sistema Operativo Financiero** que usa los mismos principios contables que los bancos (doble partida), pero con una experiencia de usuario simplificada.

### El Problema que Resuelve

**Freelancers y trabajadores remotos en Latinoamérica** que cobran en USD/EUR pero gastan en moneda local sufren:
- 💸 **Ilusión de liquidez**: "Tengo $5,000 en el banco" → pero $2,000 son para impuestos
- 📉 **Pérdidas invisibles**: No rastrean pérdidas por tipo de cambio y comisiones
- ⚠️ **Sorpresas fiscales**: Llega el SAT y no tienen dinero apartado
- 📊 **Sin visibilidad real**: No saben cuántos meses pueden sobrevivir con sus ahorros

### La Solución: Arabella

```
✅ Calcula tu "Runway" (meses de supervivencia) en tiempo real
✅ Aparta impuestos automáticamente (Tax Shield)
✅ Rastrea pérdidas por tipo de cambio
✅ Contabilidad real (como un banco) pero fácil de usar
✅ Multi-moneda nativo
```

---

## ⭐ Feature Estrella: Runway Calculation

```
Runway = (Activos Líquidos - Deudas a Corto Plazo) / Gastos Mensuales Promedio

Ejemplo:
- Activos líquidos: $10,000
- Deudas pendientes: $2,000
- Gastos promedio: $2,000/mes

Runway = ($10,000 - $2,000) / $2,000 = 4 meses ⚠️
```

**El usuario ve:** "Te quedan 4 meses de runway. Considera reducir gastos."

---

## 🚀 Estado del Proyecto

### ✅ Implementado (Backend - 95%)

| Componente | Estado | Descripción |
|------------|--------|-------------|
| **Motor Contable** | ✅ 100% | Double-entry bookkeeping funcionando |
| **API REST** | ✅ 95% | 30+ endpoints implementados |
| **Dashboard** | ✅ 100% | Con Runway y métricas clave |
| **Multi-moneda** | ✅ 80% | Básico funcionando |
| **Users** | ✅ 90% | CRUD completo con bcrypt |
| **Auth JWT** | ⚠️ 30% | Modelos listos, falta middleware |
| **Tests** | ❌ 0% | **CRÍTICO** - Pendiente |
| **Frontend** | ❌ 0% | Fase 3 (no iniciada) |

### 📊 Números del Proyecto

- **35+ archivos** Go
- **7 modelos** de datos
- **9 services** con lógica de negocio
- **8 handlers** HTTP
- **30+ endpoints** REST
- **0 tests** 😱 (próxima prioridad)

---

## 🏗️ Arquitectura

```
arabella-api/
├── cmd/
│   ├── arabella-api/    # Main API server
│   └── console/         # CLI tools
├── internal/
│   ├── app/
│   │   ├── models/      # 7 modelos (User, Account, Transaction...)
│   │   ├── dtos/        # Data Transfer Objects
│   │   ├── repositories/# Capa de datos (GORM)
│   │   ├── services/    # Lógica de negocio ⭐
│   │   └── handlers/    # HTTP handlers (Gin)
│   ├── database/        # DB setup, migrations, seeders
│   ├── platform/
│   │   ├── config/      # Configuración
│   │   └── server/      # Server setup, routes
│   └── shared/
│       └── middleware/  # CORS, Auth (WIP)
└── docs/                # Documentación completa
```

**Arquitectura:** Clean Architecture (Hexagonal)  
**ORM:** GORM  
**Framework:** Gin  
**Base de Datos:** PostgreSQL 14+

---

## 🏃‍♂️ Inicio Rápido

### Prerequisitos

- Go 1.22+
- PostgreSQL 14+
- (Opcional) Docker y Docker Compose

### Instalación

```bash
# 1. Clonar el repositorio
git clone [tu-repo]
cd arabella-api

# 2. Configurar variables de entorno
cp .env.example .env
# Editar .env con tus credenciales de PostgreSQL

# 3. Instalar dependencias
go mod tidy

# 4. Ejecutar el servidor
go run cmd/arabella-api/main.go
```

El servidor estará disponible en: **http://localhost:8080**

### Comandos Disponibles

```bash
# Ver todos los comandos
make help

# Compilar
make build

# Ejecutar
make run

# Tests (cuando estén implementados)
make test

# Formatear código
make fmt

# Analizar código
make vet

# Limpiar binarios
make clean
```

---

## 🔌 API Endpoints Principales

### Health & Info
- `GET /` - Información de bienvenida
- `GET /api/v1/health` - Health check
- `GET /api/v1/health/ready` - Readiness check

### Dashboard ⭐
- `GET /api/v1/dashboard` - Dashboard completo
- `GET /api/v1/dashboard/runway` - Cálculo de Runway
- `GET /api/v1/dashboard/monthly-stats` - Estadísticas mensuales

### Users
- `GET /api/v1/users` - Listar usuarios
- `POST /api/v1/users` - Crear usuario
- `GET /api/v1/users/:id` - Obtener usuario
- `PUT /api/v1/users/:id` - Actualizar usuario
- `DELETE /api/v1/users/:id` - Eliminar usuario

### Accounts
- `GET /api/v1/accounts` - Listar cuentas
- `POST /api/v1/accounts` - Crear cuenta
- `GET /api/v1/accounts/:id` - Obtener cuenta
- `PUT /api/v1/accounts/:id` - Actualizar cuenta
- `DELETE /api/v1/accounts/:id` - Eliminar cuenta

### Transactions
- `GET /api/v1/transactions` - Listar transacciones
- `POST /api/v1/transactions` - Crear transacción (pasa por motor contable)
- `GET /api/v1/transactions/:id` - Obtener transacción
- `PUT /api/v1/transactions/:id` - Actualizar transacción
- `DELETE /api/v1/transactions/:id` - Eliminar transacción

### Categories, Currencies, Journal Entries
Ver documentación completa en: **[docs/API.md](docs/API.md)**

---

## 🧪 Ejemplos de Uso

### 1. Crear Usuario
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "María",
    "last_name": "González",
    "email": "maria@example.com"
  }'
```

### 2. Crear Cuenta Bancaria
```bash
curl -X POST http://localhost:8080/api/v1/accounts \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Banco BBVA USD",
    "account_type": "BANK",
    "classification": "ASSET",
    "currency_code": "USD",
    "initial_balance": "5000.00"
  }'
```

### 3. Registrar un Gasto
```bash
curl -X POST http://localhost:8080/api/v1/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "type": "EXPENSE",
    "amount": "150.00",
    "account_id": 1,
    "category_id": 2,
    "description": "Compras de supermercado",
    "transaction_date": "2026-01-03T10:00:00Z"
  }'
```

### 4. Ver Dashboard con Runway
```bash
curl http://localhost:8080/api/v1/dashboard

# Respuesta incluye:
# - total_assets
# - total_liabilities
# - net_worth
# - runway (en meses) ⭐
# - monthly_income
# - monthly_expenses
# - account_balances[]
```

---

## 📚 Documentación

- **[DOC.md](docs/DOC.md)** - Documentación técnica completa
- **[PROJECT_STATUS.md](docs/PROJECT_STATUS.md)** - Estado actual del proyecto
- **[CHECKLIST.md](CHECKLIST.md)** - Roadmap y tareas pendientes
- **[BUSINESS_MODEL.md](docs/BUSINESS_MODEL.md)** - Modelo de negocio
- **[API.md](docs/API.md)** - Documentación de API
- **[USER_GUIDE.md](docs/USER_GUIDE.md)** - Guía de usuario

---

## 🎯 Próximos Pasos

### Inmediatos (Fase 2 - Semanas 5-7)
1. ⬜ **Tests unitarios** del motor contable (CRÍTICO)
2. ⬜ **Autenticación JWT** completa
3. ⬜ **Docker setup** para desarrollo
4. ⬜ Resolver TODOs de userID hardcodeado

### Medio Plazo (Fase 3 - Semanas 8-11)
5. ⬜ Frontend con Next.js 14
6. ⬜ PWA básico
7. ⬜ UI optimizada para móvil

### Largo Plazo (Fases 4-6)
8. ⬜ Tax Shield automático
9. ⬜ Email parsing (AWS SES)
10. ⬜ Beta cerrada con usuarios reales

Ver [CHECKLIST.md](CHECKLIST.md) para roadmap completo.

---

## 🤝 Contribuir

Este es un proyecto personal/side project, pero si estás interesado:

1. Fork el repositorio
2. Crea una rama (`git checkout -b feature/amazing-feature`)
3. Commit tus cambios (`git commit -m 'Add amazing feature'`)
4. Push a la rama (`git push origin feature/amazing-feature`)
5. Abre un Pull Request

---

## 📝 Licencia

[MIT](LICENSE) - Proyecto personal de Marcos Ramos

---

## 👤 Autor

**Marcos Ramos** - Senior Software Engineer  
Trabajando en construir el sistema financiero que siempre quise tener.

---

## 🙏 Agradecimientos

- Inspirado en la necesidad real de freelancers latinoamericanos
- Construido con las mejores prácticas de Go y Clean Architecture
- Feature estrella (Runway) basado en metodología de startups

---

**Estado:** 🚧 En desarrollo activo - Fase 1 casi completa (95%)  
**Última actualización:** Enero 3, 2026

## 🏗️ Arquitectura

Este proyecto sigue el patrón de **arquitectura por capas** inspirado en frameworks como NestJS:

- **Handlers**: Manejan las peticiones HTTP y las respuestas
- **Services**: Contienen la lógica de negocio
- **Repositories**: Manejan la persistencia de datos  
- **DTOs**: Definen la estructura de datos de entrada/salida
- **Models**: Representan las entidades del dominio
- **Middleware**: Procesan las peticiones de forma transversal

## 📦 Helpers de Loom

Este proyecto usa los helpers opcionales de Loom para desarrollo más rápido:

- `helpers.RespondJSON()` - Respuestas HTTP estandarizadas
- `helpers.ValidateStruct()` - Validación de structs
- `helpers.Logger` - Logging estructurado
- `helpers.AppError` - Manejo de errores mejorado

Para actualizar los helpers:
```bash
go get -u github.com/geomark27/loom-go
```

## 🔧 Configuración

Las variables de entorno se definen en `.env`:

```bash
PORT=8080
ENVIRONMENT=development
LOG_LEVEL=info
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080
```

## 📚 Próximos Pasos

1. **Agregar base de datos**: Reemplazar el repositorio en memoria
2. **Implementar autenticación**: JWT, OAuth, etc.
3. **Agregar validaciones**: Validador de structs más robusto
4. **Tests**: Crear tests unitarios e integración  
5. **Logging**: Implementar logging estructurado
6. **Métricas**: Prometheus, health checks avanzados

## 🛠️ Generado con

Este proyecto fue generado con [**Loom**](https://github.com/geomark27/loom-go) - El tejedor de proyectos Go.

¡Disfruta desarrollando con Go! 🎉
