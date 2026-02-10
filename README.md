# Stock Analysis System

Sistema completo de análisis de acciones con backend en Golang, frontend en Vue3 y base de datos CockroachDB.

## 🎉 Estado del Proyecto

- ✅ **Backend**: **100% COMPLETADO Y FUNCIONAL**
- ⏳ **Frontend**: Pendiente
- ⏳ **Tests**: Pendiente

## 🚀 Quick Start

### Iniciar el Backend

```powershell
# 1. Iniciar CockroachDB
docker start cockroachdb
# O si no existe:
docker-compose up -d cockroachdb

# 2. Iniciar servidor backend
cd backend
go run cmd/api/main.go

# 3. Probar la API
Invoke-WebRequest -Uri http://localhost:8080/health
```

### 📚 Documentación Completa

| Documento | Descripción | Estado |
|-----------|-------------|--------|
| **[BACKEND_COMPLETADO.md](Documentación\BACKEND_COMPLETADO.md)** | ⭐ Resumen completo del backend implementado | ✅ |
| **[backend/README.md](backend/README.md)** | Guía de uso y testing del backend | ✅ |
| **[GETTING_STARTED.md](Documentación\GETTING_STARTED.md)** | Guía paso a paso desde cero | ✅ |
| **[POSTMAN_GUIDE.md](Documentación\POSTMAN_GUIDE.md)** | Testing con Postman | ✅ |
| **[PLAN_DE_ACTIVIDADES.md](Documentación\PLAN_DE_ACTIVIDADES.md)** | Plan de desarrollo completo | ✅ |

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
│   ├── tests/                     ⏳ Pendiente
│   └── README.md                  ✅ Documentación completa
├── frontend/                       ⏳ Pendiente
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

### Frontend (Pendiente ⏳)
- **Vue 3**: Framework de UI
- **TypeScript**: Tipado estático
- **Pinia**: State management
- **Tailwind CSS**: Framework de CSS
- **Vite**: Build tool
- Ver `.env.example` para la configuración necesaria

## 📝 Licencia

Proyecto de reto técnico - Uso educativo
