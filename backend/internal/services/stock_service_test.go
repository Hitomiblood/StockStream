package services

import (
	"errors"
	"testing"
	"time"

	"github.com/Hitomiblood/StockStream/internal/models"
	"github.com/Hitomiblood/StockStream/internal/repositories"
)

type fakeFetcher struct {
	stocks []models.Stock
	err    error
}

func (f *fakeFetcher) FetchAllStocks() ([]models.Stock, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.stocks, nil
}

type fakeRepo struct {
	getByTickerAndTimeFn func(ticker string, at time.Time) (*models.Stock, error)
	createFn             func(stock *models.Stock) error
	saveFn               func(stock *models.Stock) error
	countAllFn           func() (int64, error)
	listFn               func(limit, offset int, sortField string, desc bool) ([]models.Stock, error)
	findByIDFn           func(id uint64) (*models.Stock, error)
	findByTickerFn       func(ticker string) ([]models.Stock, error)
	searchFn             func(query string, limit int) ([]models.Stock, error)
	filterFn             func(action, rating string, limit, offset int) ([]models.Stock, int64, error)
	distinctActionsFn    func() ([]string, error)
	distinctRatingsFn    func() ([]string, error)
	latestFn             func(limit int) ([]models.Stock, error)
	findSinceFn          func(since time.Time) ([]models.Stock, error)
}

func (f *fakeRepo) GetByTickerAndTime(ticker string, at time.Time) (*models.Stock, error) {
	if f.getByTickerAndTimeFn != nil {
		return f.getByTickerAndTimeFn(ticker, at)
	}
	return nil, repositories.ErrNotFound
}
func (f *fakeRepo) Create(stock *models.Stock) error {
	if f.createFn != nil {
		return f.createFn(stock)
	}
	return nil
}
func (f *fakeRepo) Save(stock *models.Stock) error {
	if f.saveFn != nil {
		return f.saveFn(stock)
	}
	return nil
}
func (f *fakeRepo) CountAll() (int64, error) {
	if f.countAllFn != nil {
		return f.countAllFn()
	}
	return 0, nil
}
func (f *fakeRepo) List(limit, offset int, sortField string, desc bool) ([]models.Stock, error) {
	if f.listFn != nil {
		return f.listFn(limit, offset, sortField, desc)
	}
	return nil, nil
}
func (f *fakeRepo) FindByID(id uint64) (*models.Stock, error) {
	if f.findByIDFn != nil {
		return f.findByIDFn(id)
	}
	return nil, repositories.ErrNotFound
}
func (f *fakeRepo) FindByTicker(ticker string) ([]models.Stock, error) {
	if f.findByTickerFn != nil {
		return f.findByTickerFn(ticker)
	}
	return nil, nil
}
func (f *fakeRepo) Search(query string, limit int) ([]models.Stock, error) {
	if f.searchFn != nil {
		return f.searchFn(query, limit)
	}
	return nil, nil
}
func (f *fakeRepo) Filter(action, rating string, limit, offset int) ([]models.Stock, int64, error) {
	if f.filterFn != nil {
		return f.filterFn(action, rating, limit, offset)
	}
	return nil, 0, nil
}
func (f *fakeRepo) DistinctActions() ([]string, error) {
	if f.distinctActionsFn != nil {
		return f.distinctActionsFn()
	}
	return nil, nil
}
func (f *fakeRepo) DistinctRatings() ([]string, error) {
	if f.distinctRatingsFn != nil {
		return f.distinctRatingsFn()
	}
	return nil, nil
}
func (f *fakeRepo) Latest(limit int) ([]models.Stock, error) {
	if f.latestFn != nil {
		return f.latestFn(limit)
	}
	return nil, nil
}
func (f *fakeRepo) FindSince(since time.Time) ([]models.Stock, error) {
	if f.findSinceFn != nil {
		return f.findSinceFn(since)
	}
	return nil, nil
}

func TestNormalizeSortOrder(t *testing.T) {
	tests := []struct {
		name  string
		order string
		want  bool
	}{
		{name: "asc", order: "asc", want: false},
		{name: "ASC spaces", order: " ASC ", want: false},
		{name: "desc", order: "desc", want: true},
		{name: "empty default desc", order: "", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeSortOrder(tt.order); got != tt.want {
				t.Fatalf("normalizeSortOrder(%q) = %v, want %v", tt.order, got, tt.want)
			}
		})
	}
}

func TestSyncStocksFromAPI_CreateAndUpdate(t *testing.T) {
	now := time.Now()
	incoming := []models.Stock{
		{Ticker: "AAPL", Time: now, TargetFrom: "$100", TargetTo: "$110", Action: "Upgrade", RatingFrom: "Hold", RatingTo: "Buy", Brokerage: "X"},
		{Ticker: "MSFT", Time: now.Add(time.Minute), TargetFrom: "$200", TargetTo: "$205", Action: "Maintain", RatingFrom: "Buy", RatingTo: "Buy", Brokerage: "Y"},
	}

	creates := 0
	saves := 0
	repo := &fakeRepo{
		getByTickerAndTimeFn: func(ticker string, _ time.Time) (*models.Stock, error) {
			if ticker == "AAPL" {
				return nil, repositories.ErrNotFound
			}
			return &models.Stock{ID: 55, Ticker: "MSFT", Time: now.Add(time.Minute), TargetFrom: "$190", TargetTo: "$200", Action: "Maintain", RatingFrom: "Hold", RatingTo: "Buy", Brokerage: "Y"}, nil
		},
		createFn: func(stock *models.Stock) error {
			creates++
			return nil
		},
		saveFn: func(stock *models.Stock) error {
			saves++
			if stock.ID != 55 {
				t.Fatalf("expected updated stock keep ID=55, got %d", stock.ID)
			}
			return nil
		},
	}

	svc := NewStockService(&fakeFetcher{stocks: incoming}, repo)
	newCount, updatedCount, err := svc.SyncStocksFromAPI()
	if err != nil {
		t.Fatalf("SyncStocksFromAPI error: %v", err)
	}

	if newCount != 1 || updatedCount != 1 {
		t.Fatalf("counts got new=%d updated=%d, want 1/1", newCount, updatedCount)
	}
	if creates != 1 || saves != 1 {
		t.Fatalf("repo ops got create=%d save=%d, want 1/1", creates, saves)
	}
}

func TestSyncStocksFromAPI_FetchError(t *testing.T) {
	svc := NewStockService(&fakeFetcher{err: errors.New("boom")}, &fakeRepo{})
	_, _, err := svc.SyncStocksFromAPI()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestGetAllStocks_UsesSortNormalization(t *testing.T) {
	repo := &fakeRepo{
		countAllFn: func() (int64, error) { return 2, nil },
		listFn: func(limit, offset int, sortField string, desc bool) ([]models.Stock, error) {
			if limit != 10 || offset != 5 || sortField != "time" || !desc {
				t.Fatalf("unexpected args limit=%d offset=%d sort=%s desc=%v", limit, offset, sortField, desc)
			}
			return []models.Stock{{Ticker: "AAPL"}}, nil
		},
	}

	svc := NewStockService(&fakeFetcher{}, repo)
	stocks, total, err := svc.GetAllStocks(10, 5, "time", "unexpected")
	if err != nil {
		t.Fatalf("GetAllStocks error: %v", err)
	}
	if len(stocks) != 1 || total != 2 {
		t.Fatalf("unexpected result stocks=%d total=%d", len(stocks), total)
	}
}

func TestStockService_WrapperMethods(t *testing.T) {
	repo := &fakeRepo{
		findByIDFn:     func(id uint64) (*models.Stock, error) { return &models.Stock{ID: id, Ticker: "AAPL"}, nil },
		findByTickerFn: func(ticker string) ([]models.Stock, error) { return []models.Stock{{Ticker: ticker}}, nil },
		searchFn:       func(query string, limit int) ([]models.Stock, error) { return []models.Stock{{Ticker: query}}, nil },
		filterFn: func(action, rating string, limit, offset int) ([]models.Stock, int64, error) {
			return []models.Stock{{Action: action, RatingTo: rating}}, 1, nil
		},
		distinctActionsFn: func() ([]string, error) { return []string{"Upgrade"}, nil },
		distinctRatingsFn: func() ([]string, error) { return []string{"Buy"}, nil },
		latestFn:          func(limit int) ([]models.Stock, error) { return []models.Stock{{Ticker: "MSFT"}}, nil },
	}

	svc := NewStockService(&fakeFetcher{}, repo)

	if stock, err := svc.GetStockByID(7); err != nil || stock.ID != 7 {
		t.Fatalf("GetStockByID failed: stock=%+v err=%v", stock, err)
	}
	if stocks, err := svc.GetStocksByTicker("AAPL"); err != nil || len(stocks) != 1 {
		t.Fatalf("GetStocksByTicker failed: len=%d err=%v", len(stocks), err)
	}
	if stocks, err := svc.SearchStocks("AAP", 10); err != nil || len(stocks) != 1 {
		t.Fatalf("SearchStocks failed: len=%d err=%v", len(stocks), err)
	}
	if stocks, total, err := svc.FilterStocks("Upgrade", "Buy", 10, 0); err != nil || total != 1 || len(stocks) != 1 {
		t.Fatalf("FilterStocks failed: total=%d len=%d err=%v", total, len(stocks), err)
	}
	if actions, err := svc.GetUniqueActions(); err != nil || len(actions) != 1 {
		t.Fatalf("GetUniqueActions failed: len=%d err=%v", len(actions), err)
	}
	if ratings, err := svc.GetUniqueRatings(); err != nil || len(ratings) != 1 {
		t.Fatalf("GetUniqueRatings failed: len=%d err=%v", len(ratings), err)
	}
	if latest, err := svc.GetLatestStocks(5); err != nil || len(latest) != 1 {
		t.Fatalf("GetLatestStocks failed: len=%d err=%v", len(latest), err)
	}
}
