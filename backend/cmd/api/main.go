package main

import (
	"log"
	"net/http"

	"github.com/Hitomiblood/StockStream/internal/config"
	"github.com/Hitomiblood/StockStream/internal/database"
	"github.com/Hitomiblood/StockStream/internal/handlers"
	"github.com/Hitomiblood/StockStream/internal/middleware"
	"github.com/Hitomiblood/StockStream/internal/repositories/gormrepo"
	"github.com/Hitomiblood/StockStream/internal/services"
	"github.com/gin-gonic/gin"

	_ "github.com/Hitomiblood/StockStream/docs" // docs generados por swag init
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	stockRepo := gormrepo.NewStockRepository(database.GetDB())
	stockService := services.NewStockService(apiClient, stockRepo)
	recommendationService := services.NewRecommendationService(stockRepo)
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

	// Redirigir la raíz a Swagger
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/swagger/index.html")
	})

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
	addr := cfg.APIHost + ":" + cfg.APIPort
	log.Printf("\n🚀 Server starting on http://%s", addr)
	log.Printf("📚 Swagger UI: http://%s/swagger/index.html", addr)
	log.Printf("🧪 Test the API: curl http://%s/health\n", addr)

	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
