package database

import (
	"fmt"
	"log"
	"net/url"

	"github.com/Hitomiblood/StockStream/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func buildDSN(cfg *config.Config) string {
	userInfo := url.User(cfg.DBUser)
	if cfg.DBPassword != "" {
		userInfo = url.UserPassword(cfg.DBUser, cfg.DBPassword)
	}

	dbURL := &url.URL{
		Scheme: "postgres",
		User:   userInfo,
		Host:   fmt.Sprintf("%s:%d", cfg.DBHost, cfg.DBPort),
		Path:   cfg.DBName,
	}

	query := dbURL.Query()
	query.Set("sslmode", cfg.DBSSLMode)
	query.Set("search_path", cfg.DBSchema)
	dbURL.RawQuery = query.Encode()

	return dbURL.String()
}

func resolveLogLevel(level string) logger.LogLevel {
	if level == "debug" {
		return logger.Info
	}
	return logger.Warn
}

// Connect establece la conexión a CockroachDB y ejecuta las migraciones
func Connect(cfg *config.Config) error {
	dsn := buildDSN(cfg)

	// Configurar nivel de log según configuración
	logLevel := resolveLogLevel(cfg.LogLevel)

	// Conectar a la base de datos
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("✅ Database connection established")

	return nil
}

// GetDB retorna la instancia de la base de datos
func GetDB() *gorm.DB {
	return DB
}

// Close cierra la conexión a la base de datos
func Close() error {
	if DB == nil {
		return nil
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
