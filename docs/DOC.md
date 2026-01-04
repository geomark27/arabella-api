# Documentación de Proyecto: Financial OS (Nombre Clave: \"Arabella-fos\")

Versión: 1.0

Autor: Marcos Ramos (Senior Software Engineer)

Fecha: 17 de Diciembre, 2025

## 1. Resumen Ejecutivo

**Arabella-fos** no es una simple aplicación de seguimiento de gastos.
Es un **Sistema Operativo Financiero** diseñado bajo principios de
contabilidad de doble entrada (*double-entry bookkeeping*), pero con una
capa de experiencia de usuario simplificada.

Su objetivo principal es resolver la \"ilusión de liquidez\" que sufren
los trabajadores remotos y freelancers en Latinoamérica que perciben
ingresos en moneda extranjera, automatizando la gestión de impuestos,
tipos de cambio y proyecciones de flujo de caja (*Runway*).

## 2. Modelo de Negocio y Mercado

### Público Objetivo (Target)

1.  **Principal (Nicho):** Desarrolladores, Diseñadores y Freelancers en
    > Latam que trabajan para el exterior (ingresos en USD/EUR, gastos
    > en moneda local).

2.  **Secundario:** Usuarios generales que buscan una gestión financiera
    > estricta y real (no solo \"lista de compras\").

### Principales Competidores

-   **Excel / Google Sheets:** Máxima flexibilidad, pero manual y
    > propenso a errores. Sin automatización.

-   **Apps de Gastos (Wallet, Bluecoins, Monefy):** Simples, pero
    > carecen de lógica contable real (no manejan activos/pasivos
    > correctamente) y gestión fiscal.

-   **ERPs (QuickBooks, Xero):** Demasiado complejos y costosos para un
    > individuo.

###  

### Propuesta de Valor Única (Ventajas Competitivas)

1.  **Integridad Contable Invisible:** Usa el mismo sistema que un banco
    > (Doble Entrada) para garantizar que el dinero nunca desaparece,
    > pero el usuario solo ve \"Ingresos y Gastos\".

2.  **Realidad Multi-moneda:** Cálculo automático de pérdidas por tipo
    > de cambio y comisiones (*Spread*).

3.  **Escudo Fiscal (Tax Shield):** Apartado virtual automático de
    > impuestos basado en reglas configurables.

4.  **Automatización \"Serverless\":** Ingesta de transacciones vía
    > reenvío de correos electrónicos (sin APIs bancarias costosas).

## 3. Lógica de Negocio (Core Domain)

### Principios Fundamentales

-   **La Ecuación Contable:** Activos = Pasivos + Patrimonio. Cada
    > transacción mueve dinero de una cuenta a otra. Nada se crea ni se
    > destruye sin rastro.

-   **Abstracción de Cuentas:**

    -   *Cuentas Reales:* Bancos, Efectivo, Tarjetas de Crédito.

    -   *Cuentas Nominales:* Categorías de Gasto (Comida, Servicios),
        > Categorías de Ingreso.

    -   *Cuentas Virtuales:* \"Buckets\" de impuestos o ahorro.

### Flujos Críticos

1.  **Registro de Gasto:** Debita una cuenta de Gasto, Acredita una
    > cuenta de Activo (Banco) o Pasivo (Tarjeta).

2.  **Gestión de Deuda:** Al pagar con tarjeta, el saldo bancario no
    > baja inmediatamente. Se crea una deuda. Al pagar la tarjeta, se
    > reduce el Activo (Banco) y se reduce el Pasivo (Tarjeta).

3.  **Cálculo de Runway:** (Total Activos Líquidos - Pasivos a Corto
    > Plazo) / Promedio de Gastos Mensuales.

## 4. Stack Tecnológico & Arquitectura

### Frontend (Cliente)

-   **Framework:** **Next.js 14+** (App Router).

-   **Lenguaje:** TypeScript.

-   **Styling:** Tailwind CSS (para desarrollo rápido).

-   **Estrategia:** PWA (Progressive Web App) para capacidades
    > Offline-First.

-   **State Management:** TanStack Query (React Query) para
    > sincronización eficiente con el backend.

### Backend (API)

-   **Lenguaje:** **Go (Golang)** 1.22+.

-   **Arquitectura:** Clean Architecture / Hexagonal.

    -   *Domain:* Entidades puras y lógica de negocio.

    -   *Application:* Casos de uso.

    -   *Infrastructure:* Implementación de base de datos y adaptadores
        > AWS.

-   **Concurrencia:** Uso de Goroutines para procesamiento batch de
    > correos/CSVs.

### Base de Datos

-   **Motor:** **PostgreSQL 14+** (cambio de SQL Server por simplicidad y costo).

-   **ORM:** **GORM** (Go Object-Relational Mapping).

-   **Features Clave:**

    -   **Soft Deletes:** GORM maneja eliminación lógica con `deleted_at`.

    -   **Auto Migrations:** Sincronización automática de esquema con modelos.

    -   **Constraints:** Validaciones a nivel de BD para integridad.

    -   **Tipos de Datos:** DECIMAL para precisión monetaria absoluta.

    -   **Timestamps:** created_at, updated_at automáticos.

### Infraestructura (AWS - Serverless First)

-   **Compute:** AWS Lambda (funciones Go) para la API y workers.

-   **API Gateway:** Entrada para el Frontend.

-   **Ingesta:** **AWS SES (Simple Email Service)** -\> Regla de
    > recepción -\> S3 Bucket -\> Lambda Trigger (Parser).

-   **Almacenamiento:** S3 (para almacenar los correos crudos o recibos
    > escaneados).

## 5. Definición del MVP (Producto Mínimo Viable)

El MVP se centra en la **solidez del dato** e **integridad contable**. 
No buscamos automatización total en día 1, sino integridad absoluta del sistema.

### Estado Actual del MVP

#### ✅ IMPLEMENTADO (Backend - 95%)

1.  **✅ Gestión de Cuentas Completa**
    - CRUD de cuentas bancarias (Multi-moneda: USD, MXN, EUR, etc.)
    - Tipos: BANK, CASH, CREDIT_CARD, INVESTMENT, CATEGORY
    - Clasificación: ASSET, LIABILITY, EQUITY, INCOME, EXPENSE
    - Balance automático actualizado por motor contable

2.  **✅ Motor de Doble Partida (Accounting Engine)**
    - Procesamiento automático de transacciones
    - Generación de Journal Entries (Debit + Credit)
    - Validación: SUM(Debits) = SUM(Credits)
    - Transacciones atómicas (rollback en errores)
    - Actualización automática de balances

3.  **✅ Gestión de Transacciones**
    - Tipos: INCOME, EXPENSE, TRANSFER, DEBT_PAYMENT
    - Manejo de categorías
    - Multi-moneda con tipo de cambio manual
    - Notas y metadatos
    - Audit trail completo (Journal Entries)

4.  **✅ Dashboard con Feature Estrella**
    - **Cálculo de Runway** (meses de supervivencia)
    - Net Worth (Activos - Pasivos)
    - Activos líquidos vs totales
    - Gastos/Ingresos mensuales
    - Balance por cuenta

5.  **✅ Multi-moneda Básico**
    - Soporte de múltiples divisas
    - Tipo de cambio manual al momento de transacción
    - Cálculo de balances en diferentes monedas

6.  **✅ Sistema de Usuarios**
    - Registro de usuarios
    - Autenticación con bcrypt
    - Multi-tenant ready (userID en todas las entidades)

#### ⚠️ PENDIENTE (Para MVP Completo)

1.  **Frontend (Fase 3 - 0%)**
    - Interfaz web en Next.js
    - Formularios de registro rápido
    - Dashboard visual
    - Gráficas de tendencias

2.  **Autenticación JWT (Fase 2 - 0%)**
    - Login endpoint con generación de tokens
    - Middleware de protección
    - Refresh tokens
    - Extracción de userID del contexto

3.  **Tests Unitarios (Fase 2 - 0%)**
    - Tests del Accounting Engine
    - Tests de validaciones
    - Tests de integración de endpoints
    - Cobertura objetivo: >70%

4.  **Tax Shield Automático (Fase 4 - 0%)**
    - Configuración de reglas fiscales
    - Apartado automático de impuestos
    - Cuentas virtuales

5.  **Email Parsing (Fase 5 - 0%)**
    - AWS SES + Lambda
    - Parsers de emails bancarios
    - Queue de aprobación

*Nota: La automatización por correo (AWS SES) queda para Fase 5, 
el foco inicial es solidez del sistema manual.*

## 6. Fases y Tiempos Estimados (Roadmap)

Considerando un dedicación "Side Project" (10-15 horas semanales).

### ✅ FASE 1: COMPLETADA - Fundamentos y Motor Contable

**Estado:** 95% Completo | **Tiempo Real:** Semanas 1-4

#### Implementado:
- ✅ **Arquitectura Clean Architecture** (handlers → services → repositories → models)
- ✅ **Base de Datos PostgreSQL** con GORM (migraciones automáticas)
- ✅ **Motor de Contabilidad de Doble Partida** (`AccountingEngineService`)
  - Generación automática de asientos contables (Journal Entries)
  - Validación Débito = Crédito
  - Transacciones atómicas en BD
- ✅ **Modelos Completos:**
  - Users (con bcrypt para passwords)
  - Accounts (Bancos, Efectivo, Tarjetas, Categorías)
  - Transactions (INCOME, EXPENSE, TRANSFER, DEBT_PAYMENT)
  - JournalEntries (Audit trail completo)
  - Categories (Gastos/Ingresos)
  - Currencies (Multi-moneda)
  - SystemValues (Configuración)

- ✅ **API REST Completa (30+ endpoints):**
  - Health checks
  - Users CRUD
  - Accounts CRUD
  - Transactions CRUD (pasa por Accounting Engine)
  - Categories CRUD
  - Currencies (read-only)
  - Journal Entries (read-only, auditoría)
  - **Dashboard con Runway Calculation** ⭐ (Feature Estrella)

- ✅ **Features Avanzadas:**
  - Cálculo de Runway (meses de supervivencia financiera)
  - Stats mensuales (ingresos, gastos, balance)
  - Separación de Activos Líquidos vs No Líquidos
  - Manejo de múltiples monedas
  - CORS configurado
  - Seeders de base de datos

#### Pendiente de Fase 1 (5%):
- ⚠️ Tests unitarios del motor contable
- ⚠️ Documentación Swagger/OpenAPI completa
- ⚠️ Dockerfile para desarrollo local

---

### 🔄 FASE 2: Backend Hardening y Testing (Semanas 5-7)

**Objetivo:** Asegurar calidad y robustez del backend antes del frontend

#### Tareas Críticas:
1. **Testing (Semana 5)**
   - Tests unitarios del `AccountingEngineService`
   - Tests de integración para endpoints críticos
   - Tests de validación de balance contable
   - Cobertura mínima objetivo: 70%

2. **Validaciones y Seguridad (Semana 6)**
   - Middleware de autenticación JWT
   - Rate limiting básico
   - Validación robusta de DTOs
   - Sanitización de inputs
   - Error handling consistente

3. **Optimizaciones (Semana 7)**
   - Índices en BD para queries frecuentes
   - Paginación en endpoints de listado
   - Caching básico (si es necesario)
   - Logging estructurado
   - Docker Compose para desarrollo

#### Entregables:
- ✅ Suite de tests con >70% cobertura
- ✅ API documentada (Swagger)
- ✅ Sistema de autenticación funcional
- ✅ Docker setup para desarrollo

---

### 🎨 FASE 3: Frontend MVP (Semanas 8-11)

**Objetivo:** Interfaz funcional para las operaciones core

#### Stack Frontend:
- Next.js 14+ (App Router)
- TypeScript
- Tailwind CSS
- TanStack Query (React Query)
- Shadcn/UI o similar

#### Pantallas Mínimas (Semana 8-9):
1. **Login/Registro**
2. **Dashboard Principal**
   - Resumen financiero
   - **Runway Display** (feature estrella)
   - Gráfica de tendencia mensual
3. **Cuentas**
   - Lista de cuentas
   - Crear/Editar cuenta
   - Ver detalle y transacciones
4. **Transacciones**
   - Formulario rápido de registro
   - Lista con filtros
   - Detalle de transacción

#### Features UX (Semana 10-11):
- Formulario optimizado "one-thumb" para móvil
- Categorías con íconos
- Conversión de monedas en tiempo real
- Validaciones inline
- Estados de carga y errores
- PWA básico (Service Workers)

#### Entregables:
- ✅ App funcional desplegada
- ✅ CRUD completo desde UI
- ✅ Dashboard con métricas clave
- ✅ Responsive design

---

### 🚀 FASE 4: Features Avanzadas y Tax Shield (Semanas 12-14)

**Objetivo:** Implementar diferenciadores clave del producto

#### Features a Implementar:
1. **Tax Shield Automático** ⭐
   - Configuración de reglas de impuestos por usuario
   - Cuentas virtuales automáticas
   - Apartado automático en cada ingreso
   - Dashboard de obligaciones fiscales

2. **Reportes y Exports**
   - Export a CSV/Excel
   - Reporte PDF mensual
   - Gráficas avanzadas (Recharts/Tremor)
   - Tendencias por categoría

3. **Gestión de Deuda Inteligente**
   - Vista de deudas próximas a vencer
   - Recordatorios de pagos
   - Simulador de pago de deudas

4. **Multi-moneda Avanzado**
   - Integración API de tipos de cambio
   - Detección automática de pérdidas por spread
   - Conversión automática en reportes

#### Entregables:
- ✅ Tax Shield funcional
- ✅ Sistema de reportes
- ✅ Alertas y notificaciones
- ✅ Manejo avanzado de divisas

---

### ☁️ FASE 5: Automatización y Escalabilidad (Semanas 15-17)

**Objetivo:** Reducir fricción de entrada manual

#### Implementación AWS:
1. **Email Parsing (AWS SES + Lambda)**
   - Configuración de AWS SES
   - Lambda parser de emails bancarios
   - S3 para almacenar emails originales
   - Queue de transacciones pendientes de aprobar

2. **Infraestructura como Código**
   - Terraform o AWS CDK
   - CI/CD con GitHub Actions
   - Ambientes (dev, staging, prod)
   - Backups automáticos de BD

3. **Monitoreo y Observabilidad**
   - CloudWatch Logs
   - Métricas de uso
   - Alertas de errores
   - APM básico

#### Entregables:
- ✅ Ingesta automática de emails
- ✅ Despliegue automatizado
- ✅ Monitoreo en producción
- ✅ Backups configurados

---

### 🎯 FASE 6: Pulido y Lanzamiento Beta (Semanas 18-20)

**Objetivo:** Preparar para primeros usuarios reales

#### Actividades:
1. **Testing End-to-End**
   - Casos de uso completos
   - Testing con usuarios beta
   - Corrección de bugs

2. **Documentación**
   - Guía de usuario
   - FAQs
   - Videos tutoriales

3. **Performance**
   - Optimización de queries lentas
   - Lazy loading en frontend
   - Compresión de assets

4. **Legal y Términos**
   - Términos y condiciones
   - Política de privacidad
   - GDPR básico

#### Entregables:
- ✅ Beta cerrada con 10-20 usuarios
- ✅ Documentación completa
- ✅ Sistema estable y monitoreado

---

**Tiempo Total Estimado al MVP Completo:** 4-5 meses

**Estado Actual:** Fase 1 casi completa, listo para Fase 2

## 7. Posibles Áreas de Mejora y Riesgos

### Riesgos Técnicos

1.  **Ausencia de Tests** ⚠️ CRÍTICO
    -   *Impacto:* El motor contable sin tests puede tener bugs que causen inconsistencias de datos.
    -   *Mitigación:* Priorizar tests del `AccountingEngineService` en Fase 2 antes de continuar.

2.  **Fricción de Usuario:** Si el registro manual es lento, el usuario abandona.
    -   *Mitigación:* UI optimizada "One-thumb" (uso con una mano), formularios inteligentes con defaults.

3.  **Variabilidad de Emails:** Los bancos cambian formatos de correo.
    -   *Mitigación:* Arquitectura de parsers modulares (Strategy Pattern) fácil de actualizar sin redeployar todo el backend.

4.  **Costos de Infraestructura:**
    -   PostgreSQL en RDS puede ser costoso inicialmente.
    -   *Mitigación:* Iniciar con PostgreSQL en Docker local/Railway/Supabase, migrar a RDS solo al escalar.
    -   Lambda puede generar costos inesperados con alto tráfico.
    -   *Mitigación:* Implementar rate limiting y monitoreo de costos desde día 1.

### Riesgos de Producto

5.  **Complejidad Oculta:** Los usuarios pueden no entender por qué su balance difiere de lo que ven en el banco.
    -   *Mitigación:* UI clara que explique transacciones pendientes, tipos de cambio, etc.

6.  **Abandono Temprano:** Si no ven valor en los primeros 7 días, no regresan.
    -   *Mitigación:* Onboarding guiado que muestre el Runway desde el primer día con datos de ejemplo.

---

## 8. Siguiente Paso Inmediato

### Prioridad #1: Tests del Motor Contable

Antes de avanzar al frontend, **es crítico** tener tests del `AccountingEngineService`:

```go
// Casos de prueba mínimos requeridos:
1. Test_ProcessExpense_DebitCredit_Balance
2. Test_ProcessIncome_DebitCredit_Balance
3. Test_ProcessTransfer_MultiCurrency
4. Test_ProcessDebtPayment_CreditCard
5. Test_ReverseTransaction_RestoresBalances
6. Test_InvalidTransaction_RollbackOnError
7. Test_ConcurrentTransactions_ThreadSafety
```

### Prioridad #2: Autenticación JWT

Implementar middleware de autenticación para proteger endpoints:
- Login endpoint que genere JWT
- Middleware que valide token en cada request
- Extracción de `userID` del token (resolver los TODOs actuales)

### Prioridad #3: Docker Development Setup

Crear `docker-compose.yml` para levantar:
- PostgreSQL
- API de Arabella
- (Opcional) Adminer/pgAdmin para visualizar BD

Esto facilitará onboarding de nuevos desarrolladores y testing local.

---

## 9. Estado Actual del Proyecto (Enero 2026)

### Lo que Ya Funciona ✅
- Backend API completo con 30+ endpoints
- Motor de contabilidad de doble partida funcionando
- Dashboard con cálculo de Runway
- Multi-moneda básico
- Seeders para datos de prueba
- CORS configurado

### Lo que Falta para MVP ⚠️
- Tests unitarios (CRÍTICO)
- Autenticación JWT
- Frontend (Next.js)
- Tax Shield automático
- Email parsing (AWS SES)
- Despliegue en cloud

### Deuda Técnica Conocida
- TODOs en handlers: `userID` hardcodeado a 1
- Falta validación exhaustiva de DTOs
- Sin logging estructurado
- Sin manejo de errores consistente
- Documentación Swagger incompleta

---

## 10. Métricas de Éxito del MVP

Para considerar el MVP "listo para beta":

1. **Funcionalidad:**
   - [ ] Registro de 100 transacciones sin errores
   - [ ] Cálculo de Runway preciso
   - [ ] Balance contable siempre cuadra (Debits = Credits)

2. **Calidad:**
   - [ ] Cobertura de tests >70%
   - [ ] API response time <200ms (p95)
   - [ ] Zero downtime en 7 días

3. **UX:**
   - [ ] Registrar gasto toma <10 segundos
   - [ ] Dashboard carga en <1 segundo
   - [ ] Funciona en móvil

4. **Usuarios:**
   - [ ] 10 usuarios beta activos
   - [ ] Retención 7 días >60%
   - [ ] Al menos 1 usuario con >30 días de uso

---