package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Hitomiblood/StockStream/internal/services"
	"github.com/gin-gonic/gin"
)

type StockHandler struct {
	stockService          *services.StockService
	recommendationService *services.RecommendationService
}

// NewStockHandler crea una nueva instancia del handler
func NewStockHandler(stockService *services.StockService, recService *services.RecommendationService) *StockHandler {
	return &StockHandler{
		stockService:          stockService,
		recommendationService: recService,
	}
}

// GetAllStocks maneja GET /api/v1/stocks
func (h *StockHandler) GetAllStocks(c *gin.Context) {
	// Parsear parámetros de query
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	sortBy := c.DefaultQuery("sort", "time")
	order := c.DefaultQuery("order", "desc")

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

// GetStockByID maneja GET /api/v1/stocks/:id
func (h *StockHandler) GetStockByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid stock ID",
		})
		return
	}

	stock, err := h.stockService.GetStockByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Stock not found",
		})
		return
	}

	c.JSON(http.StatusOK, stock)
}

// GetStocksByTicker maneja GET /api/v1/stocks/ticker/:ticker
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
func (h *StockHandler) FetchStocks(c *gin.Context) {
	startTime := time.Now()

	totalNew, totalUpdated, err := h.stockService.SyncStocksFromAPI()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to sync stocks: " + err.Error(),
		})
		return
	}

	duration := time.Since(startTime)

	c.JSON(http.StatusOK, gin.H{
		"message":       "Data fetched successfully",
		"total_new":     totalNew,
		"total_updated": totalUpdated,
		"total_fetched": totalNew + totalUpdated,
		"duration_ms":   duration.Milliseconds(),
	})
}

// GetRecommendations maneja GET /api/v1/recommendations
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
		"generated_at":    time.Now(),
		"count":           len(recommendations),
		"criteria": gin.H{
			"target_change_weight":   0.4,
			"rating_change_weight":   0.3,
			"recent_activity_weight": 0.3,
		},
	})
}

// GetMetadata maneja GET /api/v1/metadata
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
func (h *StockHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"message":   "Stock Analysis System API",
		"timestamp": time.Now(),
		"version":   "1.0.0",
	})
}
