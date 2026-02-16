package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Hitomiblood/StockStream/internal/models"
	"github.com/gin-gonic/gin"
)

type fakeStockService struct {
	getAllStocksFn    func(limit, offset int, sortBy, order string) ([]models.Stock, int64, error)
	getStockByIDFn    func(id uint64) (*models.Stock, error)
	searchStocksFn    func(query string, limit int) ([]models.Stock, error)
	filterStocksFn    func(action, rating string, limit, offset int) ([]models.Stock, int64, error)
	syncStocksFn      func() (int, int, error)
	getActionsFn      func() ([]string, error)
	getRatingsFn      func() ([]string, error)
	getLatestStocksFn func(limit int) ([]models.Stock, error)
	getByTickerFn     func(ticker string) ([]models.Stock, error)
}

func (f *fakeStockService) GetAllStocks(limit, offset int, sortBy, order string) ([]models.Stock, int64, error) {
	if f.getAllStocksFn != nil {
		return f.getAllStocksFn(limit, offset, sortBy, order)
	}
	return nil, 0, nil
}
func (f *fakeStockService) GetStockByID(id uint64) (*models.Stock, error) {
	if f.getStockByIDFn != nil {
		return f.getStockByIDFn(id)
	}
	return nil, errors.New("not found")
}
func (f *fakeStockService) GetStocksByTicker(ticker string) ([]models.Stock, error) {
	if f.getByTickerFn != nil {
		return f.getByTickerFn(ticker)
	}
	return nil, nil
}
func (f *fakeStockService) SearchStocks(query string, limit int) ([]models.Stock, error) {
	if f.searchStocksFn != nil {
		return f.searchStocksFn(query, limit)
	}
	return nil, nil
}
func (f *fakeStockService) FilterStocks(action, rating string, limit, offset int) ([]models.Stock, int64, error) {
	if f.filterStocksFn != nil {
		return f.filterStocksFn(action, rating, limit, offset)
	}
	return nil, 0, nil
}
func (f *fakeStockService) SyncStocksFromAPI() (int, int, error) {
	if f.syncStocksFn != nil {
		return f.syncStocksFn()
	}
	return 0, 0, nil
}
func (f *fakeStockService) GetUniqueActions() ([]string, error) {
	if f.getActionsFn != nil {
		return f.getActionsFn()
	}
	return []string{}, nil
}
func (f *fakeStockService) GetUniqueRatings() ([]string, error) {
	if f.getRatingsFn != nil {
		return f.getRatingsFn()
	}
	return []string{}, nil
}
func (f *fakeStockService) GetLatestStocks(limit int) ([]models.Stock, error) {
	if f.getLatestStocksFn != nil {
		return f.getLatestStocksFn(limit)
	}
	return nil, nil
}

type fakeRecommendationService struct {
	getRecommendationsFn func(limit int) ([]models.StockRecommendation, error)
}

func (f *fakeRecommendationService) GetRecommendations(limit int) ([]models.StockRecommendation, error) {
	if f.getRecommendationsFn != nil {
		return f.getRecommendationsFn(limit)
	}
	return nil, nil
}

func TestSanitizeSortParams(t *testing.T) {
	sortBy, order := sanitizeSortParams("  INVALID_FIELD  ", "  ASC ")
	if sortBy != "time" || order != "asc" {
		t.Fatalf("got sort=%q order=%q", sortBy, order)
	}
}

func TestGetAllStocks_SanitizesParamsAndReturns200(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	h := NewStockHandlerWithServices(&fakeStockService{
		getAllStocksFn: func(limit, offset int, sortBy, order string) ([]models.Stock, int64, error) {
			if limit != 50 || offset != 0 || sortBy != "time" || order != "desc" {
				t.Fatalf("unexpected params: limit=%d offset=%d sort=%s order=%s", limit, offset, sortBy, order)
			}
			return []models.Stock{{Ticker: "AAPL"}}, 1, nil
		},
	}, &fakeRecommendationService{})

	r.GET("/stocks", h.GetAllStocks)
	req := httptest.NewRequest(http.MethodGet, "/stocks?limit=-1&sort=bad&order=bad", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestGetStockByID_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewStockHandlerWithServices(&fakeStockService{}, &fakeRecommendationService{})
	r.GET("/stocks/:id", h.GetStockByID)

	req := httptest.NewRequest(http.MethodGet, "/stocks/not-a-number", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestFetchStocks_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewStockHandlerWithServices(&fakeStockService{
		syncStocksFn: func() (int, int, error) {
			return 3, 2, nil
		},
	}, &fakeRecommendationService{})

	start := time.Date(2026, 2, 13, 12, 0, 0, 0, time.UTC)
	end := start.Add(150 * time.Millisecond)
	calls := 0
	h.now = func() time.Time {
		calls++
		if calls == 1 {
			return start
		}
		return end
	}

	r.POST("/fetch", h.FetchStocks)
	req := httptest.NewRequest(http.MethodPost, "/fetch", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestGetRecommendations_ClampsLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	h := NewStockHandlerWithServices(&fakeStockService{}, &fakeRecommendationService{
		getRecommendationsFn: func(limit int) ([]models.StockRecommendation, error) {
			if limit != 50 {
				t.Fatalf("expected clamped limit 50, got %d", limit)
			}
			return []models.StockRecommendation{}, nil
		},
	})
	h.now = func() time.Time { return time.Date(2026, 2, 13, 12, 0, 0, 0, time.UTC) }

	r.GET("/recs", h.GetRecommendations)
	req := httptest.NewRequest(http.MethodGet, "/recs?limit=1000", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestGetStockByID_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewStockHandlerWithServices(&fakeStockService{
		getStockByIDFn: func(id uint64) (*models.Stock, error) {
			return nil, errors.New("not found")
		},
	}, &fakeRecommendationService{})
	r.GET("/stocks/:id", h.GetStockByID)

	req := httptest.NewRequest(http.MethodGet, "/stocks/10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestGetStocksByTicker_ValidationsAndNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewStockHandlerWithServices(&fakeStockService{
		getByTickerFn: func(ticker string) ([]models.Stock, error) {
			return []models.Stock{}, nil
		},
	}, &fakeRecommendationService{})
	r.GET("/ticker/:ticker", h.GetStocksByTicker)

	req := httptest.NewRequest(http.MethodGet, "/ticker/AAPL", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestGetStocksByTicker_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewStockHandlerWithServices(&fakeStockService{
		getByTickerFn: func(ticker string) ([]models.Stock, error) {
			return []models.Stock{{Ticker: ticker, Company: "Apple"}}, nil
		},
	}, &fakeRecommendationService{})
	r.GET("/ticker/:ticker", h.GetStocksByTicker)

	req := httptest.NewRequest(http.MethodGet, "/ticker/AAPL", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestSearchStocks_MissingQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewStockHandlerWithServices(&fakeStockService{}, &fakeRecommendationService{})
	r.GET("/search", h.SearchStocks)

	req := httptest.NewRequest(http.MethodGet, "/search", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestFilterStocks_MissingFilters(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewStockHandlerWithServices(&fakeStockService{}, &fakeRecommendationService{})
	r.GET("/filter", h.FilterStocks)

	req := httptest.NewRequest(http.MethodGet, "/filter", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestFilterStocks_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewStockHandlerWithServices(&fakeStockService{
		filterStocksFn: func(action, rating string, limit, offset int) ([]models.Stock, int64, error) {
			return []models.Stock{{Ticker: "AAPL"}}, 1, nil
		},
	}, &fakeRecommendationService{})
	r.GET("/filter", h.FilterStocks)

	req := httptest.NewRequest(http.MethodGet, "/filter?action=Upgrade", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestGetMetadata_ErrorAndSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	errHandler := NewStockHandlerWithServices(&fakeStockService{
		getActionsFn: func() ([]string, error) { return nil, errors.New("x") },
	}, &fakeRecommendationService{})
	r.GET("/meta-err", errHandler.GetMetadata)

	reqErr := httptest.NewRequest(http.MethodGet, "/meta-err", nil)
	wErr := httptest.NewRecorder()
	r.ServeHTTP(wErr, reqErr)
	if wErr.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", wErr.Code, http.StatusInternalServerError)
	}

	okHandler := NewStockHandlerWithServices(&fakeStockService{
		getActionsFn: func() ([]string, error) { return []string{"Upgrade"}, nil },
		getRatingsFn: func() ([]string, error) { return []string{"Buy"}, nil },
	}, &fakeRecommendationService{})
	r.GET("/meta-ok", okHandler.GetMetadata)

	reqOK := httptest.NewRequest(http.MethodGet, "/meta-ok", nil)
	wOK := httptest.NewRecorder()
	r.ServeHTTP(wOK, reqOK)
	if wOK.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", wOK.Code, http.StatusOK)
	}
}

func TestGetLatestStocks_AndHealth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewStockHandlerWithServices(&fakeStockService{
		getLatestStocksFn: func(limit int) ([]models.Stock, error) {
			if limit != 100 {
				t.Fatalf("expected clamped limit=100, got %d", limit)
			}
			return []models.Stock{{Ticker: "AAPL"}}, nil
		},
	}, &fakeRecommendationService{})
	h.now = func() time.Time { return time.Date(2026, 2, 13, 12, 0, 0, 0, time.UTC) }

	r.GET("/latest", h.GetLatestStocks)
	r.GET("/health", h.HealthCheck)

	wLatest := httptest.NewRecorder()
	r.ServeHTTP(wLatest, httptest.NewRequest(http.MethodGet, "/latest?limit=200", nil))
	if wLatest.Code != http.StatusOK {
		t.Fatalf("status latest = %d, want %d", wLatest.Code, http.StatusOK)
	}

	wHealth := httptest.NewRecorder()
	r.ServeHTTP(wHealth, httptest.NewRequest(http.MethodGet, "/health", nil))
	if wHealth.Code != http.StatusOK {
		t.Fatalf("status health = %d, want %d", wHealth.Code, http.StatusOK)
	}
}
