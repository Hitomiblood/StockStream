package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Hitomiblood/StockStream/internal/config"
	"github.com/Hitomiblood/StockStream/internal/models"
	"github.com/go-resty/resty/v2"
)

type APIClient struct {
	client *resty.Client
	config *config.Config
}

// NewAPIClient crea una nueva instancia del cliente de API externa
func NewAPIClient(cfg *config.Config) *APIClient {
	client := resty.New()
	client.SetTimeout(30 * time.Second)
	client.SetRetryCount(3)
	client.SetRetryWaitTime(5 * time.Second)

	return &APIClient{
		client: client,
		config: cfg,
	}
}

// FetchStocks obtiene una página de stocks desde la API externa
func (ac *APIClient) FetchStocks(nextPage string) (*models.APIResponse, error) {
	url := ac.config.ExternalAPIURL
	
	req := ac.client.R().
		SetHeader("Authorization", "Bearer "+ac.config.ExternalAPIToken).
		SetHeader("Content-Type", "application/json")

	// Agregar parámetro de paginación si existe
	if nextPage != "" {
		req.SetQueryParam("next_page", nextPage)
	}

	resp, err := req.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stocks: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("API returned status code %d: %s", resp.StatusCode(), resp.String())
	}

	var apiResp models.APIResponse
	if err := json.Unmarshal(resp.Body(), &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	return &apiResp, nil
}

// FetchAllStocks obtiene todos los stocks disponibles de la API externa
// haciendo requests paginados hasta que no haya más páginas
func (ac *APIClient) FetchAllStocks() ([]models.Stock, error) {
	log.Println("🔄 Starting to fetch all stocks from external API...")
	
	var allStocks []models.Stock
	nextPage := ""
	pageCount := 0

	for {
		pageCount++
		log.Printf("📄 Fetching page %d...", pageCount)

		apiResp, err := ac.FetchStocks(nextPage)
		if err != nil {
			return nil, err
		}

		allStocks = append(allStocks, apiResp.Items...)
		log.Printf("✅ Page %d fetched: %d stocks", pageCount, len(apiResp.Items))

		// Si no hay más páginas, salir del loop
		if apiResp.NextPage == "" {
			break
		}

		nextPage = apiResp.NextPage

		// Pequeña pausa para no saturar la API
		time.Sleep(500 * time.Millisecond)
	}

	log.Printf("✅ Finished fetching all stocks. Total: %d stocks from %d pages", len(allStocks), pageCount)
	return allStocks, nil
}
