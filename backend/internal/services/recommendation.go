package services

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Hitomiblood/StockStream/internal/database"
	"github.com/Hitomiblood/StockStream/internal/models"
	"gorm.io/gorm"
)

type RecommendationService struct {
	db *gorm.DB
}

// NewRecommendationService crea una nueva instancia del servicio de recomendaciones
func NewRecommendationService() *RecommendationService {
	return &RecommendationService{
		db: database.GetDB(),
	}
}

// GetRecommendations obtiene las mejores recomendaciones de inversión
func (rs *RecommendationService) GetRecommendations(limit int) ([]models.StockRecommendation, error) {
	log.Println("🔍 Calculating stock recommendations...")

	// Obtener stocks recientes (últimos 30 días)
	var stocks []models.Stock
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	if err := rs.db.Where("time > ?", thirtyDaysAgo).Find(&stocks).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch stocks: %w", err)
	}

	if len(stocks) == 0 {
		return []models.StockRecommendation{}, nil
	}

	// Agrupar stocks por ticker
	stocksByTicker := make(map[string][]models.Stock)
	for _, stock := range stocks {
		stocksByTicker[stock.Ticker] = append(stocksByTicker[stock.Ticker], stock)
	}

	// Calcular score para cada ticker
	var recommendations []models.StockRecommendation

	for _, tickerStocks := range stocksByTicker {
		if len(tickerStocks) == 0 {
			continue
		}

		// Usar el stock más reciente como referencia
		latestStock := tickerStocks[0]
		for _, s := range tickerStocks {
			if s.Time.After(latestStock.Time) {
				latestStock = s
			}
		}

		score, reason, confidence := rs.calculateScore(latestStock, tickerStocks)

		if score > 0 { // Solo incluir recomendaciones positivas
			recommendations = append(recommendations, models.StockRecommendation{
				Stock:      latestStock,
				Score:      score,
				Reason:     reason,
				Confidence: confidence,
			})
		}
	}

	// Ordenar por score descendente
	recommendations = rs.sortRecommendations(recommendations)

	// Limitar resultados
	if limit > 0 && len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}

	log.Printf("✅ Generated %d recommendations", len(recommendations))
	return recommendations, nil
}

// calculateScore calcula el score de una acción basado en múltiples criterios
func (rs *RecommendationService) calculateScore(stock models.Stock, history []models.Stock) (float64, string, string) {
	score := 0.0
	reasons := []string{}

	// Criterio 1: Incremento/Decremento de Target Price (40% del peso)
	targetScore := rs.evaluateTargetChange(stock)
	score += targetScore * 0.4
	if targetScore > 0 {
		reasons = append(reasons, fmt.Sprintf("Target price increase (+%.0f points)", targetScore*0.4))
	} else if targetScore < 0 {
		reasons = append(reasons, fmt.Sprintf("Target price decrease (%.0f points)", targetScore*0.4))
	}

	// Criterio 2: Mejora/Empeoramiento de Rating (30% del peso)
	ratingScore := rs.evaluateRatingChange(stock)
	score += ratingScore * 0.3
	if ratingScore > 0 {
		reasons = append(reasons, fmt.Sprintf("Rating upgraded to %s (+%.0f points)", stock.RatingTo, ratingScore*0.3))
	} else if ratingScore < 0 {
		reasons = append(reasons, fmt.Sprintf("Rating downgraded (%.0f points)", ratingScore*0.3))
	}

	// Criterio 3: Actividad Reciente (30% del peso)
	activityScore := rs.evaluateRecentActivity(stock, history)
	score += activityScore * 0.3
	if activityScore > 0 {
		reasons = append(reasons, fmt.Sprintf("High recent activity (+%.0f points)", activityScore*0.3))
	}

	// Determinar nivel de confianza
	confidence := rs.determineConfidence(score, len(history))

	// Construir razón final
	finalReason := strings.Join(reasons, ". ")
	if finalReason == "" {
		finalReason = "No significant changes detected"
	}

	return score, finalReason, confidence
}

// evaluateTargetChange evalúa el cambio en el precio objetivo
func (rs *RecommendationService) evaluateTargetChange(stock models.Stock) float64 {
	fromPrice := rs.parsePrice(stock.TargetFrom)
	toPrice := rs.parsePrice(stock.TargetTo)

	if fromPrice == 0 || toPrice == 0 {
		return 0
	}

	percentChange := ((toPrice - fromPrice) / fromPrice) * 100

	// Convertir porcentaje a score (max 100 puntos)
	if percentChange > 0 {
		return min(percentChange*10, 100) // Incremento positivo
	} else {
		return max(percentChange*10, -100) // Decremento negativo
	}
}

// evaluateRatingChange evalúa el cambio en el rating
func (rs *RecommendationService) evaluateRatingChange(stock models.Stock) float64 {
	fromRating := rs.ratingToScore(stock.RatingFrom)
	toRating := rs.ratingToScore(stock.RatingTo)

	diff := toRating - fromRating

	// Convertir diferencia a score
	return diff * 20 // Cada nivel vale 20 puntos
}

// evaluateRecentActivity evalúa la actividad reciente del stock
func (rs *RecommendationService) evaluateRecentActivity(stock models.Stock, history []models.Stock) float64 {
	daysSinceUpdate := time.Since(stock.Time).Hours() / 24

	// Más reciente = mejor score
	if daysSinceUpdate < 7 {
		return 100 // Muy reciente
	} else if daysSinceUpdate < 14 {
		return 70 // Reciente
	} else if daysSinceUpdate < 21 {
		return 40 // Moderadamente reciente
	} else {
		return 10 // Antiguo
	}
}

// ratingToScore convierte un rating en un valor numérico
func (rs *RecommendationService) ratingToScore(rating string) float64 {
	rating = strings.ToLower(strings.TrimSpace(rating))

	switch {
	case strings.Contains(rating, "strong buy"):
		return 5
	case strings.Contains(rating, "buy"):
		return 4
	case strings.Contains(rating, "outperform") || strings.Contains(rating, "overweight"):
		return 3.5
	case strings.Contains(rating, "hold") || strings.Contains(rating, "neutral"):
		return 2
	case strings.Contains(rating, "underperform") || strings.Contains(rating, "underweight"):
		return 1
	case strings.Contains(rating, "sell"):
		return 0
	default:
		return 2 // Neutral por defecto
	}
}

// parsePrice extrae el valor numérico de un precio (ej: "$150.00" -> 150.00)
func (rs *RecommendationService) parsePrice(price string) float64 {
	// Remover símbolos de moneda y espacios
	clean := strings.TrimSpace(price)
	clean = strings.ReplaceAll(clean, "$", "")
	clean = strings.ReplaceAll(clean, "€", "")
	clean = strings.ReplaceAll(clean, ",", "")

	value, err := strconv.ParseFloat(clean, 64)
	if err != nil {
		return 0
	}
	return value
}

// determineConfidence determina el nivel de confianza de la recomendación
func (rs *RecommendationService) determineConfidence(score float64, historyCount int) string {
	// Más historial = mayor confianza
	if score >= 60 && historyCount >= 3 {
		return "high"
	} else if score >= 40 && historyCount >= 2 {
		return "medium"
	} else {
		return "low"
	}
}

// sortRecommendations ordena las recomendaciones por score descendente
func (rs *RecommendationService) sortRecommendations(recs []models.StockRecommendation) []models.StockRecommendation {
	// Bubble sort simple (suficiente para conjuntos pequeños)
	for i := 0; i < len(recs)-1; i++ {
		for j := 0; j < len(recs)-i-1; j++ {
			if recs[j].Score < recs[j+1].Score {
				recs[j], recs[j+1] = recs[j+1], recs[j]
			}
		}
	}
	return recs
}

// Helper functions
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
