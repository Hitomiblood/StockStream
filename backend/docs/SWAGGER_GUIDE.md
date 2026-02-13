# 📚 Documentación de API con Swagger

¡Tu API ahora tiene documentación interactiva completa con Swagger UI!

## 🎯 Swagger UI

Una vez que inicies el servidor backend, podrás acceder a la documentación interactiva en:

```
http://localhost:8080/swagger/index.html
```

## ✨ Características de Swagger UI

### Interfaz Interactiva
- 📖 **Documentación completa**: Todos los endpoints documentados con descripciones detalladas
- 🧪 **Prueba directa**: Ejecuta peticiones directamente desde el navegador
- 📝 **Esquemas de datos**: Ve los modelos de request/response
- 💡 **Ejemplos incluidos**: Cada endpoint tiene ejemplos de uso

### Endpoints Documentados

#### 🏥 Health
- `GET /health` - Verificar estado de la API

#### 📊 Stocks
- `GET /api/v1/stocks` - Listar stocks con paginación
- `GET /api/v1/stocks/{id}` - Obtener stock por ID
- `GET /api/v1/stocks/ticker/{ticker}` - Historial de un ticker
- `GET /api/v1/stocks/latest` - Últimos stocks añadidos
- `GET /api/v1/stocks/search` - Buscar stocks
- `GET /api/v1/stocks/filter` - Filtrar por action/rating
- `POST /api/v1/stocks/fetch` - Sincronizar desde API externa

#### 🎯 Recommendations
- `GET /api/v1/recommendations` - Obtener recomendaciones de inversión

#### ⚙️ Metadata
- `GET /api/v1/metadata` - Obtener metadata de filtros disponibles

## 🚀 Cómo Usar Swagger UI

### 1. Iniciar el Servidor

```powershell
cd backend
go run cmd/api/main.go
```

Verás en los logs:
```
📚 Swagger UI: http://localhost:8080/swagger/index.html
```

### 2. Abrir Swagger UI

Abre tu navegador y ve a: `http://localhost:8080/swagger/index.html`

### 3. Explorar la API

1. **Ver endpoints**: Todos los endpoints están organizados por tags (stocks, recommendations, health)
2. **Expandir endpoint**: Click en cualquier endpoint para ver detalles
3. **Ver parámetros**: Ve qué parámetros requiere cada endpoint
4. **Ver respuestas**: Ve ejemplos de respuestas exitosas y de error

### 4. Probar Endpoints

1. Click en **"Try it out"** en cualquier endpoint
2. Completa los parámetros requeridos
3. Click en **"Execute"**
4. Ve la respuesta en la sección "Response"

#### Ejemplo: Buscar Stocks

1. Expandir `GET /api/v1/stocks/search`
2. Click en "Try it out"
3. En el campo `q`, escribe: `Apple`
4. Click en "Execute"
5. Ve los resultados en la sección Response

#### Ejemplo: Sincronizar Datos

1. Expandir `POST /api/v1/stocks/fetch`
2. Click en "Try it out"
3. Click en "Execute" (no requiere parámetros)
4. Espera la respuesta (puede tardar 1-2 minutos)

## 🔧 Regenerar Documentación

Si haces cambios en las anotaciones de Swagger en el código, regenera la documentación:

```powershell
cd backend
swag init -g cmd/api/main.go -o docs
```

## 📖 Anotaciones de Swagger

Las anotaciones de Swagger se agregan como comentarios especiales sobre cada handler:

```go
// GetAllStocks maneja GET /api/v1/stocks
// @Summary      List stocks
// @Description  Get all stocks with pagination, sorting and filtering
// @Tags         stocks
// @Accept       json
// @Produce      json
// @Param        limit   query  int  false  "Number of results"
// @Success      200  {object}  map[string]interface{}
// @Router       /api/v1/stocks [get]
func (h *StockHandler) GetAllStocks(c *gin.Context) {
    // ...
}
```

### Anotaciones Principales

- `@Summary`: Título corto del endpoint
- `@Description`: Descripción detallada
- `@Tags`: Agrupa endpoints relacionados
- `@Accept`: Tipo de contenido que acepta (json, xml, etc.)
- `@Produce`: Tipo de contenido que devuelve
- `@Param`: Parámetros del endpoint
  - Formato: `name location type required "description"`
  - Locations: `query`, `path`, `header`, `body`
- `@Success`: Respuestas exitosas (código y tipo)
- `@Failure`: Respuestas de error
- `@Router`: Ruta y método HTTP

## 🎨 Personalización

### Información General de la API

La información general se configura en `cmd/api/main.go`:

```go
// @title           Stock Analysis System API
// @version         1.0
// @description     Sistema completo de análisis de acciones...
// @host            localhost:8080
// @BasePath        /
```

### Tags Personalizados

Los tags se usan para agrupar endpoints en la UI:

```go
// @Tags stocks
// @Tags recommendations
// @Tags health
```

## 📊 Archivos Generados

El comando `swag init` genera estos archivos en `backend/docs/`:

- `docs.go` - Código Go con la documentación
- `swagger.json` - Especificación OpenAPI en JSON
- `swagger.yaml` - Especificación OpenAPI en YAML

**⚠️ No edites estos archivos manualmente**, se regeneran automáticamente.

## 🔗 Especificación OpenAPI

Puedes acceder a la especificación OpenAPI en formato JSON:

```
http://localhost:8080/swagger/doc.json
```

Esto es útil para:
- Importar en herramientas como Postman
- Generar clientes automáticamente
- Integrar con otras herramientas de API

## 💡 Tips y Mejores Prácticas

### 1. Documentación Clara
- Usa descripciones concisas pero completas
- Incluye ejemplos en los parámetros
- Documenta todos los posibles códigos de error

### 2. Modelos de Datos
- Define structs con tags JSON para documentación automática
- Usa anotaciones de validación para documentar restricciones

### 3. Agrupación Lógica
- Usa tags para organizar endpoints relacionados
- Mantén una estructura consistente

### 4. Mantén Actualizado
- Regenera la documentación después de cambios
- Prueba los endpoints en Swagger UI antes de deploy

## 🆚 Swagger vs Postman

| Característica | Swagger UI | Postman |
|----------------|------------|---------|
| **Documentación** | ✅ Automática desde código | ⚠️ Manual |
| **Tests** | ✅ Pruebas básicas | ✅ Tests avanzados |
| **Interfaz** | ✅ Web integrada | ✅ Aplicación dedicada |
| **Compartir** | ✅ URL pública | ⚠️ Requiere cuenta |
| **Versionado** | ✅ Con el código | ⚠️ Separado |
| **Colaboración** | ✅ Automática | ✅ Con workspace |

**Recomendación**: Usa ambos:
- **Swagger** para documentación y pruebas rápidas
- **Postman** para testing exhaustivo y colecciones

## 🎓 Recursos

### Documentación Oficial
- [Swag Documentation](https://github.com/swaggo/swag)
- [Gin-Swagger](https://github.com/swaggo/gin-swagger)
- [OpenAPI Specification](https://swagger.io/specification/)

### Tutoriales
- [Swagger with Gin Tutorial](https://github.com/swaggo/swag#getting-started)
- [OpenAPI Best Practices](https://swagger.io/docs/specification/api-design-best-practices/)

## ❓ Troubleshooting

### Error: "docs package not found"

**Solución:**
```powershell
cd backend
swag init -g cmd/api/main.go -o docs
```

### Error: Swagger UI no carga

**Verificar:**
1. El servidor está corriendo
2. La ruta es correcta: `/swagger/index.html`
3. Los archivos docs fueron generados

### Error: Endpoint no aparece en Swagger

**Verificar:**
1. Las anotaciones están correctas
2. Regeneraste la documentación
3. Reiniciaste el servidor

## 🎉 Ventajas de Swagger en tu Proyecto

✅ **Para el desarrollador:**
- Documentación siempre actualizada
- Pruebas rápidas durante el desarrollo
- Menos tiempo escribiendo docs manuales

✅ **Para el equipo:**
- Punto de referencia centralizado
- Fácil de compartir (solo una URL)
- Tests básicos disponibles para todos

✅ **Para la entrevista:**
- Demuestra profesionalismo
- Facilita la demo del proyecto
- Muestra conocimiento de buenas prácticas

---

**¡Tu API ahora está completamente documentada! 📚**

Accede a: `http://localhost:8080/swagger/index.html`
