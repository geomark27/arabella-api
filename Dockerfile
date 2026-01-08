# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Instalar dependencias del sistema necesarias
RUN apk add --no-cache git ca-certificates tzdata

# Copiar archivos de dependencias primero (para cache de Docker)
COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto del código
COPY . .

# Compilar la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o arabella-api ./cmd/arabella-api/main.go

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Instalar certificados CA y timezone data
RUN apk --no-cache add ca-certificates tzdata

# Crear usuario no-root para seguridad
RUN adduser -D -g '' appuser

# Copiar el binario compilado
COPY --from=builder /app/arabella-api .

# Copiar archivos necesarios
COPY --from=builder /app/docs ./docs

# Cambiar al usuario no-root
USER appuser

# Exponer el puerto
EXPOSE 8080

# Comando de inicio
CMD ["./arabella-api"]
