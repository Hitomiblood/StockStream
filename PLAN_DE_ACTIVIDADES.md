# 📋 Plan de Actividades - Stock Analysis System

## 🎯 Resumen Ejecutivo

**Objetivo**: Desarrollar un sistema completo de análisis de acciones  
**Duración Estimada**: 10 días  
**Enfoque**: Backend First (ideal para desarrolladores backend)  
**Nivel**: Principiante en las tecnologías del stack

---

## 📅 Cronograma Detallado

### **DÍA 1: Configuración del Entorno** (4-6 horas)

#### Mañana (2-3 horas)
- [x] Instalar Golang
- [x] Instalar Docker Desktop
- [x] Instalar Postman
- [x] Instalar VS Code + extensiones (Go, Git)
- [x] Verificar instalaciones

#### Tarde (2-3 horas)
- [x] Configurar CockroachDB con Docker
- [x] Crear estructura de directorios del proyecto
- [x] Inicializar repositorio Git
- [x] Crear archivos de configuración (.env, .gitignore)
- [x] Hacer primer commit

**Entregable**: Entorno completamente configurado y funcionando

---

### **DÍA 2: Fundamentos de Golang + Primera Conexión** (6-8 horas)

#### Mañana (3-4 horas)
- [x] Tutorial básico de Go ([Tour of Go](https://go.dev/tour/))
  - Variables y tipos
  - Estructuras de control
  - Funciones
  - Structs e interfaces
- [X] Leer sobre Go modules y packages

#### Tarde (3-4 horas)
- [X] Crear proyecto Go con `go mod init`
- [X] Instalar dependencias (Gin, GORM, godotenv)
- [X] Crear archivo `config.go` para leer variables de entorno
- [X] Crear `main.go` con endpoint `/health`
- [X] Probar endpoint con Postman

**Entregable**: API básica funcionando con endpoint de salud

**Código de referencia**: Ver `GETTING_STARTED.md` Paso 4 y 7

---

### **DÍA 3: Conexión a DB + Modelos** (6-8 horas)

#### Mañana (3-4 horas)
- [ ] Aprender básicos de GORM
- [ ] Crear archivo `database/db.go`
- [ ] Implementar función de conexión a CockroachDB
- [ ] Probar conexión (crear tabla de prueba)

#### Tarde (3-4 horas)
- [X] Diseñar modelo `Stock` en `models/stock.go`
- [ ] Implementar migraciones automáticas con GORM
- [ ] Crear operaciones CRUD básicas (Create, Read)
- [ ] Probar inserción manual de datos

**Entregable**: Conexión a DB funcionando con modelo Stock

**Recursos**:
- [GORM Documentation](https://gorm.io/docs/)
- [CockroachDB + Go Tutorial](https://www.cockroachlabs.com/docs/stable/build-a-go-app-with-cockroachdb-gorm.html)

---

### **DÍA 4: Cliente API Externa + Almacenamiento** (6-8 horas)

#### Mañana (3-4 horas)
- [ ] Crear `services/api_client.go`
- [ ] Implementar función para hacer GET a API externa
- [ ] Manejar autenticación (Bearer token)
- [ ] Parsear respuesta JSON a struct `APIResponse`
- [ ] Probar con Postman que la API externa funciona

#### Tarde (3-4 horas)
- [ ] Implementar lógica de paginación (usar `next_page`)
- [ ] Crear función para traer TODOS los datos (loop de páginas)
- [ ] Guardar datos en CockroachDB
- [ ] Evitar duplicados (check por ticker + time)
- [ ] Crear endpoint POST `/stocks/fetch` que ejecute la sincronización

**Entregable**: Sistema que consume API externa y guarda en DB

**Testing**:
```bash
# Probar fetch
curl -X POST http://localhost:8080/api/v1/stocks/fetch

# Verificar en DB
docker exec -it cockroachdb ./cockroach sql --insecure -e "SELECT COUNT(*) FROM stockdb.stocks;"
```

---

### **DÍA 5: API REST Completa** (6-8 horas)

#### Mañana (3-4 horas)
- [ ] Crear `handlers/stock_handlers.go`
- [ ] Implementar endpoint GET `/stocks` (listar todos)
- [ ] Agregar paginación (limit, offset)
- [ ] Implementar endpoint GET `/stocks/:id` (detalle)
- [ ] Probar con Postman

#### Tarde (3-4 horas)
- [ ] Implementar búsqueda GET `/stocks/search?q=Apple`
- [ ] Implementar filtros (por action, rating)
- [ ] Implementar ordenamiento (sort, order)
- [ ] Agregar middleware CORS
- [ ] Documentar todos los endpoints en Postman

**Entregable**: API REST completa y funcional

**Endpoints a implementar**:
- `GET /api/v1/stocks` - Listar con paginación
- `GET /api/v1/stocks/:id` - Detalle
- `GET /api/v1/stocks/ticker/:ticker` - Historial por ticker
- `GET /api/v1/stocks/search` - Búsqueda
- `POST /api/v1/stocks/fetch` - Sincronizar datos

---

### **DÍA 6: Algoritmo de Recomendación** (6-8 horas)

#### Mañana (3-4 horas)
- [ ] Diseñar criterios de scoring:
  - Target aumenta = +puntos
  - Target disminuye = -puntos
  - Rating upgrade (Buy → Strong Buy) = +puntos
  - Rating downgrade = -puntos
  - Actividad reciente = +puntos
- [ ] Crear `services/recommendation.go`
- [ ] Implementar función de cálculo de score

#### Tarde (3-4 horas)
- [ ] Implementar endpoint GET `/recommendations`
- [ ] Agregar parámetro `limit` para top N
- [ ] Incluir explicación del por qué se recomienda
- [ ] Probar con datos reales
- [ ] Ajustar pesos del algoritmo según resultados

**Entregable**: Sistema de recomendaciones funcionando

**Ejemplo de algoritmo básico**:
```go
score := 0.0

// Target increase
if targetTo > targetFrom {
    score += 30
}

// Rating upgrade
if ratingTo > ratingFrom {  // Strong Buy > Buy > Neutral
    score += 40
}

// Recent activity (últimos 7 días)
if time.Since(stock.Time) < 7*24*time.Hour {
    score += 30
}
```

---

### **DÍA 7: Testing Backend** (4-6 horas)

#### Mañana (2-3 horas)
- [ ] Aprender básicos de testing en Go
- [ ] Crear tests unitarios para `api_client.go`
- [ ] Crear tests para `recommendation.go`
- [ ] Usar mocks para DB

#### Tarde (2-3 horas)
- [ ] Tests de integración para handlers
- [ ] Tests de conexión a DB
- [ ] Ejecutar `go test -cover ./...`
- [ ] Objetivo: >70% de cobertura

**Entregable**: Suite de tests con buena cobertura

**Recursos**:
- [Go Testing Tutorial](https://go.dev/doc/tutorial/add-a-test)
- [Testify Documentation](https://github.com/stretchr/testify)

---

### **DÍA 8: Introducción a Frontend - Setup** (4-6 horas)

#### Mañana (2-3 horas)
- [ ] Tutorial básico de Vue 3 (1-2 horas)
- [ ] Tutorial básico de TypeScript (1 hora)
- [ ] Leer sobre Composition API

#### Tarde (2-3 horas)
- [ ] Crear proyecto Vue 3 con Vite
- [ ] Configurar TypeScript
- [ ] Instalar y configurar Tailwind CSS
- [ ] Instalar Pinia
- [ ] Crear estructura base de directorios

**Comandos**:
```bash
cd frontend
npm create vite@latest . -- --template vue-ts
npm install
npm install -D tailwindcss postcss autoprefixer
npm install pinia axios
npx tailwindcss init -p
```

**Entregable**: Proyecto Vue 3 configurado

---

### **DÍA 9: Frontend - UI Básica** (6-8 horas)

#### Mañana (3-4 horas)
- [ ] Crear servicio API client (`services/api.ts`)
- [ ] Crear Pinia store (`stores/stockStore.ts`)
- [ ] Crear componente `StockTable.vue`
- [ ] Mostrar lista de acciones desde el backend

#### Tarde (3-4 horas)
- [ ] Crear componente `StockFilters.vue`
- [ ] Agregar búsqueda en tiempo real
- [ ] Agregar filtros por action y rating
- [ ] Estilizar con Tailwind CSS

**Entregable**: UI básica funcional con tabla y filtros

---

### **DÍA 10: Frontend - Recomendaciones + Pulido** (6-8 horas)

#### Mañana (3-4 horas)
- [ ] Crear componente `RecommendationPanel.vue`
- [ ] Mostrar top 5 recomendaciones
- [ ] Agregar indicadores visuales (score, confidence)
- [ ] Crear vista de detalles de una acción

#### Tarde (3-4 horas)
- [ ] Pulir diseño general
- [ ] Agregar loading states
- [ ] Agregar manejo de errores
- [ ] Probar integración completa
- [ ] Preparar documentación final

**Entregable**: Sistema completo funcionando

---

## 📊 Hitos Principales

| Hito | Día | Validación |
|------|-----|------------|
| ✅ Entorno configurado | 1 | Docker + Go + DB funcionando |
| ✅ API básica | 2 | Endpoint /health responde |
| ✅ DB conectada | 3 | Datos se guardan correctamente |
| ✅ Sync API externa | 4 | POST /fetch trae datos |
| ✅ API REST completa | 5 | Todos los endpoints funcionan |
| ✅ Recomendaciones | 6 | GET /recommendations funciona |
| ✅ Tests backend | 7 | >70% cobertura |
| ✅ Frontend setup | 8 | npm run dev funciona |
| ✅ UI básica | 9 | Se ven datos en navegador |
| ✅ Sistema completo | 10 | Todo integrado |

---

## 🎓 Recursos de Aprendizaje por Día

### Golang
- **Día 1-2**: [Tour of Go](https://go.dev/tour/) (3-4 horas)
- **Día 2**: [Go by Example](https://gobyexample.com/)
- **Día 3**: [GORM Getting Started](https://gorm.io/docs/)
- **Día 4**: [Go HTTP Client Tutorial](https://gobyexample.com/http-clients)
- **Día 5**: [Gin Framework Tutorial](https://gin-gonic.com/docs/quickstart/)

### Frontend
- **Día 8**: [Vue 3 Tutorial](https://vuejs.org/tutorial/) (1-2 horas)
- **Día 8**: [TypeScript in 5 minutes](https://www.typescriptlang.org/docs/handbook/typescript-in-5-minutes.html)
- **Día 9**: [Pinia Tutorial](https://pinia.vuejs.org/getting-started.html)
- **Día 9**: [Tailwind CSS Tutorial](https://tailwindcss.com/docs/utility-first)

---

## ⚡ Consejos para Principiantes

### Si nunca has usado Golang:
1. **No te preocupes**: Go es uno de los lenguajes más fáciles de aprender
2. **Dedica 3-4 horas al Tour of Go**: Vale la pena
3. **Copia y adapta ejemplos**: No trates de memorizarlo todo
4. **Usa ChatGPT/Copilot**: Para entender errores de compilación

### Si nunca has usado Vue:
1. **Start simple**: Un componente a la vez
2. **Vue 3 es más fácil que Vue 2**: No te confundas con tutoriales antiguos
3. **Composition API**: Es la forma moderna, úsala
4. **Tailwind facilita el CSS**: No necesitas ser experto en CSS

### Si nunca has usado CockroachDB:
1. **Es como PostgreSQL**: Si sabes SQL, ya sabes 90%
2. **Usa Docker**: Es la forma más fácil de empezar
3. **GORM lo abstrae**: No necesitas conocer detalles internos

---

## 🚨 Posibles Bloqueos y Soluciones

### Bloqueo 1: "No sé cómo empezar con Go"
**Solución**: 
- Dedica 3 horas al Tour of Go (obligatorio)
- Copia el código de ejemplo del `GETTING_STARTED.md`
- Ejecuta `go run cmd/api/main.go` y ve qué pasa
- Los errores te enseñarán

### Bloqueo 2: "No puedo conectar a CockroachDB"
**Solución**:
```bash
# Verificar que Docker está corriendo
docker ps

# Verificar logs
docker logs cockroachdb

# Probar conexión manual
docker exec -it cockroachdb ./cockroach sql --insecure

# Si todo falla, usa CockroachDB Cloud (gratis)
```

### Bloqueo 3: "La API externa no responde"
**Solución**:
- Verificar token en Postman primero
- Verificar que incluyes "Bearer " antes del token
- Ver si el token expiró (revisar campo "exp" en el JWT)
- Contactar a quien dio el challenge

### Bloqueo 4: "Los tests no pasan"
**Solución**:
- Empieza con tests simples
- No te obsesiones con 100% de cobertura
- Usa `go test -v` para ver detalles
- Pregunta a ChatGPT sobre el error específico

### Bloqueo 5: "Frontend no se comunica con Backend"
**Solución**:
- Verificar CORS en backend
- Verificar que backend está corriendo
- Usar DevTools del navegador → Network tab
- Verificar URL del API en el código frontend

---

## ✅ Checklist Final de Entrega

### Código
- [ ] Backend completamente funcional
- [ ] Frontend completamente funcional
- [ ] Tests escritos y pasando
- [ ] .env.example creado (sin credenciales)
- [ ] .gitignore configurado
- [ ] Sin archivos sensibles en Git

### Documentación
- [ ] README.md completo
- [ ] Instrucciones de instalación
- [ ] Instrucciones de ejecución
- [ ] Endpoints documentados
- [ ] Screenshots de la UI (opcional pero impresiona)

### Demo
- [ ] Script de presentación preparado
- [ ] Postman collection exportada
- [ ] Datos de ejemplo en la DB
- [ ] Sistema corriendo sin errores

---

## 🎯 Objetivos de Aprendizaje

Al final de este proyecto habrás aprendido:

### Backend
- ✅ Crear una API REST en Golang
- ✅ Usar Gin framework
- ✅ Conectar a una base de datos con GORM
- ✅ Consumir APIs externas
- ✅ Manejar paginación
- ✅ Escribir tests en Go
- ✅ Estructurar un proyecto Go

### Frontend
- ✅ Crear aplicaciones Vue 3 con TypeScript
- ✅ Usar Composition API
- ✅ Manejar estado con Pinia
- ✅ Estilizar con Tailwind CSS
- ✅ Consumir APIs REST

### DevOps
- ✅ Usar Docker para desarrollo
- ✅ Configurar bases de datos
- ✅ Manejar variables de entorno
- ✅ Usar Git apropiadamente

### Soft Skills
- ✅ Leer documentación técnica
- ✅ Debuggear problemas
- ✅ Estructurar un proyecto desde cero
- ✅ Presentar tu trabajo

---

## 💪 Motivación

Este reto puede parecer intimidante, especialmente si es tu primera vez con estas tecnologías. Pero:

1. **Es totalmente factible**: Miles de desarrolladores han aprendido estas tecnologías desde cero
2. **La documentación es excelente**: Todas estas tecnologías tienen docs de primera calidad
3. **La comunidad es grande**: Stack Overflow tiene respuestas para casi todo
4. **Es una inversión**: Estas tecnologías son muy demandadas en el mercado

**Consejo final**: No trates de ser perfecto. Enfócate en que funcione primero, luego puedes mejorar el código. ¡Tú puedes! 🚀

---

**¡Mucho éxito en tu reto técnico!**
