package gormrepo

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Hitomiblood/StockStream/internal/repositories"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newMockedRepo(t *testing.T) (*StockRepository, sqlmock.Sqlmock, func()) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}

	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("gorm.Open error: %v", err)
	}

	cleanup := func() {
		_ = sqlDB.Close()
	}

	return NewStockRepository(gdb), mock, cleanup
}

func TestGetByTickerAndTime_NotFoundMapsDomainError(t *testing.T) {
	repo, mock, cleanup := newMockedRepo(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "stocks" WHERE ticker = $1 AND time = $2 ORDER BY "stocks"."id" LIMIT $3`)).
		WithArgs("AAPL", sqlmock.AnyArg(), 1).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := repo.GetByTickerAndTime("AAPL", time.Now())
	if !errors.Is(err, repositories.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestCountAll_Success(t *testing.T) {
	repo, mock, cleanup := newMockedRepo(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(7)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "stocks"`)).WillReturnRows(rows)

	total, err := repo.CountAll()
	if err != nil {
		t.Fatalf("CountAll error: %v", err)
	}
	if total != 7 {
		t.Fatalf("total = %d, want 7", total)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestList_Success(t *testing.T) {
	repo, mock, cleanup := newMockedRepo(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "ticker"}).AddRow(1, "AAPL")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "stocks" ORDER BY "time" DESC LIMIT $1`)).
		WithArgs(10).
		WillReturnRows(rows)

	stocks, err := repo.List(10, 0, "time", true)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if len(stocks) != 1 || stocks[0].Ticker != "AAPL" {
		t.Fatalf("unexpected stocks: %+v", stocks)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestFindSince_WrapsError(t *testing.T) {
	repo, mock, cleanup := newMockedRepo(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "stocks" WHERE time > $1`)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnError(errors.New("db error"))

	_, err := repo.FindSince(time.Now().AddDate(0, 0, -30))
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestFindByID_NotFound(t *testing.T) {
	repo, mock, cleanup := newMockedRepo(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "stocks" WHERE "stocks"."id" = $1 ORDER BY "stocks"."id" LIMIT $2`)).
		WithArgs(9, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := repo.FindByID(9)
	if !errors.Is(err, repositories.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestFindByTicker_Success(t *testing.T) {
	repo, mock, cleanup := newMockedRepo(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "ticker"}).AddRow(1, "AAPL")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "stocks" WHERE ticker = $1 ORDER BY time DESC`)).
		WithArgs("AAPL").
		WillReturnRows(rows)

	stocks, err := repo.FindByTicker("AAPL")
	if err != nil {
		t.Fatalf("FindByTicker error: %v", err)
	}
	if len(stocks) != 1 {
		t.Fatalf("expected 1 stock, got %d", len(stocks))
	}
}

func TestSearch_Success(t *testing.T) {
	repo, mock, cleanup := newMockedRepo(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "ticker"}).AddRow(1, "AAPL")
	mock.ExpectQuery(`SELECT \* FROM "stocks" WHERE ticker ILIKE \$1 OR company ILIKE \$2 ORDER BY time DESC LIMIT \$3`).
		WithArgs("%app%", "%app%", 5).
		WillReturnRows(rows)

	stocks, err := repo.Search("app", 5)
	if err != nil {
		t.Fatalf("Search error: %v", err)
	}
	if len(stocks) != 1 {
		t.Fatalf("expected 1 stock, got %d", len(stocks))
	}
}

func TestFilter_Success(t *testing.T) {
	repo, mock, cleanup := newMockedRepo(t)
	defer cleanup()

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
	mock.ExpectQuery(`SELECT count\(\*\) FROM "stocks" WHERE action = \$1`).
		WithArgs("Upgrade").
		WillReturnRows(countRows)

	dataRows := sqlmock.NewRows([]string{"id", "ticker", "action"}).AddRow(1, "AAPL", "Upgrade")
	mock.ExpectQuery(`SELECT \* FROM "stocks" WHERE action = \$1 ORDER BY time DESC LIMIT \$2`).
		WithArgs("Upgrade", 10).
		WillReturnRows(dataRows)

	stocks, total, err := repo.Filter("Upgrade", "", 10, 0)
	if err != nil {
		t.Fatalf("Filter error: %v", err)
	}
	if total != 1 || len(stocks) != 1 {
		t.Fatalf("unexpected filter result total=%d len=%d", total, len(stocks))
	}
}

func TestDistinctAndLatest(t *testing.T) {
	repo, mock, cleanup := newMockedRepo(t)
	defer cleanup()

	actionsRows := sqlmock.NewRows([]string{"action"}).AddRow("Upgrade")
	mock.ExpectQuery(`SELECT DISTINCT "action" FROM "stocks"`).WillReturnRows(actionsRows)

	ratingsRows := sqlmock.NewRows([]string{"rating_to"}).AddRow("Buy")
	mock.ExpectQuery(`SELECT DISTINCT "rating_to" FROM "stocks"`).WillReturnRows(ratingsRows)

	latestRows := sqlmock.NewRows([]string{"id", "ticker"}).AddRow(1, "AAPL")
	mock.ExpectQuery(`SELECT \* FROM "stocks" ORDER BY created_at DESC LIMIT \$1`).WithArgs(5).WillReturnRows(latestRows)

	actions, err := repo.DistinctActions()
	if err != nil || len(actions) != 1 {
		t.Fatalf("DistinctActions error=%v len=%d", err, len(actions))
	}

	ratings, err := repo.DistinctRatings()
	if err != nil || len(ratings) != 1 {
		t.Fatalf("DistinctRatings error=%v len=%d", err, len(ratings))
	}

	latest, err := repo.Latest(5)
	if err != nil || len(latest) != 1 {
		t.Fatalf("Latest error=%v len=%d", err, len(latest))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}
