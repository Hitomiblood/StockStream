package services

import (
	"errors"
	"testing"
	"time"

	"github.com/Hitomiblood/StockStream/internal/models"
)

type recoRepo struct {
	findSinceFn func(since time.Time) ([]models.Stock, error)
}

func (r *recoRepo) GetByTickerAndTime(string, time.Time) (*models.Stock, error) { return nil, nil }
func (r *recoRepo) Create(*models.Stock) error                                  { return nil }
func (r *recoRepo) Save(*models.Stock) error                                    { return nil }
func (r *recoRepo) CountAll() (int64, error)                                    { return 0, nil }
func (r *recoRepo) List(int, int, string, bool) ([]models.Stock, error)         { return nil, nil }
func (r *recoRepo) FindByID(uint64) (*models.Stock, error)                      { return nil, nil }
func (r *recoRepo) FindByTicker(string) ([]models.Stock, error)                 { return nil, nil }
func (r *recoRepo) Search(string, int) ([]models.Stock, error)                  { return nil, nil }
func (r *recoRepo) Filter(string, string, int, int) ([]models.Stock, int64, error) {
	return nil, 0, nil
}
func (r *recoRepo) DistinctActions() ([]string, error) { return nil, nil }
func (r *recoRepo) DistinctRatings() ([]string, error) { return nil, nil }
func (r *recoRepo) Latest(int) ([]models.Stock, error) { return nil, nil }
func (r *recoRepo) FindSince(since time.Time) ([]models.Stock, error) {
	if r.findSinceFn != nil {
		return r.findSinceFn(since)
	}
	return nil, nil
}

func TestNormalizeText(t *testing.T) {
	if got := normalizeText("  Strong-Buy!!! "); got != "strong buy" {
		t.Fatalf("normalizeText got %q", got)
	}
}

func TestClamp(t *testing.T) {
	if got := clamp(120, 0, 100); got != 100 {
		t.Fatalf("clamp upper got %v", got)
	}
	if got := clamp(-10, 0, 100); got != 0 {
		t.Fatalf("clamp lower got %v", got)
	}
	if got := clamp(55, 0, 100); got != 55 {
		t.Fatalf("clamp mid got %v", got)
	}
}

func TestRecommendationHelpers(t *testing.T) {
	rs := NewRecommendationService(&recoRepo{})

	if got := rs.parsePrice("$1,234.50 USD"); got != 1234.5 {
		t.Fatalf("parsePrice got %v", got)
	}

	if got := rs.ratingToScore("Strong Buy"); got < 4.9 {
		t.Fatalf("ratingToScore strong buy too low: %v", got)
	}

	if got := rs.evaluateActionSignal("Target raised by broker"); got <= 0 {
		t.Fatalf("evaluateActionSignal expected positive, got %v", got)
	}
}

func TestGetRecommendations_FallbackWindowAndLimit(t *testing.T) {
	fixedNow := time.Date(2026, 2, 13, 12, 0, 0, 0, time.UTC)
	call := 0
	repo := &recoRepo{
		findSinceFn: func(_ time.Time) ([]models.Stock, error) {
			call++
			if call == 1 {
				return []models.Stock{}, nil
			}
			return []models.Stock{
				{Ticker: "AAA", Company: "A", TargetFrom: "$100", TargetTo: "$130", Action: "Target raised", RatingFrom: "Hold", RatingTo: "Buy", Time: fixedNow.AddDate(0, 0, -1)},
				{Ticker: "BBB", Company: "B", TargetFrom: "$100", TargetTo: "$80", Action: "Downgraded", RatingFrom: "Buy", RatingTo: "Sell", Time: fixedNow.AddDate(0, 0, -2)},
			}, nil
		},
	}

	rs := NewRecommendationService(repo)
	rs.now = func() time.Time { return fixedNow }

	recs, err := rs.GetRecommendations(1)
	if err != nil {
		t.Fatalf("GetRecommendations error: %v", err)
	}

	if call != 2 {
		t.Fatalf("FindSince calls = %d, want 2", call)
	}

	if len(recs) != 1 {
		t.Fatalf("len(recs) = %d, want 1", len(recs))
	}
}

func TestGetRecommendations_RepoError(t *testing.T) {
	rs := NewRecommendationService(&recoRepo{findSinceFn: func(_ time.Time) ([]models.Stock, error) {
		return nil, errors.New("db down")
	}})
	_, err := rs.GetRecommendations(10)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
