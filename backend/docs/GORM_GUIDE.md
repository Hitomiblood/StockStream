# 📚 Guía Práctica de GORM (Consultas y Métodos Comunes)

Esta guía resume lo esencial de **GORM** (ORM para Go) con foco en **cómo consultar** y **cómo escribir** a la base de datos usando su API.

> Proyecto: Stock Analysis System (CockroachDB/PostgreSQL-compatible)

---

## 0) Migraciones SQL con CLI (recomendado)

En este proyecto las migraciones viven en `backend/migrations/` y se ejecutan con el runner de Go en `cmd/migrate/main.go`.

### Instalar la CLI de migraciones

Usamos la CLI oficial de `golang-migrate`:

```powershell
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

En Windows, normalmente se instala en:

```text
C:\Users\<tu_usuario>\go\bin\migrate.exe
```

Validar instalación:

```powershell
migrate -version
```

### Generar una nueva migración

Desde `backend/`:

```powershell
migrate create -ext sql -dir migrations -seq add_portfolios_table
```

Esto crea dos archivos:

- `migrations/0002_add_portfolios_table.up.sql`
- `migrations/0002_add_portfolios_table.down.sql`

Convención recomendada:

- `up`: cambios hacia adelante (crear tablas, columnas, índices, etc.)
- `down`: rollback exacto del `up`

### Aplicar migraciones en este proyecto

Este backend usa su propio runner (CockroachDB-compatible):

```powershell
go run ./cmd/migrate/main.go up
go run ./cmd/migrate/main.go version
go run ./cmd/migrate/main.go down
go run ./cmd/migrate/main.go steps 1
go run ./cmd/migrate/main.go steps -1
```

### Configuración por `.env`

El runner toma conexión y esquema desde variables de entorno:

- `DB_HOST`
- `DB_PORT`
- `DB_USER`
- `DB_PASSWORD`
- `DB_NAME`
- `DB_SSLMODE`
- `DB_SCHEMA` (ejemplo: `public`)

Si cambias `DB_SCHEMA`, el runner usará ese `search_path` al ejecutar migraciones.

---

## 1) Conceptos clave que debes entender

### ¿Qué es `*gorm.DB`?
- Es el *handle* principal para ejecutar queries.
- Es **concurrent-safe**: puedes reutilizarlo entre requests/goroutines.
- Cada llamada devuelve un nuevo `*gorm.DB` (inmutable en la práctica):
  ```go
  q := db.Where("ticker = ?", "AAPL")
  // db no cambia; q es la nueva cadena de query
  ```

### Resultado de una operación
Casi todas las operaciones retornan un `*gorm.DB` con:
- `result.Error` → error (si existe)
- `result.RowsAffected` → cuántas filas tocó la operación

### Manejo de “no encontrado”
Para lecturas puntuales:
- `First`, `Take`, `Last` → si no encuentra, `result.Error == gorm.ErrRecordNotFound`
- `Find` → si no encuentra, **NO** es error; devuelve slice vacío

---

## 2) Modelo y convenciones (rápido)

GORM asume:
- Tabla pluralizada (configurable)
- Clave primaria `ID`
- Campos timestamps: `CreatedAt`, `UpdatedAt`, `DeletedAt` (soft delete)

Ejemplo típico:
```go
// internal/models/stock.go (ejemplo conceptual)
type Stock struct {
  ID         uint      `gorm:"primaryKey"`
  Ticker     string    `gorm:"index"`
  Time       time.Time `gorm:"index"`
  Company    string
  CreatedAt  time.Time
  UpdatedAt  time.Time
}
```

---

## 3) Lecturas (SELECT)

### `First`, `Take`, `Last`
- `First` ordena por primary key asc.
- `Last` ordena por primary key desc.
- `Take` no aplica orden.

```go
var stock models.Stock
err := db.First(&stock, 10).Error          // by PK
err := db.Where("ticker = ?", "AAPL").First(&stock).Error

// Take: útil cuando el orden no importa
err := db.Where("ticker = ?", "AAPL").Take(&stock).Error

// Last
err := db.Where("ticker = ?", "AAPL").Last(&stock).Error
```

### `Find` (listas)
```go
var stocks []models.Stock
err := db.Where("ticker = ?", "AAPL").Order("time DESC").Find(&stocks).Error
```

### `Select` (columnas específicas)
```go
type Lite struct {
  Ticker string
  Time   time.Time
}

var rows []Lite
err := db.Model(&models.Stock{}).
  Select("ticker, time").
  Where("ticker = ?", "AAPL").
  Order("time DESC").
  Limit(50).
  Find(&rows).Error
```

### `Pluck` (una sola columna)
```go
var tickers []string
err := db.Model(&models.Stock{}).
  Distinct("ticker").
  Pluck("ticker", &tickers).Error
```

### `Count`
```go
var total int64
err := db.Model(&models.Stock{}).Where("action = ?", "upgraded by").Count(&total).Error
```

---

## 4) Filtros (WHERE) y condiciones

### Placeholder seguro (evita SQL injection)
✅ Bien:
```go
db.Where("ticker = ?", userInput)
```
❌ Mal:
```go
db.Where("ticker = '" + userInput + "'")
```

### `Where` con SQL + args
```go
db.Where("ticker = ? AND time = ?", stock.Ticker, stock.Time)
```

### `Where` con struct (solo campos no-zero)
```go
db.Where(&models.Stock{Ticker: "AAPL"}).Find(&stocks)
```

### `Where` con map
```go
db.Where(map[string]any{"ticker": "AAPL", "action": "upgraded by"}).Find(&stocks)
```

### `Or`, `Not`
```go
q := db.Where("ticker = ?", "AAPL").Or("ticker = ?", "MSFT")
q.Find(&stocks)

db.Not("brokerage = ?", "").Find(&stocks)
```

### IN
```go
db.Where("ticker IN ?", []string{"AAPL", "MSFT"}).Find(&stocks)
```

### LIKE / ILIKE (PostgreSQL/Cockroach)
```go
pattern := "%" + query + "%"
// ILIKE = case-insensitive en Postgres/Cockroach
err := db.Where("ticker ILIKE ? OR company ILIKE ?", pattern, pattern).
  Order("time DESC").
  Limit(25).
  Find(&stocks).Error
```

### Rango de fechas
```go
from, to := time.Now().AddDate(0, -1, 0), time.Now()
err := db.Where("time BETWEEN ? AND ?", from, to).
  Order("time DESC").
  Find(&stocks).Error
```

---

## 5) Ordenamiento y paginación

### `Order`
```go
db.Order("time DESC").Find(&stocks)
```

### `Limit` + `Offset`
```go
err := db.Limit(limit).Offset(offset).
  Order("time DESC").
  Find(&stocks).Error
```

> Tip: valida `sortBy`/`order` desde backend antes de concatenar (para evitar inyección por ORDER BY).

---

## 6) Escrituras (INSERT/UPDATE/DELETE)

### Crear (INSERT)
```go
err := db.Create(&stock).Error
```

### Actualizar

#### `Save` (actualiza todo el struct)
- Si `stock.ID` existe → UPDATE
- Si no existe → INSERT
```go
err := db.Save(&stock).Error
```

#### `Updates` (solo campos no-zero si usas struct; o todos si usas map)
```go
err := db.Model(&models.Stock{}).
  Where("id = ?", id).
  Updates(map[string]any{"rating_to": "Buy", "brokerage": "X"}).Error
```

#### `Select` / `Omit` en updates
```go
err := db.Model(&stock).
  Select("rating_to", "brokerage").
  Updates(stock).Error

err := db.Model(&stock).
  Omit("created_at").
  Updates(stock).Error
```

#### `Update` / `UpdateColumn`
- `Update` respeta hooks/updated_at
- `UpdateColumn` no dispara hooks ni actualiza `UpdatedAt`
```go
db.Model(&stock).Update("rating_to", "Buy")
db.Model(&stock).UpdateColumn("rating_to", "Buy")
```

### Borrar (DELETE)

#### Soft delete (si el modelo tiene `gorm.DeletedAt`)
```go
db.Delete(&stock)
```

#### Hard delete
```go
db.Unscoped().Delete(&stock)
```

---

## 7) Upsert (insert-or-update) recomendado para sincronizaciones

Cuando sincronizas desde una API (como en `SyncStocksFromAPI`), a veces conviene un **upsert** en una sola operación.

GORM lo soporta con `Clauses(OnConflict ...)`:
```go
import "gorm.io/gorm/clause"

err := db.Clauses(clause.OnConflict{
  Columns: []clause.Column{{Name: "ticker"}, {Name: "time"}},
  DoUpdates: clause.AssignmentColumns([]string{
    "target_from", "target_to", "action", "rating_from", "rating_to", "brokerage", "company",
    "updated_at",
  }),
}).Create(&stock).Error
```

> Para que esto sea robusto, suele ser buena idea tener un **índice/constraint único** en `(ticker, time)`.

---

## 8) Joins y asociaciones

### `Joins`
```go
// Ejemplo genérico (depende de tus tablas)
db.Joins("JOIN users u ON u.id = stocks.user_id").Where("u.active = ?", true).Find(&stocks)
```

### `Preload` (eager loading)
Si tienes relaciones definidas (has many, belongs to, etc.):
```go
db.Preload("Orders").First(&user, userID)
```

---

## 9) Transacciones

Usa `Transaction` para operaciones que deben ser atómicas:
```go
err := db.Transaction(func(tx *gorm.DB) error {
  if err := tx.Create(&a).Error; err != nil { return err }
  if err := tx.Create(&b).Error; err != nil { return err }
  return nil
})
```

---

## 10) SQL crudo (Raw / Exec) y mapeo

### `Raw` (SELECT)
```go
type Row struct {
  Ticker string
  Cnt    int64
}
var rows []Row
err := db.Raw(`
  SELECT ticker, COUNT(*) as cnt
  FROM stocks
  GROUP BY ticker
  ORDER BY cnt DESC
  LIMIT ?`, 10).
  Scan(&rows).Error
```

### `Exec` (INSERT/UPDATE/DDL)
```go
err := db.Exec("UPDATE stocks SET brokerage = '' WHERE brokerage IS NULL").Error
```

---

## 11) Scopes (reusar filtros)

Los scopes son funciones que devuelven `func(*gorm.DB) *gorm.DB`:
```go
func ByTicker(t string) func(*gorm.DB) *gorm.DB {
  return func(db *gorm.DB) *gorm.DB {
    return db.Where("ticker = ?", t)
  }
}

func RecentFirst() func(*gorm.DB) *gorm.DB {
  return func(db *gorm.DB) *gorm.DB {
    return db.Order("time DESC")
  }
}

db.Scopes(ByTicker("AAPL"), RecentFirst()).Limit(20).Find(&stocks)
```

---

## 12) Debug y logging de queries

### `Debug()` (imprime SQL)
```go
db.Debug().Where("ticker = ?", "AAPL").First(&stock)
```

### `Session` útil (DryRun)
```go
stmt := db.Session(&gorm.Session{DryRun: true}).Where("ticker = ?", "AAPL").Find(&stocks).Statement
fmt.Println(stmt.SQL.String())
fmt.Println(stmt.Vars)
```

---

## 13) Buenas prácticas rápidas

- **Siempre** usa placeholders (`?`) para valores del usuario.
- Para búsquedas, crea índices en columnas consultadas frecuentemente (`ticker`, `time`, etc.).
- Para paginación estable, ordena por una columna determinista (`time DESC`, `id DESC`).
- Valida `sortBy`/`order` con una whitelist antes de construir `Order(...)`.
- Si pasas `context.Context` (por request): `db.WithContext(ctx)`.

---

## 14) Mini “cheat sheet” de métodos comunes

| Objetivo | Método(s) | Ejemplo corto |
|---|---|---|
| Un registro por PK | `First(&m, id)` | `db.First(&stock, 10)` |
| Un registro por filtro | `Where().First()` | `db.Where("ticker = ?", "AAPL").First(&s)` |
| Lista | `Where().Find()` | `db.Where("action = ?", a).Find(&xs)` |
| Ordenar | `Order()` | `db.Order("time DESC")...` |
| Paginación | `Limit()`, `Offset()` | `db.Limit(20).Offset(40)` |
| Contar | `Count()` | `db.Model(&Stock{}).Count(&n)` |
| Columnas específicas | `Select()` | `db.Select("ticker,time")...` |
| Una columna | `Pluck()` | `db.Pluck("ticker", &tickers)` |
| Distintos | `Distinct()` | `db.Distinct("action").Pluck(...)` |
| SQL crudo | `Raw()`, `Exec()` | `db.Raw("...").Scan(&out)` |
| Transacción | `Transaction()` | `db.Transaction(func(tx *gorm.DB) error { ... })` |

