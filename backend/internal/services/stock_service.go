package services

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Hitomiblood/StockStream/internal/models"
	"github.com/Hitomiblood/StockStream/internal/repositories"
)

type StockService struct {
	repo      repositories.StockRepository
	apiClient StockFetcher
}

type StockFetcher interface {
	FetchAllStocks() ([]models.Stock, error)
}

// NewStockService crea una nueva instancia del servicio de stocks
func NewStockService(apiClient StockFetcher, repo repositories.StockRepository) *StockService {
	return &StockService{
		repo:      repo,
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
		existing, err := s.repo.GetByTickerAndTime(stock.Ticker, stock.Time)
		if err != nil && err != repositories.ErrNotFound {
			log.Printf("⚠️  Error checking existing stock %s: %v", stock.Ticker, err)
			continue
		}

		if err == repositories.ErrNotFound {
			// No existe, crear nuevo
			if err := s.repo.Create(&stock); err != nil {
				log.Printf("⚠️  Error creating stock %s: %v", stock.Ticker, err)
				continue
			}
			log.Printf("🆕 Created new stock: ID=%d, Ticker=%s", stock.ID, stock.Ticker)
			totalNew++
		} else if existing != nil {
			// Ya existe, actualizar si hay cambios
			if s.hasChanges(existing, &stock) {
				stock.ID = existing.ID // Mantener el ID
				if err := s.repo.Save(&stock); err != nil {
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
	// Contar total
	total, err := s.repo.CountAll()
	if err != nil {
		return nil, 0, err
	}

	desc := normalizeSortOrder(order)

	stocks, err := s.repo.List(limit, offset, sortBy, desc)
	if err != nil {
		return nil, 0, err
	}

	return stocks, total, nil
}

func normalizeSortOrder(order string) bool {
	normalized := strings.ToLower(strings.TrimSpace(order))
	return normalized != "asc"
}

// GetStockByID obtiene un stock por su ID
func (s *StockService) GetStockByID(id uint64) (*models.Stock, error) {
	return s.repo.FindByID(id)
}

// GetStocksByTicker obtiene el historial de un stock por su ticker
func (s *StockService) GetStocksByTicker(ticker string) ([]models.Stock, error) {
	return s.repo.FindByTicker(ticker)
}

// SearchStocks busca stocks por ticker o nombre de compañía
func (s *StockService) SearchStocks(query string, limit int) ([]models.Stock, error) {
	return s.repo.Search(query, limit)
}

// FilterStocks filtra stocks por acción y/o rating
func (s *StockService) FilterStocks(action, rating string, limit, offset int) ([]models.Stock, int64, error) {
	return s.repo.Filter(action, rating, limit, offset)
}

// GetUniqueActions obtiene todas las acciones únicas disponibles
func (s *StockService) GetUniqueActions() ([]string, error) {
	return s.repo.DistinctActions()
}

// GetUniqueRatings obtiene todos los ratings únicos disponibles
func (s *StockService) GetUniqueRatings() ([]string, error) {
	return s.repo.DistinctRatings()
}

// GetLatestStocks obtiene los últimos N stocks añadidos
func (s *StockService) GetLatestStocks(limit int) ([]models.Stock, error) {
	return s.repo.Latest(limit)
}
