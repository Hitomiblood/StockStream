package services

import (
	"fmt"
	"log"
	"time"

	"github.com/Hitomiblood/StockStream/internal/database"
	"github.com/Hitomiblood/StockStream/internal/models"
	"gorm.io/gorm"
)

type StockService struct {
	db        *gorm.DB
	apiClient *APIClient
}

// NewStockService crea una nueva instancia del servicio de stocks
func NewStockService(apiClient *APIClient) *StockService {
	return &StockService{
		db:        database.GetDB(),
		apiClient: apiClient,
	}
}

// SyncStocksFromAPI sincroniza los datos desde la API externa a la base de datos
func (s *StockService) SyncStocksFromAPI() (int, int, error) {
	startTime := time.Now()
	log.Println("🔄 Starting stock synchronization...")

	// Obtener todos los stocks de la API
	stocks, err := s.apiClient.FetchAllStocks()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to fetch stocks: %w", err)
	}

	totalNew := 0
	totalUpdated := 0

	// Procesar cada stock
	for _, stock := range stocks {
		// Verificar si ya existe (por ticker y time)
		var existing models.Stock
		result := s.db.Where("ticker = ? AND time = ?", stock.Ticker, stock.Time).First(&existing)

		if result.Error == gorm.ErrRecordNotFound {
			// No existe, crear nuevo
			if err := s.db.Create(&stock).Error; err != nil {
				log.Printf("⚠️  Error creating stock %s: %v", stock.Ticker, err)
				continue
			}
			totalNew++
		} else if result.Error == nil {
			// Ya existe, actualizar si hay cambios
			if s.hasChanges(&existing, &stock) {
				stock.ID = existing.ID // Mantener el ID
				if err := s.db.Save(&stock).Error; err != nil {
					log.Printf("⚠️  Error updating stock %s: %v", stock.Ticker, err)
					continue
				}
				totalUpdated++
			}
		}
	}

	duration := time.Since(startTime)
	log.Printf("✅ Sync completed: %d new, %d updated in %v", totalNew, totalUpdated, duration)

	return totalNew, totalUpdated, nil
}

// hasChanges verifica si hay diferencias entre dos stocks
func (s *StockService) hasChanges(old, new *models.Stock) bool {
	return old.TargetFrom != new.TargetFrom ||
		old.TargetTo != new.TargetTo ||
		old.Action != new.Action ||
		old.RatingFrom != new.RatingFrom ||
		old.RatingTo != new.RatingTo ||
		old.Brokerage != new.Brokerage
}

// GetAllStocks obtiene todos los stocks con paginación
func (s *StockService) GetAllStocks(limit, offset int, sortBy, order string) ([]models.Stock, int64, error) {
	var stocks []models.Stock
	var total int64

	// Contar total
	s.db.Model(&models.Stock{}).Count(&total)

	// Aplicar paginación y ordenamiento
	query := s.db.Limit(limit).Offset(offset)

	if sortBy != "" {
		orderClause := sortBy
		if order == "desc" {
			orderClause += " DESC"
		} else {
			orderClause += " ASC"
		}
		query = query.Order(orderClause)
	} else {
		query = query.Order("time DESC") // Por defecto ordenar por fecha descendente
	}

	if err := query.Find(&stocks).Error; err != nil {
		return nil, 0, err
	}

	return stocks, total, nil
}

// GetStockByID obtiene un stock por su ID
func (s *StockService) GetStockByID(id uint) (*models.Stock, error) {
	var stock models.Stock
	if err := s.db.First(&stock, id).Error; err != nil {
		return nil, err
	}
	return &stock, nil
}

// GetStocksByTicker obtiene el historial de un stock por su ticker
func (s *StockService) GetStocksByTicker(ticker string) ([]models.Stock, error) {
	var stocks []models.Stock
	if err := s.db.Where("ticker = ?", ticker).Order("time DESC").Find(&stocks).Error; err != nil {
		return nil, err
	}
	return stocks, nil
}

// SearchStocks busca stocks por ticker o nombre de compañía
func (s *StockService) SearchStocks(query string, limit int) ([]models.Stock, error) {
	var stocks []models.Stock
	searchPattern := "%" + query + "%"

	if err := s.db.Where("ticker ILIKE ? OR company ILIKE ?", searchPattern, searchPattern).
		Limit(limit).
		Order("time DESC").
		Find(&stocks).Error; err != nil {
		return nil, err
	}

	return stocks, nil
}

// FilterStocks filtra stocks por acción y/o rating
func (s *StockService) FilterStocks(action, rating string, limit, offset int) ([]models.Stock, int64, error) {
	var stocks []models.Stock
	var total int64

	query := s.db.Model(&models.Stock{})

	if action != "" {
		query = query.Where("action = ?", action)
	}
	if rating != "" {
		query = query.Where("rating_to = ?", rating)
	}

	query.Count(&total)

	if err := query.Limit(limit).Offset(offset).Order("time DESC").Find(&stocks).Error; err != nil {
		return nil, 0, err
	}

	return stocks, total, nil
}

// GetUniqueActions obtiene todas las acciones únicas disponibles
func (s *StockService) GetUniqueActions() ([]string, error) {
	var actions []string
	if err := s.db.Model(&models.Stock{}).Distinct("action").Pluck("action", &actions).Error; err != nil {
		return nil, err
	}
	return actions, nil
}

// GetUniqueRatings obtiene todos los ratings únicos disponibles
func (s *StockService) GetUniqueRatings() ([]string, error) {
	var ratings []string
	if err := s.db.Model(&models.Stock{}).Distinct("rating_to").Pluck("rating_to", &ratings).Error; err != nil {
		return nil, err
	}
	return ratings, nil
}

// GetLatestStocks obtiene los últimos N stocks añadidos
func (s *StockService) GetLatestStocks(limit int) ([]models.Stock, error) {
	var stocks []models.Stock
	if err := s.db.Order("created_at DESC").Limit(limit).Find(&stocks).Error; err != nil {
		return nil, err
	}
	return stocks, nil
}
