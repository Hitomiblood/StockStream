package database

import (
	"fmt"
	"log"

	"github.com/Hitomiblood/StockStream/internal/config"
	"github.com/Hitomiblood/StockStream/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Connect establece la conexión a CockroachDB y ejecuta las migraciones
func Connect(cfg *config.Config) error {
	// Construir DSN (Data Source Name)
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	// Configurar nivel de log según configuración
	var logLevel logger.LogLevel
	if cfg.LogLevel == "debug" {
		logLevel = logger.Info
	} else {
		logLevel = logger.Warn
	}

	// Conectar a la base de datos
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("✅ Database connection established")

	// Ejecutar migraciones automáticas
	if err := AutoMigrate(); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	return nil
}

// AutoMigrate ejecuta las migraciones automáticas de GORM
func AutoMigrate() error {
	log.Println("🔄 Running database migrations...")

	if err := DB.AutoMigrate(&models.Stock{}); err != nil {
		return err
	}

	log.Println("✅ Database migrations completed")
	return nil
}

// GetDB retorna la instancia de la base de datos
func GetDB() *gorm.DB {
	return DB
}

// Close cierra la conexión a la base de datos
func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
