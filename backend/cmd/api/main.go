package main

import (
	"log"

	"github.com/Hitomiblood/StockStream/internal/config"
	"github.com/Hitomiblood/StockStream/internal/database"
	"github.com/Hitomiblood/StockStream/internal/handlers"
	"github.com/Hitomiblood/StockStream/internal/middleware"
	"github.com/Hitomiblood/StockStream/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	// Cargar configuración
	cfg := config.Load()
	log.Println("✅ Configuration loaded")

	// Conectar a la base de datos
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	log.Println("✅ Database connected")

	// Crear servicios
	apiClient := services.NewAPIClient(cfg)
	stockService := services.NewStockService(apiClient)
	recommendationService := services.NewRecommendationService()
	log.Println("✅ Services initialized")

	// Crear handlers
	stockHandler := handlers.NewStockHandler(stockService, recommendationService)

	// Configurar Gin según el nivel de log
	if cfg.LogLevel != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Crear router
	r := gin.Default()

	// Aplicar middleware
	r.Use(middleware.CORS())

	// Rutas públicas
	r.GET("/health", stockHandler.HealthCheck)

	// Grupo de rutas API v1
	v1 := r.Group("/api/v1")
	{
		// Stocks
		v1.GET("/stocks", stockHandler.GetAllStocks)
		v1.GET("/stocks/latest", stockHandler.GetLatestStocks)
		v1.GET("/stocks/search", stockHandler.SearchStocks)
		v1.GET("/stocks/filter", stockHandler.FilterStocks)
		v1.GET("/stocks/ticker/:ticker", stockHandler.GetStocksByTicker)
		v1.GET("/stocks/:id", stockHandler.GetStockByID)
		v1.POST("/stocks/fetch", stockHandler.FetchStocks)

		// Recomendaciones
		v1.GET("/recommendations", stockHandler.GetRecommendations)

		// Metadata
		v1.GET("/metadata", stockHandler.GetMetadata)
	}

	// Información de rutas disponibles
	log.Println("\n📋 Available endpoints:")
	log.Println("   GET  /health - Health check")
	log.Println("   GET  /api/v1/stocks - List all stocks")
	log.Println("   GET  /api/v1/stocks/latest - Get latest stocks")
	log.Println("   GET  /api/v1/stocks/search?q=query - Search stocks")
	log.Println("   GET  /api/v1/stocks/filter?action=X&rating=Y - Filter stocks")
	log.Println("   GET  /api/v1/stocks/ticker/:ticker - Get stock history by ticker")
	log.Println("   GET  /api/v1/stocks/:id - Get stock by ID")
	log.Println("   POST /api/v1/stocks/fetch - Fetch from external API")
	log.Println("   GET  /api/v1/recommendations - Get investment recommendations")
	log.Println("   GET  /api/v1/metadata - Get available actions and ratings")

	// Iniciar servidor
	addr := cfg.APIHost + ":" + cfg.APIPort
	log.Printf("\n🚀 Server starting on http://%s", addr)
	log.Printf("📚 Test the API: curl http://%s/health\n", addr)
	
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
