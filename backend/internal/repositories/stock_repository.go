package repositories

import (
	"errors"
	"time"

	"github.com/Hitomiblood/StockStream/internal/models"
)

// ErrNotFound is returned when a repository lookup yields no results.
var ErrNotFound = errors.New("stock repository: not found")

// StockRepository abstracts persistence for stocks (services must not query the DB directly).
type StockRepository interface {
	GetByTickerAndTime(ticker string, t time.Time) (*models.Stock, error)
	Create(stock *models.Stock) error
	Save(stock *models.Stock) error

	CountAll() (int64, error)
	List(limit, offset int, sortField string, desc bool) ([]models.Stock, error)

	FindByID(id uint64) (*models.Stock, error)
	FindByTicker(ticker string) ([]models.Stock, error)
	Search(query string, limit int) ([]models.Stock, error)
	Filter(action, rating string, limit, offset int) ([]models.Stock, int64, error)

	DistinctActions() ([]string, error)
	DistinctRatings() ([]string, error)
	Latest(limit int) ([]models.Stock, error)

	FindSince(since time.Time) ([]models.Stock, error)
}
