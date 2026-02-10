package main

import (
	"log"

	"github.com/Hitomiblood/StockStream/internal/config"
	"github.com/gin-gonic/gin"
)

func main() {
	// Cargar configuración
	cfg := config.Load()

	// Crear router
	r := gin.Default()

	// Ruta de prueba
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Stock Analysis System API",
		})
	})

	// Iniciar servidor
	addr := cfg.APIHost + ":" + cfg.APIPort
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
