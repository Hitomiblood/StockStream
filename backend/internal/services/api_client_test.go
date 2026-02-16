package services

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Hitomiblood/StockStream/internal/config"
	"github.com/go-resty/resty/v2"
)

func newTestAPIClient(url string) *APIClient {
	return &APIClient{
		client: resty.New(),
		config: &config.Config{
			ExternalAPIURL:   url,
			ExternalAPIToken: "token",
		},
	}
}

func TestFetchStocks_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"items":[{"ticker":"AAPL"}],"next_page":""}`))
	}))
	defer ts.Close()

	client := newTestAPIClient(ts.URL)
	resp, err := client.FetchStocks("")
	if err != nil {
		t.Fatalf("FetchStocks error: %v", err)
	}

	if len(resp.Items) != 1 || resp.Items[0].Ticker != "AAPL" {
		t.Fatalf("unexpected payload: %+v", resp.Items)
	}
}

func TestFetchStocks_Non200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(`bad gateway`))
	}))
	defer ts.Close()

	client := newTestAPIClient(ts.URL)
	_, err := client.FetchStocks("")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestFetchAllStocks_Pagination(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextPage := r.URL.Query().Get("next_page")
		w.Header().Set("Content-Type", "application/json")
		if nextPage == "" {
			_, _ = w.Write([]byte(`{"items":[{"ticker":"AAA"}],"next_page":"2"}`))
			return
		}
		if nextPage == "2" {
			_, _ = w.Write([]byte(`{"items":[{"ticker":"BBB"}],"next_page":""}`))
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"error":"unexpected next_page %s"}`, nextPage)))
	}))
	defer ts.Close()

	client := newTestAPIClient(ts.URL)
	stocks, err := client.FetchAllStocks()
	if err != nil {
		t.Fatalf("FetchAllStocks error: %v", err)
	}

	if len(stocks) != 2 {
		t.Fatalf("len(stocks) = %d, want 2", len(stocks))
	}
}
