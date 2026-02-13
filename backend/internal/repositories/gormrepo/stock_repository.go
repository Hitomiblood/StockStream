package gormrepo

import (
	"errors"
	"fmt"
	"time"

	"github.com/Hitomiblood/StockStream/internal/models"
	"github.com/Hitomiblood/StockStream/internal/repositories"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StockRepository struct {
	db *gorm.DB
}

func NewStockRepository(db *gorm.DB) *StockRepository {
	return &StockRepository{db: db}
}

func (r *StockRepository) GetByTickerAndTime(ticker string, t time.Time) (*models.Stock, error) {
	var stock models.Stock
	err := r.db.Where("ticker = ? AND time = ?", ticker, t).First(&stock).Error
	if err == nil {
		return &stock, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repositories.ErrNotFound
	}
	return nil, err
}

func (r *StockRepository) Create(stock *models.Stock) error {
	return r.db.Create(stock).Error
}

func (r *StockRepository) Save(stock *models.Stock) error {
	return r.db.Save(stock).Error
}

func (r *StockRepository) CountAll() (int64, error) {
	var total int64
	if err := r.db.Model(&models.Stock{}).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (r *StockRepository) List(limit, offset int, sortField string, desc bool) ([]models.Stock, error) {
	var stocks []models.Stock

	q := r.db.
		Limit(limit).
		Offset(offset).
		Order(clause.OrderByColumn{Column: clause.Column{Name: sortField}, Desc: desc})

	if err := q.Find(&stocks).Error; err != nil {
		return nil, err
	}
	return stocks, nil
}

func (r *StockRepository) FindByID(id uint64) (*models.Stock, error) {
	var stock models.Stock
	if err := r.db.First(&stock, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repositories.ErrNotFound
		}
		return nil, err
	}
	return &stock, nil
}

func (r *StockRepository) FindByTicker(ticker string) ([]models.Stock, error) {
	var stocks []models.Stock
	if err := r.db.Where("ticker = ?", ticker).Order("time DESC").Find(&stocks).Error; err != nil {
		return nil, err
	}
	return stocks, nil
}

func (r *StockRepository) Search(query string, limit int) ([]models.Stock, error) {
	var stocks []models.Stock
	searchPattern := "%" + query + "%"

	if err := r.db.Where("ticker ILIKE ? OR company ILIKE ?", searchPattern, searchPattern).
		Limit(limit).
		Order("time DESC").
		Find(&stocks).Error; err != nil {
		return nil, err
	}

	return stocks, nil
}

func (r *StockRepository) Filter(action, rating string, limit, offset int) ([]models.Stock, int64, error) {
	var stocks []models.Stock
	var total int64

	q := r.db.Model(&models.Stock{})
	if action != "" {
		q = q.Where("action = ?", action)
	}
	if rating != "" {
		q = q.Where("rating_to = ?", rating)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := q.Limit(limit).Offset(offset).Order("time DESC").Find(&stocks).Error; err != nil {
		return nil, 0, err
	}

	return stocks, total, nil
}

func (r *StockRepository) DistinctActions() ([]string, error) {
	var actions []string
	if err := r.db.Model(&models.Stock{}).Distinct("action").Pluck("action", &actions).Error; err != nil {
		return nil, err
	}
	return actions, nil
}

func (r *StockRepository) DistinctRatings() ([]string, error) {
	var ratings []string
	if err := r.db.Model(&models.Stock{}).Distinct("rating_to").Pluck("rating_to", &ratings).Error; err != nil {
		return nil, err
	}
	return ratings, nil
}

func (r *StockRepository) Latest(limit int) ([]models.Stock, error) {
	var stocks []models.Stock
	if err := r.db.Order("created_at DESC").Limit(limit).Find(&stocks).Error; err != nil {
		return nil, err
	}
	return stocks, nil
}

func (r *StockRepository) FindSince(since time.Time) ([]models.Stock, error) {
	var stocks []models.Stock
	if err := r.db.Where("time > ?", since).Find(&stocks).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch stocks: %w", err)
	}
	return stocks, nil
}
