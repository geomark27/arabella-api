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

-   **Motor:** **SQL Server 2022** (Express o Web Edition en AWS RDS).

-   **Features Clave:**

    -   **Temporal Tables:** Para auditoría histórica automática
        > (SYSTEM_VERSIONING = ON).

    -   **Constraints:** Para asegurar SUM(Debit) = SUM(Credit).

    -   **Tipos de Datos:** DECIMAL(19,4) para precisión monetaria
        > absoluta.

### Infraestructura (AWS - Serverless First)

-   **Compute:** AWS Lambda (funciones Go) para la API y workers.

-   **API Gateway:** Entrada para el Frontend.

-   **Ingesta:** **AWS SES (Simple Email Service)** -\> Regla de
    > recepción -\> S3 Bucket -\> Lambda Trigger (Parser).

-   **Almacenamiento:** S3 (para almacenar los correos crudos o recibos
    > escaneados).

## 5. Definición del MVP (Producto Mínimo Viable)

El MVP se centrará en la **solidez del dato**. No buscaremos la
automatización total en el día 1, sino la integridad total.

### Alcance Funcional MVP

1.  **Gestión de Cuentas:** Crear cuentas bancarias (USD/Local),
    > efectivo y tarjetas.

2.  **Registro Manual Optimizado:** Interfaz rápida para registrar
    > gastos/ingresos.

3.  **Motor de Doble Partida:** El backend procesa todo como asientos
    > contables.

4.  **Dashboard Básico:** Saldo actual real, Gastos del mes, Deuda
    > total.

5.  **Multi-moneda Básico:** Registro manual de la tasa de cambio al
    > momento de la transacción.

*Nota: La automatización por correo (AWS SES) quedará maquetada en
arquitectura pero se implementará en la Fase 2.*

## 6. Fases y Tiempos Estimados (Roadmap)

Considerando un dedicación \"Side Project\" (10-15 horas semanales).

### Fase 1: Fundamentos y Arquitectura (Semanas 1-3)

-   Diseño del Esquema de Base de Datos (ERD) en SQL Server.

-   Configuración del repo (Go + Next.js).

-   Implementación del \"Core\" en Go (Lógica de asientos contables y
    > validaciones).

-   Despliegue inicial de infraestructura \"Hello World\" en AWS
    > (Terraform/CDK).

### Fase 2: API y Lógica de Negocio (Semanas 4-6)

-   Endpoints CRUD para Cuentas y Transacciones.

-   Implementación de lógica de Tarjetas de Crédito y Deudas.

-   Tests unitarios del motor contable (Crítico).

### Fase 3: Frontend y UX (Semanas 7-10)

-   Desarrollo de la UI en Next.js.

-   Integración con la API.

-   Dashboards y Gráficos (Recharts o Tremor).

-   PWA Setup (Service Workers).

### Fase 4: Automatización AWS (Semanas 11-13)

-   Configuración de AWS SES.

-   Desarrollo de la Lambda Parser (Regex/Lógica de extracción de
    > emails).

-   Integración del flujo de \"Transacciones Pendientes de Aprobar\".

**Tiempo Total Estimado al MVP:** 3 a 3.5 meses.

## 7. Posibles Áreas de Mejora y Riesgos

1.  **Fricción de Usuario:** Si el registro manual es lento, el usuario
    > abandona.

    -   *Mitigación:* UI optimizada \"One-thumb\" (uso con una mano).

2.  **Variabilidad de Emails:** Los bancos cambian formatos de correo.

    -   *Mitigación:* Arquitectura de parsers modulares (Strategy
        > Pattern) fácil de actualizar sin redeployar todo el backend.

3.  **Costos de AWS:** SQL Server en RDS puede ser caro si no se cuida
    > la capa gratuita o instancias reservadas.

    -   *Mitigación:* Iniciar con SQL Server Express en una EC2 pequeña
        > o Dockerizado, migrar a RDS solo al escalar.

## 8. Siguiente Paso Sugerido

Para arrancar con el pie derecho, el siguiente paso lógico es definir el
**Modelo de Datos**.