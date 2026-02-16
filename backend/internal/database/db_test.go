package database

import (
	"strings"
	"testing"

	"github.com/Hitomiblood/StockStream/internal/config"
	"gorm.io/gorm/logger"
)

func TestBuildDSN(t *testing.T) {
	cfg := &config.Config{
		DBHost:     "localhost",
		DBPort:     26257,
		DBUser:     "root",
		DBPassword: "secret",
		DBName:     "stockdb",
		DBSchema:   "public",
		DBSSLMode:  "disable",
	}

	dsn := buildDSN(cfg)
	if !strings.Contains(dsn, "postgres://root:secret@localhost:26257/stockdb") {
		t.Fatalf("unexpected dsn: %s", dsn)
	}
	if !strings.Contains(dsn, "sslmode=disable") || !strings.Contains(dsn, "search_path=public") {
		t.Fatalf("missing query params in dsn: %s", dsn)
	}
}

func TestResolveLogLevel(t *testing.T) {
	if got := resolveLogLevel("debug"); got != logger.Info {
		t.Fatalf("debug level mismatch: %v", got)
	}
	if got := resolveLogLevel("info"); got != logger.Warn {
		t.Fatalf("default level mismatch: %v", got)
	}
}

func TestCloseNilDB(t *testing.T) {
	DB = nil
	if err := Close(); err != nil {
		t.Fatalf("Close with nil DB should not fail: %v", err)
	}
}
