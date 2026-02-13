# 📮 Guía Completa de Postman - Stock Analysis System

## 🎯 Objetivo
Esta guía te ayudará a probar tanto la API externa del challenge como tu propia API una vez la construyas.

---

## 📥 Instalación de Postman

### Windows
```powershell
# Usando winget
winget install Postman.Postman

# O descarga directamente desde:
# https://www.postman.com/downloads/
```

---

## 🔐 Parte 1: Testear la API Externa del Challenge

### Paso 1: Crear una Nueva Request en Postman

1. **Abrir Postman**
2. **Crear nueva colección**:
   - Click en "New" → "Collection"
   - Nombre: `Stock Challenge API`
   - Descripción: `Pruebas de API externa del challenge`

3. **Crear nueva request**:
   - Click derecho en la colección → "Add request"
   - Nombre: `Get Stocks - First Page`

### Paso 2: Configurar la Request

#### Configuración Básica
- **Method**: `GET`
- **URL**: `https://api.karenai.click/swechallenge/list`

#### Headers
Agregar los siguientes headers:

| Key | Value |
|-----|-------|
| `Authorization` | `Bearer your_token_here` |
| `Content-Type` | `application/json` |

**📸 Captura de pantalla de referencia:**
```
Tab: Headers
┌─────────────────┬──────────────────────────────────────────┐
│ KEY             │ VALUE                                    │
├─────────────────┼──────────────────────────────────────────┤
│ Authorization   │ Bearer eyJhbGciOiJIUzI1NiIsInR5cC...    │
│ Content-Type    │ application/json                         │
└─────────────────┴──────────────────────────────────────────┘
```

### Paso 3: Enviar la Request

1. Click en el botón azul **"Send"**
2. Verificar que el **Status** sea: `200 OK`
3. Ver la respuesta en formato JSON

#### Respuesta Esperada
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
    },
    // ... más items
  ],
  "next_page": "ZBIO"
}
```

### Paso 4: Testear Paginación

Para obtener la siguiente página de resultados:

1. **Duplicar la request anterior**:
   - Click derecho en la request → "Duplicate"
   - Nombre: `Get Stocks - Next Page`

2. **Agregar Query Parameter**:
   - Tab: `Params`
   - Agregar:

| Key | Value | Description |
|-----|-------|-------------|
| `next_page` | `ZBIO` | Token de la página anterior |

**URL resultante**:
```
https://api.karenai.click/swechallenge/list?next_page=ZBIO
```

3. **Enviar** y verificar nuevos resultados

---

## 🔧 Parte 2: Testear tu Propia API (Backend)

### Configurar Variables de Entorno en Postman

1. **Crear Environment**:
   - Click en "Environments" (ícono de ojo)
   - Click en "Add"
   - Nombre: `Local Development`

2. **Agregar Variables**:

| Variable | Initial Value | Current Value |
|----------|--------------|---------------|
| `base_url` | `http://localhost:8080` | `http://localhost:8080` |
| `api_version` | `v1` | `v1` |

3. **Activar Environment**:
   - Seleccionar "Local Development" en el dropdown

### Request 1: Health Check

**Propósito**: Verificar que el servidor está corriendo

- **Method**: `GET`
- **URL**: `{{base_url}}/health`
- **Headers**: Ninguno especial

**Respuesta Esperada**:
```json
{
  "status": "ok",
  "message": "Stock Analysis System API"
}
```

---

### Request 2: Obtener Todas las Acciones

**Propósito**: Listar todas las acciones almacenadas en la DB

- **Method**: `GET`
- **URL**: `{{base_url}}/api/{{api_version}}/stocks`
- **Headers**:
  - `Content-Type: application/json`

**Query Parameters (opcionales)**:
| Key | Value | Description |
|-----|-------|-------------|
| `limit` | `10` | Número de resultados |
| `offset` | `0` | Offset para paginación |
| `sort` | `time` | Campo para ordenar |
| `order` | `desc` | Dirección del ordenamiento |

**Respuesta Esperada**:
```json
{
  "data": [
    {
      "id": 1,
      "ticker": "AAPL",
      "company": "Apple Inc.",
      "target_from": "$150.00",
      "target_to": "$160.00",
      "action": "target raised by",
      "rating_from": "Buy",
      "rating_to": "Strong Buy",
      "time": "2025-11-14T00:30:06Z"
    }
  ],
  "total": 100,
  "limit": 10,
  "offset": 0
}
```

---

### Request 3: Buscar Acciones

**Propósito**: Buscar acciones por ticker o nombre de compañía

- **Method**: `GET`
- **URL**: `{{base_url}}/api/{{api_version}}/stocks/search`

**Query Parameters**:
| Key | Value | Description |
|-----|-------|-------------|
| `q` | `Apple` | Término de búsqueda |

**Respuesta Esperada**:
```json
{
  "data": [
    {
      "id": 1,
      "ticker": "AAPL",
      "company": "Apple Inc.",
      ...
    }
  ],
  "total": 1
}
```

---

### Request 4: Obtener una Acción Específica

**Propósito**: Obtener detalles de una acción por ID

- **Method**: `GET`
- **URL**: `{{base_url}}/api/{{api_version}}/stocks/1`
- **Headers**:
  - `Content-Type: application/json`

**Respuesta Esperada**:
```json
{
  "id": 1,
  "ticker": "AAPL",
  "company": "Apple Inc.",
  "target_from": "$150.00",
  "target_to": "$160.00",
  "action": "target raised by",
  "brokerage": "Goldman Sachs",
  "rating_from": "Buy",
  "rating_to": "Strong Buy",
  "time": "2025-11-14T00:30:06Z",
  "created_at": "2025-11-15T10:00:00Z",
  "updated_at": "2025-11-15T10:00:00Z"
}
```

---

### Request 5: Obtener Acciones por Ticker

**Propósito**: Obtener historial de una acción específica

- **Method**: `GET`
- **URL**: `{{base_url}}/api/{{api_version}}/stocks/ticker/AAPL`

**Respuesta Esperada**:
```json
{
  "ticker": "AAPL",
  "company": "Apple Inc.",
  "history": [
    {
      "id": 1,
      "target_from": "$150.00",
      "target_to": "$160.00",
      "action": "target raised by",
      "time": "2025-11-14T00:30:06Z"
    },
    {
      "id": 2,
      "target_from": "$145.00",
      "target_to": "$150.00",
      "action": "target raised by",
      "time": "2025-11-10T00:30:06Z"
    }
  ],
  "total": 2
}
```

---

### Request 6: Forzar Actualización de Datos

**Propósito**: Traer datos frescos de la API externa

- **Method**: `POST`
- **URL**: `{{base_url}}/api/{{api_version}}/stocks/fetch`
- **Headers**:
  - `Content-Type: application/json`
- **Body**: (ninguno)

**Respuesta Esperada**:
```json
{
  "message": "Data fetched successfully",
  "total_fetched": 145,
  "total_new": 23,
  "total_updated": 122,
  "duration_ms": 3450
}
```

---

### Request 7: Obtener Recomendaciones

**Propósito**: Obtener las mejores acciones para invertir hoy

- **Method**: `GET`
- **URL**: `{{base_url}}/api/{{api_version}}/recommendations`

**Query Parameters (opcionales)**:
| Key | Value | Description |
|-----|-------|-------------|
| `limit` | `5` | Número de recomendaciones |

**Respuesta Esperada**:
```json
{
  "recommendations": [
    {
      "stock": {
        "id": 1,
        "ticker": "AAPL",
        "company": "Apple Inc.",
        "target_from": "$150.00",
        "target_to": "$160.00",
        "action": "target raised by",
        "rating_from": "Buy",
        "rating_to": "Strong Buy"
      },
      "score": 92.5,
      "reason": "Strong upward trend with multiple target raises. Rating upgraded to Strong Buy.",
      "confidence": "high"
    },
    {
      "stock": {
        "ticker": "GOOGL",
        "company": "Alphabet Inc.",
        ...
      },
      "score": 87.3,
      "reason": "Consistent positive ratings and target increases.",
      "confidence": "high"
    }
  ],
  "generated_at": "2025-11-15T10:00:00Z",
  "criteria": {
    "target_increase_weight": 0.4,
    "rating_upgrade_weight": 0.3,
    "recent_activity_weight": 0.3
  }
}
```

---

### Request 8: Filtrar por Acción

**Propósito**: Filtrar acciones por tipo de acción (target raised/lowered)

- **Method**: `GET`
- **URL**: `{{base_url}}/api/{{api_version}}/stocks`

**Query Parameters**:
| Key | Value | Description |
|-----|-------|-------------|
| `action` | `target raised by` | Tipo de acción |

---

### Request 9: Filtrar por Rating

**Propósito**: Filtrar acciones por rating actual

- **Method**: `GET`
- **URL**: `{{base_url}}/api/{{api_version}}/stocks`

**Query Parameters**:
| Key | Value | Description |
|-----|-------|-------------|
| `rating` | `Buy` | Rating actual |

---

## 📊 Tests Automáticos en Postman

### Agregar Tests a tus Requests

En el tab "Tests" de cada request, puedes agregar scripts para validar automáticamente:

#### Test 1: Verificar Status Code
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});
```

#### Test 2: Verificar Estructura de Respuesta
```javascript
pm.test("Response has required fields", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData).to.have.property('items');
    pm.expect(jsonData).to.have.property('next_page');
});
```

#### Test 3: Verificar que hay Datos
```javascript
pm.test("Response contains data", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.items).to.be.an('array');
    pm.expect(jsonData.items.length).to.be.above(0);
});
```

#### Test 4: Guardar next_page para siguiente request
```javascript
pm.test("Save next_page token", function () {
    var jsonData = pm.response.json();
    if (jsonData.next_page) {
        pm.environment.set("next_page_token", jsonData.next_page);
    }
});
```

#### Test 5: Verificar Tiempo de Respuesta
```javascript
pm.test("Response time is less than 2000ms", function () {
    pm.expect(pm.response.responseTime).to.be.below(2000);
});
```

---

## 🔄 Collection Runner

### Ejecutar Múltiples Requests Automáticamente

1. **Click derecho en tu colección** → "Run collection"
2. **Seleccionar requests** a ejecutar
3. **Configurar**:
   - Iterations: 1
   - Delay: 0ms
4. **Click "Run"**

Esto ejecutará todas las requests en secuencia y mostrará un resumen de tests.

---

## 📁 Exportar/Importar Colección

### Exportar
1. Click derecho en colección → "Export"
2. Guardar como `Stock_Challenge_API.postman_collection.json`

### Importar
1. Click en "Import"
2. Seleccionar archivo JSON
3. La colección aparecerá en tu workspace

---

## 🎨 Organizar Colecciones

### Estructura Recomendada

```
📁 Stock Challenge API
├── 📁 External API
│   ├── Get Stocks - First Page
│   └── Get Stocks - Next Page
├── 📁 Health & Status
│   └── Health Check
├── 📁 Stocks CRUD
│   ├── Get All Stocks
│   ├── Get Stock by ID
│   ├── Get Stock by Ticker
│   ├── Search Stocks
│   └── Fetch Fresh Data
├── 📁 Recommendations
│   └── Get Recommendations
└── 📁 Filters
    ├── Filter by Action
    └── Filter by Rating
```

---

## 🐛 Troubleshooting

### Error: "Could not send request"
- Verificar que el servidor está corriendo
- Verificar que la URL es correcta
- Verificar firewall/antivirus

### Error: 401 Unauthorized
- Verificar que el token en Authorization header es correcto
- Verificar que incluiste "Bearer " antes del token

### Error: 404 Not Found
- Verificar que la ruta es correcta
- Verificar que el endpoint está implementado en tu backend

### Error: 500 Internal Server Error
- Revisar logs del servidor backend
- Verificar que la base de datos está corriendo
- Verificar que la configuración (.env) es correcta

---

## 💡 Tips Profesionales

### 1. Usar Pre-request Scripts
Ejecutar código antes de cada request:
```javascript
// Agregar timestamp
pm.environment.set("timestamp", new Date().toISOString());

// Generar token temporal
pm.environment.set("request_id", Math.random().toString(36).substring(7));
```

### 2. Usar Snippets
Postman tiene snippets predefinidos en el tab "Tests":
- Status code tests
- Response body tests
- Response time tests

### 3. Documentar tus Requests
- Agregar descripción a cada request
- Incluir ejemplos de respuestas
- Documentar parámetros requeridos/opcionales

### 4. Crear Ambientes para Diferentes Stages
```
📁 Environments
├── Local Development (localhost:8080)
├── Staging (staging.tudominio.com)
└── Production (api.tudominio.com)
```

---

## 📚 Recursos Adicionales

- [Postman Learning Center](https://learning.postman.com/)
- [Postman API Testing Guide](https://www.postman.com/api-platform/api-testing/)
- [Postman Scripts Examples](https://www.postman.com/postman/workspace/postman-answers/collection/9215231-b9133e48-73c3-4aa4-b189-e038ee4c5e00)

---

## ✅ Checklist de Testing

### API Externa
- [ ] GET primera página funciona
- [ ] GET con paginación funciona
- [ ] Headers de autorización correctos
- [ ] Response tiene estructura esperada

### Tu API (Backend)
- [ ] Health check responde
- [ ] GET /stocks lista acciones
- [ ] GET /stocks/:id obtiene detalles
- [ ] GET /stocks/ticker/:ticker funciona
- [ ] POST /stocks/fetch actualiza datos
- [ ] GET /recommendations devuelve sugerencias
- [ ] Búsqueda funciona correctamente
- [ ] Filtros funcionan
- [ ] Ordenamiento funciona
- [ ] Paginación funciona

---

**¡Happy Testing! 🚀**
