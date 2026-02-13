package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Hitomiblood/StockStream/internal/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/cockroachdb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	if len(os.Args) < 2 {
		usageAndExit()
	}

	cmd := os.Args[1]
	cfg := config.Load()

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

	dsn := dbURL.String()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("❌ Failed to open DB: %v", err)
	}
	defer db.Close()

	driver, err := cockroachdb.WithInstance(db, &cockroachdb.Config{
		DatabaseName: cfg.DBName,
	})
	if err != nil {
		log.Fatalf("❌ Failed to init migrate driver: %v", err)
	}

	sourceURL := fileSourceURL(findMigrationsDir())
	m, err := migrate.NewWithDatabaseInstance(sourceURL, "cockroachdb", driver)
	if err != nil {
		log.Fatalf("❌ Failed to init migrate: %v", err)
	}
	defer func() {
		_, _ = m.Close()
	}()

	switch cmd {
	case "up":
		err = m.Up()
	case "down":
		err = m.Down()
	case "steps":
		if len(os.Args) < 3 {
			usageAndExit()
		}
		n, convErr := strconv.Atoi(os.Args[2])
		if convErr != nil {
			log.Fatalf("❌ Invalid steps value: %v", convErr)
		}
		err = m.Steps(n)
	case "version":
		v, dirty, verr := m.Version()
		if errors.Is(verr, migrate.ErrNilVersion) {
			fmt.Println("version: none")
			return
		}
		if verr != nil {
			log.Fatalf("❌ Failed to get version: %v", verr)
		}
		fmt.Printf("version: %d (dirty=%v)\n", v, dirty)
		return
	default:
		usageAndExit()
	}

	if err == nil {
		log.Println("✅ Migrations completed")
		return
	}
	if errors.Is(err, migrate.ErrNoChange) {
		log.Println("✅ No migrations to apply")
		return
	}

	log.Fatalf("❌ Migration failed: %v", err)
}

func usageAndExit() {
	fmt.Println("Usage: go run ./cmd/migrate <up|down|version|steps N>")
	os.Exit(2)
}

func findMigrationsDir() string {
	// Prefer current working directory layout: ./migrations (when run from backend/)
	if stat, err := os.Stat("migrations"); err == nil && stat.IsDir() {
		abs, _ := filepath.Abs("migrations")
		return abs
	}

	// Fallback when run from repo root: ./backend/migrations
	if stat, err := os.Stat(filepath.Join("backend", "migrations")); err == nil && stat.IsDir() {
		abs, _ := filepath.Abs(filepath.Join("backend", "migrations"))
		return abs
	}

	// Default: return absolute path to expected dir (migrate will error clearly if missing).
	abs, _ := filepath.Abs("migrations")
	return abs
}

func fileSourceURL(absDir string) string {
	u := &url.URL{
		Scheme: "file",
		Path:   filepath.ToSlash(absDir),
	}
	return u.String()
}
