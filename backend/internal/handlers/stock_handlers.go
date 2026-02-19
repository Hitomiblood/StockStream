package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Hitomiblood/StockStream/internal/models"
	"github.com/Hitomiblood/StockStream/internal/services"
	"github.com/gin-gonic/gin"
)

type stockService interface {
	GetAllStocks(limit, offset int, sortBy, order string) ([]models.Stock, int64, error)
	GetStockByID(id uint64) (*models.Stock, error)
	GetStocksByTicker(ticker string) ([]models.Stock, error)
	SearchStocks(query string, limit int) ([]models.Stock, error)
	FilterStocks(action, rating string, limit, offset int) ([]models.Stock, int64, error)
	SyncStocksFromAPI() (int, int, error)
	GetUniqueActions() ([]string, error)
	GetUniqueRatings() ([]string, error)
	GetLatestStocks(limit int) ([]models.Stock, error)
}

type recommendationService interface {
	GetRecommendations(limit int) ([]models.StockRecommendation, error)
}

type StockHandler struct {
	stockService          stockService
	recommendationService recommendationService
	now                   func() time.Time
}

var allowedSortQueryFields = map[string]struct{}{
	"id":          {},
	"ticker":      {},
	"target_from": {},
	"target_to":   {},
	"company":     {},
	"action":      {},
	"brokerage":   {},
	"rating_from": {},
	"rating_to":   {},
	"time":        {},
	"created_at":  {},
	"updated_at":  {},
}

// NewStockHandler crea una nueva instancia del handler
func NewStockHandler(stockService *services.StockService, recService *services.RecommendationService) *StockHandler {
	return &StockHandler{
		stockService:          stockService,
		recommendationService: recService,
		now:                   time.Now,
	}
}

func NewStockHandlerWithServices(stockService stockService, recService recommendationService) *StockHandler {
	return &StockHandler{
		stockService:          stockService,
		recommendationService: recService,
		now:                   time.Now,
	}
}

// GetAllStocks maneja GET /api/v1/stocks
// @Summary      List stocks
// @Description  Get all stocks with pagination, sorting and filtering
// @Tags         stocks
// @Accept       json
// @Produce      json
// @Param        limit   query  int     false  "Number of results (default: 50, max: 200)"
// @Param        offset  query  int     false  "Offset for pagination (default: 0)"
// @Param        sort    query  string  false  "Field to sort by (default: time)"
// @Param        order   query  string  false  "Sort order: asc or desc (default: desc)"
// @Success      200  {object}  map[string]interface{}  "List of stocks"
// @Failure      500  {object}  map[string]interface{}  "Internal server error"
// @Router       /api/v1/stocks [get]
func (h *StockHandler) GetAllStocks(c *gin.Context) {
	// Parsear parámetros de query
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	sortBy := c.DefaultQuery("sort", "time")
	order := c.DefaultQuery("order", "desc")
	sortBy, order = sanitizeSortParams(sortBy, order)

	// Validar límite
	if limit > 200 {
		limit = 200
	}
	if limit < 1 {
		limit = 50
	}

	stocks, total, err := h.stockService.GetAllStocks(limit, offset, sortBy, order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch stocks",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   stocks,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func sanitizeSortParams(sortBy, order string) (string, string) {
	normalizedSort := strings.ToLower(strings.TrimSpace(sortBy))
	if _, ok := allowedSortQueryFields[normalizedSort]; !ok {
		normalizedSort = "time"
	}

	normalizedOrder := strings.ToLower(strings.TrimSpace(order))
	if normalizedOrder != "asc" && normalizedOrder != "desc" {
		normalizedOrder = "desc"
	}
	return normalizedSort, normalizedOrder
}

// GetStockByID maneja GET /api/v1/stocks/:id
// @Summary      Get stock by ID
// @Description  Get detailed information of a specific stock
// @Tags         stocks
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Stock ID (as string due to large ID values)"
// @Success      200  {object}  models.Stock
// @Failure      400  {object}  map[string]interface{}  "Invalid ID"
// @Failure      404  {object}  map[string]interface{}  "Stock not found"
// @Router       /api/v1/stocks/{id} [get]
func (h *StockHandler) GetStockByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid stock ID",
		})
		return
	}

	stock, err := h.stockService.GetStockByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Stock not found",
		})
		return
	}

	c.JSON(http.StatusOK, stock)
}

// GetStocksByTicker maneja GET /api/v1/stocks/ticker/:ticker
// @Summary      Get stocks by ticker
// @Description  Get all historical records for a specific stock ticker
// @Tags         stocks
// @Accept       json
// @Produce      json
// @Param        ticker  path      string  true  "Stock ticker symbol"
// @Success      200  {object}  map[string]interface{}  "Stock history for ticker"
// @Failure      400  {object}  map[string]interface{}  "Invalid ticker"
// @Failure      404  {object}  map[string]interface{}  "No stocks found for ticker"
// @Failure      500  {object}  map[string]interface{}  "Internal server error"
// @Router       /api/v1/stocks/ticker/{ticker} [get]
func (h *StockHandler) GetStocksByTicker(c *gin.Context) {
	ticker := c.Param("ticker")
	if ticker == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Ticker is required",
		})
		return
	}

	stocks, err := h.stockService.GetStocksByTicker(ticker)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch stocks",
		})
		return
	}

	if len(stocks) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No stocks found for ticker " + ticker,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ticker":  ticker,
		"company": stocks[0].Company,
		"history": stocks,
		"total":   len(stocks),
	})
}

// SearchStocks maneja GET /api/v1/stocks/search
// @Summary      Search stocks
// @Description  Search stocks by ticker or company name
// @Tags         stocks
// @Accept       json
// @Produce      json
// @Param        q      query  string  true   "Search query (ticker or company name)"
// @Param        limit  query  int     false  "Number of results (default: 50, max: 200)"
// @Success      200  {object}  map[string]interface{}  "Search results"
// @Failure      400  {object}  map[string]interface{}  "Missing search query"
// @Failure      500  {object}  map[string]interface{}  "Internal server error"
// @Router       /api/v1/stocks/search [get]
func (h *StockHandler) SearchStocks(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Search query 'q' is required",
		})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if limit > 200 {
		limit = 200
	}

	stocks, err := h.stockService.SearchStocks(query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to search stocks",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"query": query,
		"data":  stocks,
		"total": len(stocks),
	})
}

// FilterStocks maneja GET /api/v1/stocks/filter
// @Summary      Filter stocks
// @Description  Filter stocks by action, rating or both criteria
// @Tags         stocks
// @Accept       json
// @Produce      json
// @Param        action  query  string  false  "Filter by action (e.g., Upgrade, Downgrade, Initiated, Maintains)"
// @Param        rating  query  string  false  "Filter by rating (e.g., Buy, Sell, Hold)"
// @Param        limit   query  int     false  "Number of results (default: 50, max: 200)"
// @Param        offset  query  int     false  "Offset for pagination (default: 0)"
// @Success      200  {object}  map[string]interface{}  "Filtered results"
// @Failure      400  {object}  map[string]interface{}  "Missing filter parameters"
// @Failure      500  {object}  map[string]interface{}  "Internal server error"
// @Router       /api/v1/stocks/filter [get]
func (h *StockHandler) FilterStocks(c *gin.Context) {
	action := c.Query("action")
	rating := c.Query("rating")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if action == "" && rating == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "At least one filter parameter (action or rating) is required",
		})
		return
	}

	stocks, total, err := h.stockService.FilterStocks(action, rating, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to filter stocks",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"filters": gin.H{
			"action": action,
			"rating": rating,
		},
		"data":   stocks,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// FetchStocks maneja POST /api/v1/stocks/fetch
// @Summary      Sync stocks from external API
// @Description  Fetch and synchronize all stocks from the external API to the database
// @Tags         stocks
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "Sync completed successfully"
// @Failure      500  {object}  map[string]interface{}  "Failed to sync stocks"
// @Router       /api/v1/stocks/fetch [post]
func (h *StockHandler) FetchStocks(c *gin.Context) {
	startTime := h.now()

	totalNew, totalUpdated, err := h.stockService.SyncStocksFromAPI()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to sync stocks: " + err.Error(),
		})
		return
	}

	duration := h.now().Sub(startTime)

	c.JSON(http.StatusOK, gin.H{
		"message":       "Data fetched successfully",
		"total_new":     totalNew,
		"total_updated": totalUpdated,
		"total_fetched": totalNew + totalUpdated,
		"duration_ms":   duration.Milliseconds(),
	})
}

// GetRecommendations maneja GET /api/v1/recommendations
// @Summary      Get investment recommendations
// @Description  Get intelligent stock recommendations based on multiple criteria
// @Tags         recommendations
// @Accept       json
// @Produce      json
// @Param        limit  query  int  false  "Number of recommendations (default: 10, max: 50)"
// @Success      200  {object}  map[string]interface{}  "Recommendations payload"
// @Failure      500  {object}  map[string]interface{}  "Failed to generate recommendations"
// @Router       /api/v1/recommendations [get]
func (h *StockHandler) GetRecommendations(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit > 50 {
		limit = 50
	}
	if limit < 1 {
		limit = 10
	}

	recommendations, err := h.recommendationService.GetRecommendations(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate recommendations",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"recommendations": recommendations,
		"generated_at":    h.now(),
		"count":           len(recommendations),
		"criteria": gin.H{
			"target_change_weight":   0.4,
			"rating_change_weight":   0.3,
			"recent_activity_weight": 0.3,
		},
	})
}

// GetMetadata maneja GET /api/v1/metadata
// @Summary      Get metadata
// @Description  Get unique values for actions and ratings available in the database
// @Tags         metadata
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "Available actions and ratings"
// @Failure      500  {object}  map[string]interface{}  "Failed to fetch metadata"
// @Router       /api/v1/metadata [get]
func (h *StockHandler) GetMetadata(c *gin.Context) {
	actions, err1 := h.stockService.GetUniqueActions()
	ratings, err2 := h.stockService.GetUniqueRatings()

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch metadata",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"actions": actions,
		"ratings": ratings,
	})
}

// GetLatestStocks maneja GET /api/v1/stocks/latest
// @Summary      Get latest stocks
// @Description  Get the most recent stock updates ordered by time
// @Tags         stocks
// @Accept       json
// @Produce      json
// @Param        limit  query  int  false  "Number of results (default: 20, max: 100)"
// @Success      200  {object}  map[string]interface{}  "Latest stocks"
// @Failure      500  {object}  map[string]interface{}  "Internal server error"
// @Router       /api/v1/stocks/latest [get]
func (h *StockHandler) GetLatestStocks(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit > 100 {
		limit = 100
	}

	stocks, err := h.stockService.GetLatestStocks(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch latest stocks",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  stocks,
		"total": len(stocks),
	})
}

// HealthCheck maneja GET /health
// @Summary      Health check
// @Description  Check if the API is running
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "API is healthy"
// @Router       /health [get]
func (h *StockHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"message":   "Stock Analysis System API",
		"timestamp": h.now(),
		"version":   "1.0.0",
	})
}
