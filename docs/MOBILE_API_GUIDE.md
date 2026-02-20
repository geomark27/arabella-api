# Arabella Financial OS — Guía de API para Frontend Móvil

> Versión: **v1.0.0** · Base URL: `http://localhost:8080/api/v1` · Formato: **JSON**

---

## Índice

1. [Autenticación](#1-autenticación)
2. [Convenciones generales](#2-convenciones-generales)
3. [Módulo: Auth](#3-módulo-auth)
4. [Módulo: Users](#4-módulo-users)
5. [Módulo: Accounts](#5-módulo-accounts)
6. [Módulo: Transactions](#6-módulo-transactions)
7. [Módulo: Categories](#7-módulo-categories)
8. [Módulo: Dashboard](#8-módulo-dashboard)
9. [Módulo: Currencies](#9-módulo-currencies)
10. [Módulo: System Values](#10-módulo-system-values)
11. [Módulo: Journal Entries](#11-módulo-journal-entries)
12. [Catálogo de valores](#12-catálogo-de-valores)
13. [Códigos de error](#13-códigos-de-error)

---

## 1. Autenticación

Todos los endpoints marcados con 🔒 requieren un **JWT Access Token** en el header:

```
Authorization: Bearer <access_token>
```

El token se obtiene en `POST /auth/login` o `POST /auth/register`.  
Duración del access token: **1 hora**.  
Duración del refresh token: **30 días**.

---

## 2. Convenciones generales

### Estructura de respuesta exitosa
```json
{
  "data": { ... },
  "count": 10
}
```

### Estructura de respuesta de error
```json
{
  "error": "Mensaje corto del error",
  "details": "Detalle técnico del problema"
}
```

### Tipos de datos
| Tipo | Formato | Ejemplo |
|------|---------|---------|
| Fechas | ISO 8601 / RFC3339 | `"2026-02-19T14:35:40Z"` |
| Montos | String decimal | `"1500.0000"` |
| IDs | Integer | `1` |
| Booleanos | Boolean | `true` / `false` |

### Paginación (endpoints de lista)
Los endpoints de lista aceptan los query params:
- `page` — Número de página (default: `1`)
- `page_size` — Elementos por página (default varía, máximo: `100`)

Respuesta paginada:
```json
{
  "transactions": [...],
  "total": 150,
  "page": 1,
  "page_size": 20,
  "total_pages": 8
}
```

---

## 3. Módulo: Auth

Endpoints públicos (sin token). Manejan el ciclo completo de autenticación.

---

### 3.1 Registro

**`POST /auth/register`** — Crea una nueva cuenta y devuelve tokens.

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "marcos@arabella.app",
    "password": "MiPassword123",
    "first_name": "Marcos",
    "last_name": "Ramos"
  }'
```

**Body:**
| Campo | Tipo | Requerido | Descripción |
|-------|------|-----------|-------------|
| `email` | string | ✅ | Email único del usuario |
| `password` | string | ✅ | Mínimo 8 caracteres |
| `first_name` | string | ✅ | Nombre |
| `last_name` | string | ✅ | Apellido |

**Respuesta `201`:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "marcos@arabella.app",
    "first_name": "Marcos",
    "last_name": "Ramos",
    "is_active": true
  }
}
```

**Errores:**
| Código | Causa |
|--------|-------|
| `400` | Datos inválidos o contraseña menor a 8 caracteres |
| `409` | El email ya está registrado |

---

### 3.2 Login

**`POST /auth/login`** — Autentica al usuario y devuelve tokens.

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "marcos@arabella.app",
    "password": "MiPassword123"
  }'
```

**Body:**
| Campo | Tipo | Requerido |
|-------|------|-----------|
| `email` | string | ✅ |
| `password` | string | ✅ |

**Respuesta `200`:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "marcos@arabella.app",
    "first_name": "Marcos",
    "last_name": "Ramos",
    "is_active": true
  }
}
```

**Errores:**
| Código | Causa |
|--------|-------|
| `401` | Credenciales incorrectas o usuario inactivo |

---

### 3.3 Refrescar token

**`POST /auth/refresh`** — Genera un nuevo par de tokens usando el refresh token.

```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }'
```

**Respuesta `200`:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

> ⚠️ Guarda el nuevo `refresh_token` devuelto. El anterior queda inválido.

---

### 3.4 Cambiar contraseña 🔒

**`PUT /auth/change-password`** — Cambia la contraseña del usuario autenticado.

```bash
curl -X PUT http://localhost:8080/api/v1/auth/change-password \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "MiPassword123",
    "new_password": "NuevoPassword456"
  }'
```

**Respuesta `200`:**
```json
{
  "message": "password changed successfully"
}
```

**Errores:**
| Código | Causa |
|--------|-------|
| `400` | Contraseña actual incorrecta o nueva contraseña menor a 8 caracteres |

---

## 4. Módulo: Users

Gestión del perfil de usuario. Todos los endpoints son 🔒.

---

### 4.1 Listar usuarios

**`GET /users`** 🔒

```bash
curl -X GET http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer <access_token>"
```

**Respuesta `200`:**
```json
{
  "data": [
    {
      "id": 1,
      "first_name": "Marcos",
      "last_name": "Ramos",
      "email": "marcos@arabella.app"
    }
  ],
  "count": 1
}
```

---

### 4.2 Obtener usuario por ID

**`GET /users/:id`** 🔒

```bash
curl -X GET http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer <access_token>"
```

**Respuesta `200`:**
```json
{
  "data": {
    "id": 1,
    "first_name": "Marcos",
    "last_name": "Ramos",
    "email": "marcos@arabella.app"
  }
}
```

---

### 4.3 Crear usuario

**`POST /users`** 🔒

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "nuevo@arabella.app",
    "first_name": "Ana",
    "last_name": "García",
    "user_name": "anagarcia",
    "password": "Password123"
  }'
```

**Body:**
| Campo | Tipo | Requerido | Descripción |
|-------|------|-----------|-------------|
| `email` | string | ✅ | Email único |
| `first_name` | string | ✅ | Nombre |
| `last_name` | string | ✅ | Apellido |
| `user_name` | string | ✅ | Nombre de usuario único |
| `password` | string | ✅ | Mínimo 8 caracteres |

**Respuesta `201`:**
```json
{
  "message": "User created successfully",
  "data": {
    "id": 2,
    "first_name": "Ana",
    "last_name": "García",
    "email": "nuevo@arabella.app"
  }
}
```

---

### 4.4 Actualizar usuario

**`PUT /users/:id`** 🔒 — Solo envía los campos que quieres modificar.

```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Marcos Antonio",
    "avatar_url": "https://cdn.arabella.app/avatars/1.jpg"
  }'
```

**Body (todos opcionales):**
| Campo | Tipo | Descripción |
|-------|------|-------------|
| `first_name` | string | Nombre |
| `last_name` | string | Apellido |
| `user_name` | string | Nombre de usuario |
| `email` | string | Email (valida duplicados) |
| `avatar_url` | string | URL del avatar |

**Respuesta `200`:**
```json
{
  "message": "User updated successfully",
  "data": {
    "id": 1,
    "first_name": "Marcos Antonio",
    "last_name": "Ramos",
    "email": "marcos@arabella.app"
  }
}
```

---

### 4.5 Eliminar usuario

**`DELETE /users/:id`** 🔒 — Soft delete, no borra de la base de datos.

```bash
curl -X DELETE http://localhost:8080/api/v1/users/2 \
  -H "Authorization: Bearer <access_token>"
```

**Respuesta `200`:**
```json
{
  "message": "User deleted successfully"
}
```

---

## 5. Módulo: Accounts

Gestión de cuentas financieras (banco, efectivo, tarjeta, etc.). Todos los endpoints son 🔒.

El `user_id` se asigna automáticamente desde el JWT — no se envía en el body.

---

### 5.1 Listar cuentas

**`GET /accounts`** 🔒

```bash
curl -X GET http://localhost:8080/api/v1/accounts \
  -H "Authorization: Bearer <access_token>"
```

**Respuesta `200`:**
```json
{
  "data": [
    {
      "id": 1,
      "name": "BBVA USD",
      "account_type": "BANK",
      "balance": "5000.0000",
      "currency_id": 1,
      "currency": {
        "id": 1,
        "code": "USD",
        "symbol": "$"
      },
      "is_active": true,
      "created_at": "2026-02-01T10:00:00Z",
      "updated_at": "2026-02-19T14:35:00Z"
    }
  ],
  "count": 1
}
```

---

### 5.2 Obtener cuenta por ID

**`GET /accounts/:id`** 🔒

```bash
curl -X GET http://localhost:8080/api/v1/accounts/1 \
  -H "Authorization: Bearer <access_token>"
```

**Respuesta `200`:**
```json
{
  "data": {
    "id": 1,
    "name": "BBVA USD",
    "account_type": "BANK",
    "balance": "5000.0000",
    "currency_id": 1,
    "currency": {
      "id": 1,
      "code": "USD",
      "symbol": "$"
    },
    "is_active": true,
    "created_at": "2026-02-01T10:00:00Z",
    "updated_at": "2026-02-19T14:35:00Z"
  }
}
```

---

### 5.3 Crear cuenta

**`POST /accounts`** 🔒

```bash
curl -X POST http://localhost:8080/api/v1/accounts \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "BBVA USD",
    "account_type": "BANK",
    "currency_id": 1,
    "balance": "5000.00"
  }'
```

**Body:**
| Campo | Tipo | Requerido | Descripción |
|-------|------|-----------|-------------|
| `name` | string | ✅ | Nombre de la cuenta |
| `account_type` | string | ✅ | Ver [tipos de cuenta](#tipos-de-cuenta) |
| `currency_id` | integer | ✅ | ID de la moneda (ver [Currencies](#9-módulo-currencies)) |
| `balance` | decimal string | ❌ | Saldo inicial (default: `"0"`) |
| `is_active` | boolean | ❌ | Default: `true` |

**Respuesta `201`:**
```json
{
  "message": "Account created successfully",
  "data": {
    "id": 2,
    "name": "BBVA USD",
    "account_type": "BANK",
    "balance": "5000.0000",
    "currency_id": 1,
    "currency": { "id": 1, "code": "USD", "symbol": "$" },
    "is_active": true,
    "created_at": "2026-02-19T14:35:00Z",
    "updated_at": "2026-02-19T14:35:00Z"
  }
}
```

---

### 5.4 Actualizar cuenta

**`PUT /accounts/:id`** 🔒 — Solo envía los campos a modificar.

```bash
curl -X PUT http://localhost:8080/api/v1/accounts/1 \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "BBVA USD Principal",
    "is_active": true
  }'
```

**Body (todos opcionales):**
| Campo | Tipo | Descripción |
|-------|------|-------------|
| `name` | string | Nombre de la cuenta |
| `account_type` | string | Tipo de cuenta |
| `currency_id` | integer | ID de moneda |
| `balance` | decimal string | Saldo (usar con cuidado) |
| `is_active` | boolean | Activar/desactivar |

**Respuesta `200`:**
```json
{
  "message": "Account updated successfully"
}
```

---

### 5.5 Eliminar cuenta

**`DELETE /accounts/:id`** 🔒 — Soft delete.

```bash
curl -X DELETE http://localhost:8080/api/v1/accounts/1 \
  -H "Authorization: Bearer <access_token>"
```

**Respuesta `200`:**
```json
{
  "message": "Account deleted successfully"
}
```

---

## 6. Módulo: Transactions

⭐ **El módulo más importante.** Cada transacción pasa por el **Motor Contable de Doble Partida**, que genera automáticamente los asientos contables (Journal Entries) y actualiza los saldos de las cuentas involucradas. Todo es atómico — si algo falla, se hace rollback completo.

Todos los endpoints son 🔒.

---

### Reglas de negocio críticas

| Tipo | `category_id` | `account_to_id` | Descripción |
|------|:---:|:---:|-------------|
| `INCOME` | ✅ Requerido | ❌ No enviar | Dinero que entra a una cuenta |
| `EXPENSE` | ✅ Requerido | ❌ No enviar | Dinero que sale de una cuenta |
| `TRANSFER` | ❌ No enviar | ✅ Requerido | Movimiento entre dos cuentas propias |

> ⚠️ Violar estas reglas devuelve `400 Bad Request`.

---

### 6.1 Listar transacciones

**`GET /transactions`** 🔒

**Query params (todos opcionales):**
| Param | Tipo | Descripción |
|-------|------|-------------|
| `type` | string | `INCOME`, `EXPENSE` o `TRANSFER` |
| `account_id` | integer | Filtrar por cuenta origen |
| `category_id` | integer | Filtrar por categoría |
| `page` | integer | Número de página (default: 1) |
| `page_size` | integer | Elementos por página (default: 20, máx: 100) |

```bash
# Todas las transacciones
curl -X GET http://localhost:8080/api/v1/transactions \
  -H "Authorization: Bearer <access_token>"

# Solo gastos, página 1
curl -X GET "http://localhost:8080/api/v1/transactions?type=EXPENSE&page=1&page_size=20" \
  -H "Authorization: Bearer <access_token>"

# Por cuenta específica
curl -X GET "http://localhost:8080/api/v1/transactions?account_id=1" \
  -H "Authorization: Bearer <access_token>"
```

**Respuesta `200`:**
```json
{
  "transactions": [
    {
      "id": 10,
      "user_id": 1,
      "type": "EXPENSE",
      "description": "Supermercado Walmart",
      "amount": "150.0000",
      "amount_in_usd": "150.0000",
      "exchange_rate": "1.000000",
      "transaction_date": "2026-02-19T10:00:00Z",
      "notes": "Compra semanal",
      "is_reconciled": false,
      "created_at": "2026-02-19T10:05:00Z",
      "updated_at": "2026-02-19T10:05:00Z",
      "account_from": {
        "id": 1,
        "name": "BBVA USD",
        "type": "BANK",
        "balance": "4850.0000",
        "currency": { "id": 1, "code": "USD", "symbol": "$" }
      },
      "account_to": null,
      "category": {
        "id": 3,
        "name": "Alimentación",
        "type": "EXPENSE"
      }
    }
  ],
  "total": 45,
  "page": 1,
  "page_size": 20,
  "total_pages": 3
}
```

---

### 6.2 Obtener transacción por ID

**`GET /transactions/:id`** 🔒

```bash
curl -X GET http://localhost:8080/api/v1/transactions/10 \
  -H "Authorization: Bearer <access_token>"
```

**Respuesta `200`:**
```json
{
  "data": {
    "id": 10,
    "type": "EXPENSE",
    "description": "Supermercado Walmart",
    "amount": "150.0000",
    "transaction_date": "2026-02-19T10:00:00Z",
    "account_from": { "id": 1, "name": "BBVA USD", "type": "BANK" },
    "category": { "id": 3, "name": "Alimentación", "type": "EXPENSE" }
  }
}
```

---

### 6.3 Crear transacción — EXPENSE (Gasto)

**`POST /transactions`** 🔒

```bash
curl -X POST http://localhost:8080/api/v1/transactions \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "EXPENSE",
    "description": "Supermercado Walmart",
    "amount": "150.00",
    "account_from_id": 1,
    "category_id": 3,
    "transaction_date": "2026-02-19T10:00:00Z",
    "notes": "Compra semanal"
  }'
```

---

### 6.4 Crear transacción — INCOME (Ingreso)

```bash
curl -X POST http://localhost:8080/api/v1/transactions \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "INCOME",
    "description": "Salario febrero",
    "amount": "3500.00",
    "account_from_id": 1,
    "category_id": 1,
    "transaction_date": "2026-02-01T09:00:00Z",
    "notes": "Pago quincenal"
  }'
```

---

### 6.5 Crear transacción — TRANSFER (Transferencia)

```bash
curl -X POST http://localhost:8080/api/v1/transactions \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "TRANSFER",
    "description": "Mover dinero a efectivo",
    "amount": "500.00",
    "account_from_id": 1,
    "account_to_id": 2,
    "transaction_date": "2026-02-19T12:00:00Z"
  }'
```

**Body completo (referencia):**
| Campo | Tipo | Requerido | Descripción |
|-------|------|-----------|-------------|
| `type` | string | ✅ | `INCOME`, `EXPENSE` o `TRANSFER` |
| `description` | string | ✅ | 1–255 caracteres |
| `amount` | decimal string | ✅ | Debe ser positivo |
| `account_from_id` | integer | ✅ | Cuenta origen |
| `account_to_id` | integer | Solo TRANSFER | Cuenta destino |
| `category_id` | integer | Solo INCOME/EXPENSE | Categoría |
| `transaction_date` | string RFC3339 | ✅ | Fecha de la transacción |
| `notes` | string | ❌ | Notas adicionales (máx 1000) |
| `exchange_rate` | decimal string | ❌ | Tasa de cambio (default: `"1"`) |

**Respuesta `201`:**
```json
{
  "message": "Transaction created successfully",
  "data": {
    "id": 11,
    "type": "EXPENSE",
    "description": "Supermercado Walmart",
    "amount": "150.0000",
    "amount_in_usd": "150.0000",
    "exchange_rate": "1.000000",
    "transaction_date": "2026-02-19T10:00:00Z",
    "account_from": { "id": 1, "name": "BBVA USD", "type": "BANK" },
    "category": { "id": 3, "name": "Alimentación", "type": "EXPENSE" },
    "is_reconciled": false
  }
}
```

---

### 6.6 Actualizar transacción

**`PUT /transactions/:id`** 🔒

> Solo se pueden actualizar campos descriptivos. Para cambiar monto o cuentas, elimina y crea una nueva.

```bash
curl -X PUT http://localhost:8080/api/v1/transactions/10 \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Walmart — compra mensual",
    "notes": "Actualizado",
    "is_reconciled": true
  }'
```

**Body (todos opcionales):**
| Campo | Tipo | Descripción |
|-------|------|-------------|
| `description` | string | Nueva descripción |
| `notes` | string | Notas (máx 1000) |
| `transaction_date` | string RFC3339 | Nueva fecha |
| `is_reconciled` | boolean | Marcar como conciliada |

**Respuesta `200`:**
```json
{
  "message": "Transaction updated successfully"
}
```

---

### 6.7 Eliminar transacción (reversión contable)

**`DELETE /transactions/:id`** 🔒

> No borra físicamente. Crea **asientos de reversión** (DEBIT↔CREDIT invertidos) y restaura los saldos de las cuentas.

```bash
curl -X DELETE http://localhost:8080/api/v1/transactions/10 \
  -H "Authorization: Bearer <access_token>"
```

**Respuesta `200`:**
```json
{
  "message": "Transaction deleted successfully (reversed)"
}
```

---

## 7. Módulo: Categories

Categorías para clasificar ingresos y gastos. Todos los endpoints son 🔒.

---

### 7.1 Listar categorías

**`GET /categories`** 🔒

**Query params:**
| Param | Tipo | Descripción |
|-------|------|-------------|
| `type` | string | `INCOME` o `EXPENSE` |

```bash
# Todas las categorías
curl -X GET http://localhost:8080/api/v1/categories \
  -H "Authorization: Bearer <access_token>"

# Solo categorías de gastos
curl -X GET "http://localhost:8080/api/v1/categories?type=EXPENSE" \
  -H "Authorization: Bearer <access_token>"
```

**Respuesta `200`:**
```json
{
  "data": [
    {
      "id": 1,
      "name": "Salario",
      "type": "INCOME",
      "user_id": 1,
      "is_active": true
    },
    {
      "id": 3,
      "name": "Alimentación",
      "type": "EXPENSE",
      "user_id": 1,
      "is_active": true
    }
  ],
  "count": 2
}
```

---

### 7.2 Obtener categoría por ID

**`GET /categories/:id`** 🔒

```bash
curl -X GET http://localhost:8080/api/v1/categories/3 \
  -H "Authorization: Bearer <access_token>"
```

**Respuesta `200`:**
```json
{
  "data": {
    "id": 3,
    "name": "Alimentación",
    "type": "EXPENSE",
    "user_id": 1,
    "is_active": true
  }
}
```

---

### 7.3 Crear categoría

**`POST /categories`** 🔒

```bash
curl -X POST http://localhost:8080/api/v1/categories \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Transporte",
    "type": "EXPENSE"
  }'
```

**Body:**
| Campo | Tipo | Requerido | Descripción |
|-------|------|-----------|-------------|
| `name` | string | ✅ | 2–100 caracteres |
| `type` | string | ✅ | `INCOME` o `EXPENSE` |

**Respuesta `201`:**
```json
{
  "message": "Category created successfully",
  "data": {
    "id": 5,
    "name": "Transporte",
    "type": "EXPENSE",
    "user_id": 1,
    "is_active": true
  }
}
```

---

### 7.4 Actualizar categoría

**`PUT /categories/:id`** 🔒

```bash
curl -X PUT http://localhost:8080/api/v1/categories/5 \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Transporte y Gasolina",
    "is_active": true
  }'
```

**Body (todos opcionales):**
| Campo | Tipo | Descripción |
|-------|------|-------------|
| `name` | string | Nuevo nombre |
| `is_active` | boolean | Activar/desactivar |

**Respuesta `200`:**
```json
{
  "message": "Category updated successfully"
}
```

---

### 7.5 Eliminar categoría

**`DELETE /categories/:id`** 🔒 — Soft delete. No elimina las transacciones asociadas.

```bash
curl -X DELETE http://localhost:8080/api/v1/categories/5 \
  -H "Authorization: Bearer <access_token>"
```

**Respuesta `200`:**
```json
{
  "message": "Category deleted successfully"
}
```

---

## 8. Módulo: Dashboard

⭐ **Feature estrella.** Proporciona el resumen financiero completo incluyendo el cálculo de **Runway** (cuántos meses puede sobrevivir el usuario sin ingresos). Todos los endpoints son 🔒.

---

### 8.1 Dashboard completo

**`GET /dashboard`** 🔒

```bash
curl -X GET http://localhost:8080/api/v1/dashboard \
  -H "Authorization: Bearer <access_token>"
```

**Respuesta `200`:**
```json
{
  "data": {
    "total_assets": "10000.0000",
    "total_liabilities": "2000.0000",
    "net_worth": "8000.0000",
    "liquid_assets": "7000.0000",
    "monthly_income": "3500.0000",
    "monthly_expenses": "1800.0000",
    "monthly_net_cash_flow": "1700.0000",
    "runway": 3.88,
    "runway_days": 116,
    "average_monthly_expenses": "1800.0000",
    "account_balances": [
      {
        "id": 1,
        "name": "BBVA USD",
        "type": "BANK",
        "balance": "5000.0000",
        "currency_code": "USD",
        "currency_symbol": "$",
        "is_active": true
      },
      {
        "id": 2,
        "name": "Efectivo",
        "type": "CASH",
        "balance": "2000.0000",
        "currency_code": "USD",
        "currency_symbol": "$",
        "is_active": true
      }
    ],
    "as_of": "2026-02-19T14:35:40Z",
    "base_currency": "USD"
  }
}
```

---

### 8.2 Cálculo detallado de Runway

**`GET /dashboard/runway`** 🔒

Calcula cuántos meses puede sostenerse el usuario sin nuevos ingresos.

**Fórmula:**
```
Runway = (Activos Líquidos - Pasivos a Corto Plazo) / Promedio Gastos Mensuales (últimos 3 meses)
```

**Estados:**
| Status | Condición | Significado |
|--------|-----------|-------------|
| `HEALTHY` | ≥ 6 meses | Situación financiera estable |
| `WARNING` | 3–6 meses | Atención recomendada |
| `CRITICAL` | < 3 meses | Acción inmediata necesaria |

```bash
curl -X GET http://localhost:8080/api/v1/dashboard/runway \
  -H "Authorization: Bearer <access_token>"
```

**Respuesta `200`:**
```json
{
  "data": {
    "liquid_assets": "7000.0000",
    "short_term_liabilities": "2000.0000",
    "available_funds": "5000.0000",
    "average_monthly_expenses": "1800.0000",
    "runway_months": 2.77,
    "runway_days": 83,
    "calculation_date": "2026-02-19T14:35:40Z",
    "base_currency": "USD",
    "bank_accounts": [
      {
        "id": 1,
        "name": "BBVA USD",
        "type": "BANK",
        "balance": "5000.0000",
        "currency_code": "USD",
        "currency_symbol": "$",
        "is_active": true
      }
    ],
    "cash_accounts": [
      {
        "id": 2,
        "name": "Efectivo",
        "type": "CASH",
        "balance": "2000.0000",
        "currency_code": "USD",
        "currency_symbol": "$",
        "is_active": true
      }
    ],
    "credit_card_accounts": [
      {
        "id": 3,
        "name": "Visa Platinum",
        "type": "CREDIT_CARD",
        "balance": "2000.0000",
        "currency_code": "USD",
        "currency_symbol": "$",
        "is_active": true
      }
    ],
    "status": "CRITICAL",
    "message": "⚠️ CRITICAL: Your runway is less than 3 months. Consider increasing income or reducing expenses immediately."
  }
}
```

---

### 8.3 Estadísticas mensuales

**`GET /dashboard/monthly-stats`** 🔒

**Query params (opcionales):**
| Param | Tipo | Descripción |
|-------|------|-------------|
| `month` | integer | Mes 1–12 (default: mes actual) |
| `year` | integer | Año (default: año actual) |

```bash
# Mes actual
curl -X GET http://localhost:8080/api/v1/dashboard/monthly-stats \
  -H "Authorization: Bearer <access_token>"

# Mes específico
curl -X GET "http://localhost:8080/api/v1/dashboard/monthly-stats?month=1&year=2026" \
  -H "Authorization: Bearer <access_token>"
```

**Respuesta `200`:**
```json
{
  "data": {
    "month": 2,
    "year": 2026,
    "income": "3500.0000",
    "expenses": "1800.0000",
    "net_cash_flow": "1700.0000",
    "transaction_count": 23
  }
}
```

---

## 9. Módulo: Currencies

Catálogo de monedas disponibles. **Endpoints públicos — no requieren token.**

---

### 9.1 Listar monedas

**`GET /currencies`**

```bash
curl -X GET http://localhost:8080/api/v1/currencies
```

**Respuesta `200`:**
```json
{
  "data": [
    {
      "id": 1,
      "code": "USD",
      "name": "US Dollar",
      "symbol": "$",
      "is_active": true,
      "created_at": "2026-01-01T00:00:00Z",
      "updated_at": "2026-01-01T00:00:00Z"
    },
    {
      "id": 2,
      "code": "EUR",
      "name": "Euro",
      "symbol": "€",
      "is_active": true,
      "created_at": "2026-01-01T00:00:00Z",
      "updated_at": "2026-01-01T00:00:00Z"
    },
    {
      "id": 3,
      "code": "MXN",
      "name": "Mexican Peso",
      "symbol": "$",
      "is_active": true,
      "created_at": "2026-01-01T00:00:00Z",
      "updated_at": "2026-01-01T00:00:00Z"
    }
  ],
  "count": 3
}
```

---

### 9.2 Obtener moneda por código ISO

**`GET /currencies/:code`**

```bash
curl -X GET http://localhost:8080/api/v1/currencies/USD
```

**Respuesta `200`:**
```json
{
  "data": {
    "id": 1,
    "code": "USD",
    "name": "US Dollar",
    "symbol": "$",
    "is_active": true,
    "created_at": "2026-01-01T00:00:00Z",
    "updated_at": "2026-01-01T00:00:00Z"
  }
}
```

---

## 10. Módulo: System Values

Catálogos del sistema con los valores válidos para los distintos campos enumerados. **Endpoints públicos — no requieren token.** Úsalos para poblar selects/dropdowns en la app.

---

### 10.1 Tipos de cuenta

**`GET /system-values/account-types`**

```bash
curl -X GET http://localhost:8080/api/v1/system-values/account-types
```

**Respuesta `200`:**
```json
{
  "data": ["BANK", "CASH", "CREDIT_CARD", "SAVINGS", "INVESTMENT"],
  "count": 5
}
```

---

### 10.2 Clasificaciones de cuenta

**`GET /system-values/account-classifications`**

```bash
curl -X GET http://localhost:8080/api/v1/system-values/account-classifications
```

**Respuesta `200`:**
```json
{
  "data": ["ASSET", "LIABILITY", "EQUITY"],
  "count": 3
}
```

---

### 10.3 Tipos de transacción

**`GET /system-values/transaction-types`**

```bash
curl -X GET http://localhost:8080/api/v1/system-values/transaction-types
```

**Respuesta `200`:**
```json
{
  "data": ["INCOME", "EXPENSE", "TRANSFER"],
  "count": 3
}
```

---

### 10.4 Tipos de categoría

**`GET /system-values/category-types`**

```bash
curl -X GET http://localhost:8080/api/v1/system-values/category-types
```

**Respuesta `200`:**
```json
{
  "data": ["INCOME", "EXPENSE"],
  "count": 2
}
```

---

### 10.5 Catálogo genérico por tipo

**`GET /system-values/catalog/:catalogType`**

Obtiene cualquier catálogo por su nombre.

```bash
curl -X GET http://localhost:8080/api/v1/system-values/catalog/ACCOUNT_TYPE
curl -X GET http://localhost:8080/api/v1/system-values/catalog/TRANSACTION_TYPE
curl -X GET http://localhost:8080/api/v1/system-values/catalog/CATEGORY_TYPE
curl -X GET http://localhost:8080/api/v1/system-values/catalog/ACCOUNT_CLASSIFICATION
```

---

## 11. Módulo: Journal Entries

Pista de auditoría contable. Los asientos son **generados automáticamente** por el Motor Contable al crear o revertir transacciones. **Solo lectura** — nunca se crean ni modifican manualmente.

Todos los endpoints son 🔒.

---

### 11.1 Listar asientos contables

**`GET /journal-entries`** 🔒

**Query params (opcionales):**
| Param | Tipo | Descripción |
|-------|------|-------------|
| `transaction_id` | integer | Filtrar por transacción |
| `account_id` | integer | Filtrar por cuenta |
| `debit_or_credit` | string | `DEBIT` o `CREDIT` |
| `page` | integer | Número de página (default: 1) |
| `page_size` | integer | Elementos por página (default: 50, máx: 100) |

```bash
# Todos los asientos
curl -X GET http://localhost:8080/api/v1/journal-entries \
  -H "Authorization: Bearer <access_token>"

# Asientos de una cuenta específica
curl -X GET "http://localhost:8080/api/v1/journal-entries?account_id=1&page=1" \
  -H "Authorization: Bearer <access_token>"

# Solo débitos
curl -X GET "http://localhost:8080/api/v1/journal-entries?debit_or_credit=DEBIT" \
  -H "Authorization: Bearer <access_token>"
```

**Respuesta `200`:**
```json
{
  "entries": [
    {
      "id": 21,
      "user_id": 1,
      "transaction_id": 10,
      "account_id": 3,
      "debit_or_credit": "DEBIT",
      "amount": "150.0000",
      "entry_date": "2026-02-19T10:00:00Z",
      "description": "Expense: Supermercado Walmart",
      "created_at": "2026-02-19T10:05:00Z"
    },
    {
      "id": 22,
      "user_id": 1,
      "transaction_id": 10,
      "account_id": 1,
      "debit_or_credit": "CREDIT",
      "amount": "150.0000",
      "entry_date": "2026-02-19T10:00:00Z",
      "description": "Payment: Supermercado Walmart",
      "created_at": "2026-02-19T10:05:00Z"
    }
  ],
  "total": 88,
  "page": 1,
  "page_size": 50,
  "total_pages": 2
}
```

---

### 11.2 Asientos de una transacción

**`GET /journal-entries/transaction/:id`** 🔒

Obtiene los 2 asientos (DEBIT + CREDIT) generados para una transacción específica.

```bash
curl -X GET http://localhost:8080/api/v1/journal-entries/transaction/10 \
  -H "Authorization: Bearer <access_token>"
```

**Respuesta `200`:**
```json
{
  "data": [
    {
      "id": 21,
      "transaction_id": 10,
      "account_id": 3,
      "debit_or_credit": "DEBIT",
      "amount": "150.0000",
      "entry_date": "2026-02-19T10:00:00Z",
      "description": "Expense: Supermercado Walmart"
    },
    {
      "id": 22,
      "transaction_id": 10,
      "account_id": 1,
      "debit_or_credit": "CREDIT",
      "amount": "150.0000",
      "entry_date": "2026-02-19T10:00:00Z",
      "description": "Payment: Supermercado Walmart"
    }
  ],
  "count": 2
}
```

---

### 11.3 Verificar equilibrio contable de una transacción

**`GET /journal-entries/verify/:id`** 🔒

Verifica que `∑ Débitos = ∑ Créditos` para una transacción. Útil para auditoría.

```bash
curl -X GET http://localhost:8080/api/v1/journal-entries/verify/10 \
  -H "Authorization: Bearer <access_token>"
```

**Respuesta `200`:**
```json
{
  "data": {
    "transaction_id": 10,
    "is_balanced": true,
    "total_debits": "150.0000",
    "total_credits": "150.0000",
    "difference": "0.0000",
    "message": "Transaction is balanced"
  }
}
```

---

## 12. Catálogo de valores

Referencia rápida de todos los valores enumerados que usa la API.

### Tipos de cuenta (`account_type`)

| Valor | Descripción | Tipo contable |
|-------|-------------|---------------|
| `BANK` | Cuenta bancaria | Activo líquido |
| `CASH` | Efectivo | Activo líquido |
| `CREDIT_CARD` | Tarjeta de crédito | Pasivo |
| `SAVINGS` | Cuenta de ahorro | Activo |
| `INVESTMENT` | Inversión | Activo |

### Tipos de transacción (`type`)

| Valor | `category_id` | `account_to_id` | Efecto contable |
|-------|:---:|:---:|-----------------|
| `INCOME` | ✅ Requerido | ❌ | DEBIT cuenta → CREDIT categoría |
| `EXPENSE` | ✅ Requerido | ❌ | DEBIT categoría → CREDIT cuenta |
| `TRANSFER` | ❌ | ✅ Requerido | DEBIT cuenta destino → CREDIT cuenta origen |

### Tipos de categoría (`type`)

| Valor | Uso |
|-------|-----|
| `INCOME` | Categorías para clasificar ingresos |
| `EXPENSE` | Categorías para clasificar gastos |

### Estados de Runway (`status`)

| Valor | Condición | Color sugerido |
|-------|-----------|----------------|
| `HEALTHY` | ≥ 6 meses | 🟢 Verde |
| `WARNING` | 3–6 meses | 🟡 Amarillo |
| `CRITICAL` | < 3 meses | 🔴 Rojo |

---

## 13. Códigos de error

| Código HTTP | Significado | Causa común |
|-------------|-------------|-------------|
| `400` | Bad Request | Body inválido, campo requerido faltante, regla de negocio violada |
| `401` | Unauthorized | Token ausente, expirado o inválido |
| `404` | Not Found | Recurso no existe o fue eliminado (soft delete) |
| `409` | Conflict | Email duplicado en registro |
| `500` | Internal Server Error | Error inesperado del servidor |

### Estructura de error estándar

```json
{
  "error": "Failed to create transaction",
  "details": "category_id is required for EXPENSE transactions"
}
```

---

### Flujo de autenticación sugerido para la app móvil

```
1. App abre → verificar si hay access_token guardado
2. Si hay token → intentar request
   ├── 200 OK → continuar normalmente
   └── 401 Unauthorized → ir al paso 3
3. Intentar refresh con refresh_token guardado
   ├── 200 OK → guardar nuevos tokens → reintentar request original
   └── 401 Unauthorized → redirect a pantalla de Login
4. Login exitoso → guardar access_token y refresh_token de forma segura
   (Keychain en iOS, EncryptedSharedPreferences en Android)
```

---

> **Documentación interactiva completa:** `http://localhost:8080/docs`  
> **Generada con:** swaggo/swag — actualizar con `make swagger` tras cualquier cambio en DTOs o anotaciones.