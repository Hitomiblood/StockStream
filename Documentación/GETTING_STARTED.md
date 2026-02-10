# Guía de Inicio - Stock Analysis System

## 📋 Descripción General

Sistema completo para análisis de acciones que incluye:
- **Backend en Golang**: API REST que consume datos externos y los almacena en CockroachDB
- **Frontend en Vue3**: Interfaz de usuario con TypeScript, Pinia y Tailwind CSS
- **Base de datos**: CockroachDB
- **Recomendaciones**: Algoritmo de inversión basado en análisis de datos

---

## 🎯 Objetivos del Challenge

1. ✅ Conectar a la API y almacenar datos en CockroachDB
2. ✅ Crear API REST y UI intuitiva con búsqueda, ordenamiento y filtros
3. ✅ Implementar algoritmo de recomendación de acciones
4. ✅ Escribir tests unitarios

---

## 📊 Estructura de Datos de la API Externa

**Endpoint**: `GET https://api.karenai.click/swechallenge/list`

**Response**:
```json
{
  "items": [
    {
      "ticker": "RMTI",
      "target_from": "$3.00",
      "target_to": "$2.50",
      "company": "Rockwell Medical",
      "action": "target lowered by",
      "brokerage": "",
      "rating_from": "Buy",
      "rating_to": "Buy",
      "time": "2025-11-14T00:30:06.716890597Z"
    }
  ],
  "next_page": "ZBIO"
}
```

---

## 🚀 Plan de Desarrollo (Enfoque Backend First)

### Fase 1: Configuración del Entorno (Día 1)
1. ✅ Instalar herramientas necesarias
2. ✅ Configurar CockroachDB
3. ✅ Inicializar proyecto Golang
4. ✅ Configurar variables de entorno

### Fase 2: Backend - Core (Días 2-3)
1. ✅ Crear cliente HTTP para consumir API externa
2. ✅ Implementar paginación automática (next_page)
3. ✅ Diseñar modelos de datos
4. ✅ Implementar conexión a CockroachDB
5. ✅ Crear servicio de almacenamiento

### Fase 3: Backend - API REST (Día 4)
1. ✅ Crear endpoints CRUD
2. ✅ Implementar búsqueda y filtros
3. ✅ Implementar ordenamiento
4. ✅ Agregar middleware de autenticación (opcional)

### Fase 4: Backend - Algoritmo de Recomendación (Día 5)
1. ✅ Análisis de datos históricos
2. ✅ Implementar scoring system
3. ✅ Crear endpoint `/recommendations`

### Fase 5: Backend - Testing (Día 6)
1. ✅ Tests unitarios de servicios
2. ✅ Tests de integración con DB
3. ✅ Tests de endpoints

### Fase 6: Frontend - Setup (Día 7)
1. ✅ Crear proyecto Vue3 + TypeScript
2. ✅ Configurar Tailwind CSS
3. ✅ Configurar Pinia stores

### Fase 7: Frontend - UI Components (Días 8-9)
1. ✅ Tabla de acciones
2. ✅ Filtros y búsqueda
3. ✅ Vista de detalles
4. ✅ Panel de recomendaciones

### Fase 8: Integración y Deploy (Día 10)
1. ✅ Integrar frontend con backend
2. ✅ Documentación final
3. ✅ Preparar para presentación

---

## 🛠️ Instalación de Herramientas

### 1. Instalar Golang

**Windows:**
```powershell
# Descargar desde: https://go.dev/dl/
# O usar winget:
winget install GoLang.Go
```

**Verificar instalación:**
```powershell
go version
# Debería mostrar: go version go1.21.x windows/amd64
```

### 2. Instalar CockroachDB

**Opción A: Local (Desarrollo)**
```powershell
# Descargar desde: https://www.cockroachlabs.com/docs/stable/install-cockroachdb-windows.html
# O usar chocolatey:
choco install cockroachdb
```

**Opción B: Docker (Recomendado para desarrollo)**
```powershell
# Instalar Docker Desktop primero
docker pull cockroachdb/cockroach:latest

# Iniciar CockroachDB
docker run -d `
  --name cockroachdb `
  -p 26257:26257 `
  -p 8080:8080 `
  cockroachdb/cockroach:latest start-single-node --insecure
```

**Opción C: CockroachDB Cloud (Recomendado para producción)**
- Registrarse en: https://cockroachlabs.cloud/
- Crear cluster gratuito
- Obtener connection string

### 3. Instalar Node.js (para Frontend)

```powershell
# Descargar desde: https://nodejs.org/
# O usar winget:
winget install OpenJS.NodeJS.LTS

# Verificar
node --version
npm --version
```

### 4. Instalar Postman

```powershell
# Descargar desde: https://www.postman.com/downloads/
# O usar winget:
winget install Postman.Postman
```

### 5. Herramientas Adicionales

```powershell
# Git (si no está instalado)
winget install Git.Git

# VS Code (editor recomendado)
winget install Microsoft.VisualStudioCode

# Extensiones recomendadas para VS Code:
# - Go (golang.go)
# - Vue Language Features (Vue.volar)
# - Tailwind CSS IntelliSense
```

---

## 📁 Estructura del Proyecto

```
stock-analysis-system/
├── backend/
│   ├── cmd/
│   │   └── api/
│   │       └── main.go              # Entry point
│   ├── internal/
│   │   ├── config/
│   │   │   └── config.go            # Configuración
│   │   ├── models/
│   │   │   └── stock.go             # Modelos de datos
│   │   ├── database/
│   │   │   └── db.go                # Conexión DB
│   │   ├── services/
│   │   │   ├── api_client.go        # Cliente API externa
│   │   │   ├── stock_service.go     # Lógica de negocio
│   │   │   └── recommendation.go    # Algoritmo recomendación
│   │   ├── handlers/
│   │   │   └── stock_handlers.go    # HTTP handlers
│   │   └── middleware/
│   │       └── cors.go              # CORS middleware
│   ├── tests/
│   │   ├── unit/
│   │   └── integration/
│   ├── go.mod
│   ├── go.sum
│   ├── .env.example
│   └── README.md
│
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   ├── StockTable.vue
│   │   │   ├── StockFilters.vue
│   │   │   ├── StockDetails.vue
│   │   │   └── RecommendationPanel.vue
│   │   ├── stores/
│   │   │   └── stockStore.ts        # Pinia store
│   │   ├── services/
│   │   │   └── api.ts               # API client
│   │   ├── views/
│   │   │   ├── HomeView.vue
│   │   │   └── RecommendationsView.vue
│   │   ├── App.vue
│   │   └── main.ts
│   ├── package.json
│   ├── vite.config.ts
│   ├── tailwind.config.js
│   └── tsconfig.json
│
├── .gitignore
├── docker-compose.yml               # Para desarrollo
└── README.md
```

---

## 🏁 Por Dónde Empezar: BACKEND FIRST

### Paso 1: Configuración Inicial del Proyecto

```powershell
# Crear estructura de directorios
cd "~\Proyecto"

# Crear directorios del backend
mkdir backend
cd backend
mkdir cmd, internal, tests
mkdir cmd\api
mkdir internal\config, internal\models, internal\database, internal\services, internal\handlers, internal\middleware
mkdir tests\unit, tests\integration

# Inicializar módulo Go
go mod init github.com/tu-usuario/stock-analysis-system
```

### Paso 2: Configurar Variables de Entorno

```powershell
# Crear archivo .env en backend/
cd backend
New-Item -ItemType File -Path .env
```

**Contenido de `.env`:**
```env
# API Externa
EXTERNAL_API_URL=https://api.karenai.click/swechallenge/list
EXTERNAL_API_TOKEN=your_token_here

# CockroachDB
DB_HOST=localhost
DB_PORT=26257
DB_USER=root
DB_PASSWORD=
DB_NAME=stockdb
DB_SSLMODE=disable

# Servidor API
API_PORT=8080
API_HOST=localhost

# Configuración
FETCH_INTERVAL=3600  # Segundos entre actualizaciones
LOG_LEVEL=debug
```

### Paso 3: Instalar Dependencias de Golang

```powershell
cd backend

# Framework web (Gin)
go get -u github.com/gin-gonic/gin

# Driver PostgreSQL (CockroachDB usa protocolo PostgreSQL)
go get -u github.com/lib/pq

# ORM (opcional pero recomendado)
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres

# Variables de entorno
go get -u github.com/joho/godotenv

# Testing
go get -u github.com/stretchr/testify

# HTTP client
go get -u github.com/go-resty/resty/v2
```

### Paso 4: Crear Configuración Base

**Archivo: `internal/config/config.go`**

```go
package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// API Externa
	ExternalAPIURL   string
	ExternalAPIToken string

	// Database
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// Server
	APIPort string
	APIHost string

	// Configuración
	FetchInterval int
	LogLevel      string
}

func Load() *Config {
	// Cargar .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "26257"))
	fetchInterval, _ := strconv.Atoi(getEnv("FETCH_INTERVAL", "3600"))

	return &Config{
		ExternalAPIURL:   getEnv("EXTERNAL_API_URL", ""),
		ExternalAPIToken: getEnv("EXTERNAL_API_TOKEN", ""),
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBPort:           dbPort,
		DBUser:           getEnv("DB_USER", "root"),
		DBPassword:       getEnv("DB_PASSWORD", ""),
		DBName:           getEnv("DB_NAME", "stockdb"),
		DBSSLMode:        getEnv("DB_SSLMODE", "disable"),
		APIPort:          getEnv("API_PORT", "8080"),
		APIHost:          getEnv("API_HOST", "localhost"),
		FetchInterval:    fetchInterval,
		LogLevel:         getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
```

### Paso 5: Crear Modelos de Datos

**Archivo: `internal/models/stock.go`**

```go
package models

import (
	"time"
)

type Stock struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Ticker      string    `gorm:"index;not null" json:"ticker"`
	TargetFrom  string    `json:"target_from"`
	TargetTo    string    `json:"target_to"`
	Company     string    `json:"company"`
	Action      string    `json:"action"`
	Brokerage   string    `json:"brokerage"`
	RatingFrom  string    `json:"rating_from"`
	RatingTo    string    `json:"rating_to"`
	Time        time.Time `gorm:"index" json:"time"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type APIResponse struct {
	Items    []Stock `json:"items"`
	NextPage string  `json:"next_page"`
}

type StockRecommendation struct {
	Stock       Stock   `json:"stock"`
	Score       float64 `json:"score"`
	Reason      string  `json:"reason"`
	Confidence  string  `json:"confidence"` // "high", "medium", "low"
}
```

### Paso 6: Iniciar CockroachDB

```powershell
# Si usas Docker:
docker run -d `
  --name cockroachdb `
  -p 26257:26257 `
  -p 8080:8080 `
  -v cockroach-data:/cockroach/cockroach-data `
  cockroachdb/cockroach:latest start-single-node --insecure

# Verificar que está corriendo
docker ps

# Acceder al cliente SQL
docker exec -it cockroachdb ./cockroach sql --insecure

# Crear base de datos
CREATE DATABASE stockdb;
\q
```

### Paso 7: Primera Prueba - Hello World API

**Archivo: `cmd/api/main.go`**

```go
package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/tu-usuario/stock-analysis-system/internal/config"
)

func main() {
	// Cargar configuración
	cfg := config.Load()

	// Crear router
	r := gin.Default()

	// Ruta de prueba
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"message": "Stock Analysis System API",
		})
	})

	// Iniciar servidor
	addr := cfg.APIHost + ":" + cfg.APIPort
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
```

**Ejecutar:**
```powershell
cd backend
go run cmd/api/main.go
```

**Probar:**
```powershell
# En otra terminal
Invoke-WebRequest -Uri http://localhost:8080/health | Select-Object -ExpandProperty Content
```

---

## 🔄 Próximos Pasos Detallados

### 1. Implementar Cliente API Externa

Ver archivo `internal/services/api_client.go` (se creará en los siguientes pasos)

### 2. Implementar Conexión a Base de Datos

Ver archivo `internal/database/db.go` (se creará en los siguientes pasos)

### 3. Implementar Servicios y Handlers

Ver archivos en `internal/services/` y `internal/handlers/`

---

## 🧪 Testing

### Ejecutar Tests
```powershell
# Todos los tests
go test ./...

# Con cobertura
go test -cover ./...

# Tests específicos
go test ./internal/services/...

# Verbose
go test -v ./...
```

---

## 📚 Recursos de Aprendizaje

### Golang
- [Tour of Go](https://go.dev/tour/) - Tutorial oficial interactivo
- [Go by Example](https://gobyexample.com/) - Ejemplos prácticos
- [Effective Go](https://go.dev/doc/effective_go) - Best practices

### Gin Framework
- [Documentación oficial](https://gin-gonic.com/docs/)
- [Ejemplos](https://github.com/gin-gonic/examples)

### GORM
- [Documentación](https://gorm.io/docs/)
- [Guía de inicio](https://gorm.io/docs/index.html)

### CockroachDB
- [Documentación](https://www.cockroachlabs.com/docs/)
- [Quickstart](https://www.cockroachlabs.com/docs/stable/quickstart.html)

### Vue 3 + TypeScript
- [Vue 3 Docs](https://vuejs.org/)
- [Pinia](https://pinia.vuejs.org/)
- [Tailwind CSS](https://tailwindcss.com/docs)

---

## ⚠️ Consideraciones de Seguridad

1. ✅ **NUNCA subir .env a Git**
2. ✅ Usar `.env.example` como template (sin valores sensibles)
3. ✅ Agregar `.env` al `.gitignore`
4. ✅ Usar variables de entorno en producción
5. ✅ Considerar usar secrets managers para producción

---

## 🐛 Troubleshooting

### Error: "cannot find module"
```powershell
go mod tidy
go mod download
```

### Error conexión a CockroachDB
```powershell
# Verificar que Docker está corriendo
docker ps

# Ver logs de CockroachDB
docker logs cockroachdb

# Reiniciar contenedor
docker restart cockroachdb
```

### Puerto 8080 en uso
```powershell
# Encontrar proceso
netstat -ano | findstr :8080

# Matar proceso (reemplazar PID)
taskkill /PID <PID> /F
```

---

## 📞 Soporte

Este proyecto es parte de un reto técnico. Para dudas sobre el desarrollo:
1. Revisar documentación oficial de cada tecnología
2. Buscar en Stack Overflow
3. Consultar ejemplos en GitHub

---

## ✅ Checklist de Progreso

### Backend
- [x] Configuración inicial
- [ ] Cliente API externa
- [ ] Conexión a CockroachDB
- [ ] Modelos de datos
- [ ] Endpoint GET /stocks
- [ ] Endpoint POST /stocks/fetch
- [ ] Endpoint GET /stocks/:id
- [ ] Búsqueda y filtros
- [ ] Algoritmo de recomendación
- [ ] Tests unitarios
- [ ] Tests de integración

### Frontend
- [ ] Setup Vue 3 + Vite
- [ ] Configurar Tailwind
- [ ] Store de Pinia
- [ ] Componente tabla de stocks
- [ ] Componente filtros
- [ ] Vista de detalles
- [ ] Panel de recomendaciones
- [ ] Integración con API

### General
- [ ] Documentación completa
- [ ] .gitignore configurado
- [ ] README.md actualizado
- [ ] Preparado para demo

---

**¡Éxito en tu reto técnico! 🚀**
