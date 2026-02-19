# Stock Analysis System

Sistema completo de análisis de acciones con backend en Golang, frontend en Vue3 y base de datos CockroachDB.

## 🎉 Estado del Proyecto

- ✅ **Backend**: **100% COMPLETADO Y FUNCIONAL**
- ✅ **Frontend**: MVP implementado y funcional
- ✅ **Tests**: Implementados (backend con Go test + frontend con Vitest)

## 🚀 Quick Start

### Iniciar el Backend

```powershell
# 1. Iniciar CockroachDB
docker start cockroachdb
# O si no existe:
docker-compose up -d cockroachdb
 
# 2. Ejecutar migraciones
cd backend
go run ./cmd/migrate up

# 3. Iniciar servidor backend
go run cmd/api/main.go
 
# 4. Probar la API
Invoke-WebRequest -Uri http://localhost:8080/health
```

### Iniciar el Frontend

```powershell
cd frontend
Copy-Item .env.example .env
npm install
npm run dev
```

Frontend por defecto en:

```HTTP
http://localhost:5173
```

### 🧱 Migraciones SQL (CLI)

```powershell
cd backend
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
migrate -version

# Crear nueva migración
migrate create -ext sql -dir migrations -seq add_portfolios_table

# Aplicar migraciones con el runner del proyecto
go run ./cmd/migrate/main.go up
```

> Configura el esquema en `.env` con `DB_SCHEMA` (ejemplo: `public`).

### 📚 Swagger UI - Documentación Interactiva

**¡Prueba la API desde tu navegador sin necesidad de Postman!**

Una vez iniciado el servidor, accede a:
```
http://localhost:8080/swagger/index.html
```

Swagger UI te permite:
- 📖 Ver todos los endpoints con descripciones detalladas
- 🧪 Ejecutar peticiones directamente desde el navegador
- 📝 Ver modelos de datos y ejemplos
- 💡 Aprender cómo usar la API de forma interactiva

📚 **Guía completa**: [SWAGGER_GUIDE.md](backend/docs/SWAGGER_GUIDE.md)

### 📚 Documentación Completa

| Documento | Descripción | Estado |
|-----------|-------------|--------|
| **[BACKEND_COMPLETADO.md](backend\docs\BACKEND_COMPLETADO.md)** | ⭐ Resumen completo del backend implementado | ✅ |
| **[backend/README.md](backend/README.md)** | Guía de uso y testing del backend | ✅ |
| **[frontend/README.md](frontend/README.md)** | Guía del frontend (rutas, stores y scripts) | ✅ |
| **[backend/SWAGGER_GUIDE.md](backend\docs\SWAGGER_GUIDE.md)** | 📚 Documentación interactiva con Swagger UI | ✅ |
| **[GETTING_STARTED.md](backend\docs\GETTING_STARTED.md)** | Guía paso a paso desde cero | ✅ |
| **[POSTMAN_GUIDE.md](backend\docs\POSTMAN_GUIDE.md)** | Testing con Postman | ✅ |
| **[PLAN_DE_ACTIVIDADES.md](backend\docs\PLAN_DE_ACTIVIDADES.md)** | Plan de desarrollo completo | ✅ |

## 📊 Estructura del Proyecto

```
stock-analysis-system/
├── backend/                        ✅ COMPLETADO
│   ├── cmd/api/main.go            ✅ Entry point configurado
│   ├── internal/
│   │   ├── config/                ✅ Gestión de configuración
│   │   ├── database/              ✅ Conexión a CockroachDB
│   │   ├── models/                ✅ Modelos de datos
│   │   ├── services/              ✅ Lógica de negocio
│   │   ├── handlers/              ✅ HTTP handlers (10 endpoints)
│   │   └── middleware/            ✅ CORS configurado
│   ├── cmd/*_test.go              ✅ Tests unitarios/componente
│   ├── internal/*_test.go         ✅ Tests unitarios/componente
│   └── README.md                  ✅ Documentación completa
├── frontend/                       ✅ MVP implementado + tests
├── GETTING_STARTED.md             ✅ Guía completa
├── POSTMAN_GUIDE.md               ✅ Testing guide
├── PLAN_DE_ACTIVIDADES.md         ✅ Plan de trabajo
├── BACKEND_COMPLETADO.md          ✅ Resumen del backend
└── docker-compose.yml             ✅ CockroachDB configurado
```

## 🎯 API Endpoints Disponibles

| Endpoint | Método | Descripción |
|----------|--------|-------------|
| `/health` | GET | Health check |
| `/api/v1/stocks` | GET | Listar stocks con paginación |
| `/api/v1/stocks/latest` | GET | Últimos stocks |
| `/api/v1/stocks/search` | GET | Buscar stocks |
| `/api/v1/stocks/filter` | GET | Filtrar por action/rating |
| `/api/v1/stocks/ticker/:ticker` | GET | Historial por ticker |
| `/api/v1/stocks/:id` | GET | Obtener por ID |
| `/api/v1/stocks/fetch` | POST | Sincronizar desde API externa |
| `/api/v1/recommendations` | GET | Recomendaciones de inversión ⭐ |
| `/api/v1/metadata` | GET | Metadata (filtros disponibles) |

**Ver [backend/README.md](backend/README.md) para detalles y ejemplos de uso.**

## 📋 Tecnologías

### Backend (Implementado ✅)
- **Golang 1.25.6**: Lenguaje principal
- **Gin 1.11.0**: Framework web
- **GORM 1.31.1**: ORM para base de datos
- **CockroachDB**: Base de datos distribuida
- **Resty**: Cliente HTTP para API externa
- **godotenv**: Gestión de variables de entorno

### Frontend (Implementado ✅)
- **Vue 3**: Framework de UI
- **TypeScript**: Tipado estático
- **Pinia**: State management
- **Tailwind CSS**: Framework de CSS
- **Vite**: Build tool
- **Vitest + Vue Test Utils**: Testing unitario de API, stores, componentes y vistas
- Ver `.env.example` para la configuración necesaria

### Backend Testing (Implementado ✅)

```powershell
cd backend
go test ./...
```

### Validación Frontend ✅

```powershell
cd frontend
npm run typecheck
npm run test
npm run lint
npm run build
```

## 📝 Licencia

Proyecto de reto técnico - Uso educativo
