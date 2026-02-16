package services

import (
	"fmt"
	"log"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Hitomiblood/StockStream/internal/models"
	"github.com/Hitomiblood/StockStream/internal/repositories"
)

type RecommendationService struct {
	repo repositories.StockRepository
	now  func() time.Time
}

type featureVector struct {
	PriceDirection    float64
	PriceMomentum     float64
	ActionRatingCombo float64
	RatingQuality     float64
	RatingChange      float64
	RecentActivity    float64
	HistoryConsensus  float64
	DataQuality       float64
}

type priceFluctuation struct {
	Direction     float64
	Momentum      float64
	PercentChange float64
	HasPricePair  bool
}

type weightedFeature struct {
	label        string
	weight       float64
	value        float64
	positiveText string
	neutralText  string
	negativeText string
	contribution float64
}

var noiseTextRegex = regexp.MustCompile(`[^a-z0-9\s]+`)

const (
	minimumAcceptedScore = 55.0
	maxHistoryItems      = 8
	defaultUnknownRating = 2.0
)

var ratingSignals = []struct {
	pattern string
	score   float64
}{
	{pattern: "strong buy", score: 5.0},
	{pattern: "speculative buy", score: 4.5},
	{pattern: "market outperform", score: 4.0},
	{pattern: "sector outperform", score: 4.0},
	{pattern: "outperformer", score: 4.0},
	{pattern: "buy", score: 4.0},
	{pattern: "overweight", score: 3.5},
	{pattern: "outperform", score: 3.5},
	{pattern: "positive", score: 3.5},
	{pattern: "accumulate", score: 3.0},
	{pattern: "equal weight", score: 2.0},
	{pattern: "in line", score: 2.0},
	{pattern: "market perform", score: 2.0},
	{pattern: "sector perform", score: 2.0},
	{pattern: "neutral", score: 2.0},
	{pattern: "hold", score: 2.0},
	{pattern: "cautious", score: 1.5},
	{pattern: "underweight", score: 1.0},
	{pattern: "underperform", score: 0.5},
	{pattern: "sector underperform", score: 0.5},
	{pattern: "reduce", score: 0.5},
	{pattern: "sell", score: 0.0},
	{pattern: "strong sell", score: 0.0},
}

var actionSignals = []struct {
	pattern string
	score   float64
}{
	{pattern: "target raised by", score: 92},
	{pattern: "target raised", score: 90},
	{pattern: "raises target", score: 90},
	{pattern: "raise target", score: 90},
	{pattern: "upgraded", score: 75},
	{pattern: "upgrade", score: 75},
	{pattern: "initiated with buy", score: 65},
	{pattern: "initiated", score: 30},
	{pattern: "reiterated buy", score: 40},
	{pattern: "reiterated", score: 15},
	{pattern: "maintains buy", score: 35},
	{pattern: "maintained buy", score: 35},
	{pattern: "maintains", score: 10},
	{pattern: "target lowered by", score: -92},
	{pattern: "target lowered", score: -90},
	{pattern: "lowers target", score: -90},
	{pattern: "lower target", score: -90},
	{pattern: "downgraded", score: -75},
	{pattern: "downgrade", score: -75},
	{pattern: "suspended", score: -60},
	{pattern: "removed", score: -50},
}

// NewRecommendationService crea una nueva instancia del servicio de recomendaciones
func NewRecommendationService(repo repositories.StockRepository) *RecommendationService {
	return &RecommendationService{
		repo: repo,
		now:  time.Now,
	}
}

// GetRecommendations obtiene las mejores recomendaciones de inversión
func (rs *RecommendationService) GetRecommendations(limit int) ([]models.StockRecommendation, error) {
	log.Println("🔍 Calculating stock recommendations...")

	var stocks []models.Stock

	// Buscar primero en la ventana reciente y luego ampliar si no hay datos.
	thirtyDaysAgo := rs.now().AddDate(0, 0, -30)
	var err error
	stocks, err = rs.repo.FindSince(thirtyDaysAgo)
	if err != nil {
		return nil, err
	}

	if len(stocks) == 0 {
		ninetyDaysAgo := rs.now().AddDate(0, 0, -90)
		stocks, err = rs.repo.FindSince(ninetyDaysAgo)
		if err != nil {
			return nil, err
		}
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

		sort.Slice(tickerStocks, func(i, j int) bool {
			return tickerStocks[i].Time.After(tickerStocks[j].Time)
		})
		latestStock := tickerStocks[0]

		score, reason, confidence := rs.calculateScore(latestStock, tickerStocks)
		recommendations = append(recommendations, models.StockRecommendation{
			Stock:      latestStock,
			Score:      score,
			Reason:     reason,
			Confidence: confidence,
		})
	}

	// Ordenar por score descendente
	recommendations = rs.sortRecommendations(recommendations)

	// Filtrar recomendaciones con score suficientemente bueno.
	filtered := make([]models.StockRecommendation, 0, len(recommendations))
	for _, rec := range recommendations {
		if rec.Score >= minimumAcceptedScore {
			filtered = append(filtered, rec)
		}
	}

	if len(filtered) > 0 {
		recommendations = filtered
	}

	// Limitar resultados
	if limit > 0 && len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}

	log.Printf("✅ Generated %d recommendations", len(recommendations))
	return recommendations, nil
}

// calculateScore calcula el score de una acción basado en múltiples criterios
func (rs *RecommendationService) calculateScore(stock models.Stock, history []models.Stock) (float64, string, models.ConfidenceLevel) {
	features := rs.buildFeatureVector(stock, history)
	weighted, rawScore := rs.scoreWithWeights(features, stock)
	finalScore := rs.calibrateScore(rawScore, features.DataQuality)
	confidence := rs.determineConfidence(finalScore, len(history), features.DataQuality)
	finalReason := rs.buildReason(weighted)
	if finalReason == "" {
		finalReason = "No significant changes detected"
	}

	return finalScore, finalReason, confidence
}

func (rs *RecommendationService) buildFeatureVector(stock models.Stock, history []models.Stock) featureVector {
	fluctuation, targetQuality := rs.evaluatePriceFluctuation(stock)
	comboScore := rs.evaluateActionRatingCombination(stock, fluctuation)
	ratingQuality := rs.evaluateRatingQuality(stock.RatingTo)
	ratingChange := rs.evaluateRatingChange(stock)
	activityScore := rs.evaluateRecentActivity(stock)
	consensusScore := rs.evaluateHistoryConsensus(history)

	dataQuality := targetQuality
	if strings.TrimSpace(stock.RatingTo) != "" {
		dataQuality += 20
	}
	if strings.TrimSpace(stock.RatingFrom) != "" {
		dataQuality += 8
	}
	if strings.TrimSpace(stock.Action) != "" {
		dataQuality += 15
	}
	if !stock.Time.IsZero() {
		dataQuality += 12
	}
	if fluctuation.HasPricePair && math.Abs(fluctuation.PercentChange) >= 2 {
		dataQuality += 10
	}

	return featureVector{
		PriceDirection:    fluctuation.Direction,
		PriceMomentum:     fluctuation.Momentum,
		ActionRatingCombo: comboScore,
		RatingQuality:     ratingQuality,
		RatingChange:      ratingChange,
		RecentActivity:    activityScore,
		HistoryConsensus:  consensusScore,
		DataQuality:       clamp(dataQuality, 0, 100),
	}
}

func (rs *RecommendationService) scoreWithWeights(features featureVector, stock models.Stock) ([]weightedFeature, float64) {
	weighted := []weightedFeature{
		{label: "Price direction", weight: 0.36, value: features.PriceDirection, positiveText: "Target price moved upward", neutralText: "Target price is mostly unchanged", negativeText: "Target price moved downward"},
		{label: "Price momentum", weight: 0.22, value: features.PriceMomentum, positiveText: "Magnitude of target change is bullish", neutralText: "Target change magnitude is small", negativeText: "Magnitude of target change is bearish"},
		{label: "Action-rating combo", weight: 0.20, value: features.ActionRatingCombo, positiveText: fmt.Sprintf("Action and rating are aligned (%s / %s→%s)", stock.Action, stock.RatingFrom, stock.RatingTo), neutralText: "Action and rating combination is mixed", negativeText: fmt.Sprintf("Action and rating are bearish (%s / %s→%s)", stock.Action, stock.RatingFrom, stock.RatingTo)},
		{label: "Rating quality", weight: 0.08, value: features.RatingQuality, positiveText: fmt.Sprintf("Current rating is favorable (%s)", stock.RatingTo), neutralText: "Current rating is neutral", negativeText: fmt.Sprintf("Current rating is weak (%s)", stock.RatingTo)},
		{label: "Rating change", weight: 0.06, value: features.RatingChange, positiveText: "Rating improved", neutralText: "Rating is unchanged", negativeText: "Rating deteriorated"},
		{label: "Recency", weight: 0.05, value: features.RecentActivity, positiveText: "Very recent signal", neutralText: "Moderately recent signal", negativeText: "Signal is stale"},
		{label: "Consensus", weight: 0.03, value: features.HistoryConsensus, positiveText: "Recent history confirms bullish bias", neutralText: "Recent history is mixed", negativeText: "Recent history confirms bearish bias"},
	}

	score := 0.0
	for i := range weighted {
		weighted[i].contribution = weighted[i].value * weighted[i].weight
		score += weighted[i].contribution
	}

	return weighted, score
}

func (rs *RecommendationService) buildReason(weighted []weightedFeature) string {
	parts := make([]string, 0, len(weighted))
	for _, item := range weighted {
		explanation := item.label
		if item.value >= 15 {
			explanation = item.positiveText
		} else if item.value <= -15 {
			explanation = item.negativeText
		} else if item.neutralText != "" {
			explanation = item.neutralText
		}

		parts = append(parts, fmt.Sprintf("%s (+%.1f points)", explanation, item.contribution))
	}

	return strings.Join(parts, ". ")
}

func (rs *RecommendationService) evaluatePriceFluctuation(stock models.Stock) (priceFluctuation, float64) {
	fromPrice := rs.parsePrice(stock.TargetFrom)
	toPrice := rs.parsePrice(stock.TargetTo)
	actionSignal := rs.evaluateActionSignal(stock.Action)

	if fromPrice == 0 || toPrice == 0 {
		fallbackDirection := actionSignal * 0.35
		fallbackMomentum := actionSignal * 0.25
		return priceFluctuation{
			Direction:     clamp(fallbackDirection, -100, 100),
			Momentum:      clamp(fallbackMomentum, -100, 100),
			PercentChange: 0,
			HasPricePair:  false,
		}, 25
	}

	percentChange := ((toPrice - fromPrice) / fromPrice) * 100
	direction := math.Tanh(percentChange/14.0) * 100
	momentum := percentChange * 4.5

	if math.Abs(percentChange) < 0.6 {
		direction = actionSignal * 0.15
		momentum = actionSignal * 0.10
	}

	return priceFluctuation{
		Direction:     clamp(direction, -100, 100),
		Momentum:      clamp(momentum, -100, 100),
		PercentChange: percentChange,
		HasPricePair:  true,
	}, 70
}

func (rs *RecommendationService) evaluateActionRatingCombination(stock models.Stock, fluctuation priceFluctuation) float64 {
	actionSignal := rs.evaluateActionSignal(stock.Action)
	ratingToScore := rs.ratingToScore(stock.RatingTo)
	ratingFromScore := rs.ratingToScore(stock.RatingFrom)
	ratingSignal := rs.evaluateRatingQuality(stock.RatingTo)
	ratingTransition := clamp((ratingToScore-ratingFromScore)*35, -100, 100)

	combo := actionSignal*0.44 + ratingSignal*0.32 + ratingTransition*0.24

	if actionSignal >= 55 && ratingTransition >= 20 {
		combo += 18
	}
	if actionSignal <= -55 && ratingTransition <= -20 {
		combo -= 18
	}

	if math.Abs(ratingTransition) < 10 && math.Abs(actionSignal) >= 70 {
		if actionSignal > 0 {
			if ratingToScore >= 3.5 {
				combo += 14
			} else if ratingToScore <= 1.5 {
				combo -= 20
			} else {
				combo += 5
			}
		} else {
			if ratingToScore <= 1.5 {
				combo -= 14
			} else if ratingToScore >= 3.5 {
				combo += 12
			} else {
				combo -= 5
			}
		}
	}

	if actionSignal >= 55 && ratingToScore <= 1.5 {
		combo -= 28
	}
	if actionSignal <= -55 && ratingToScore >= 3.8 {
		combo += 20
	}

	if fluctuation.HasPricePair {
		alignment := fluctuation.Direction * actionSignal
		if alignment > 1500 {
			combo += 10
		} else if alignment < -1500 {
			combo -= 12
		}
	}

	return clamp(combo, -100, 100)
}

// evaluateActionSignal transforma la metadata de action en una señal de mercado [-100, 100]
func (rs *RecommendationService) evaluateActionSignal(action string) float64 {
	normalized := normalizeText(action)
	if normalized == "" {
		return 0
	}

	bestScore := 0.0
	bestPatternLen := 0
	for _, signal := range actionSignals {
		if strings.Contains(normalized, signal.pattern) {
			if len(signal.pattern) > bestPatternLen {
				bestPatternLen = len(signal.pattern)
				bestScore = signal.score
			}
		}
	}

	return bestScore
}

// evaluateRatingChange evalúa el cambio en el rating
func (rs *RecommendationService) evaluateRatingChange(stock models.Stock) float64 {
	fromRating := rs.ratingToScore(stock.RatingFrom)
	toRating := rs.ratingToScore(stock.RatingTo)

	diff := toRating - fromRating

	return clamp(diff*35, -100, 100)
}

func (rs *RecommendationService) evaluateRatingQuality(rating string) float64 {
	score := rs.ratingToScore(rating)

	// Mapea de [0,5] a [-100,100]
	return ((score - 2.5) / 2.5) * 100
}

// evaluateRecentActivity evalúa la actividad reciente del stock
func (rs *RecommendationService) evaluateRecentActivity(stock models.Stock) float64 {
	if stock.Time.IsZero() {
		return -40
	}

	daysSinceUpdate := rs.now().Sub(stock.Time).Hours() / 24

	if daysSinceUpdate < 3 {
		return 90
	} else if daysSinceUpdate < 7 {
		return 70
	} else if daysSinceUpdate < 14 {
		return 35
	} else if daysSinceUpdate < 30 {
		return -10
	} else {
		return -50
	}
}

// ratingToScore convierte un rating en un valor numérico
func (rs *RecommendationService) ratingToScore(rating string) float64 {
	rating = normalizeText(rating)

	if rating == "" {
		return defaultUnknownRating
	}

	bestScore := defaultUnknownRating
	bestPatternLen := 0
	for _, signal := range ratingSignals {
		if strings.Contains(rating, signal.pattern) {
			if len(signal.pattern) > bestPatternLen {
				bestPatternLen = len(signal.pattern)
				bestScore = signal.score
			}
		}
	}

	return bestScore
}

// parsePrice extrae el valor numérico de un precio (ej: "$150.00" -> 150.00)
func (rs *RecommendationService) parsePrice(price string) float64 {
	clean := strings.TrimSpace(price)
	if clean == "" {
		return 0
	}
	clean = strings.ReplaceAll(clean, "$", "")
	clean = strings.ReplaceAll(clean, "€", "")
	clean = strings.ReplaceAll(clean, "£", "")
	clean = strings.ReplaceAll(clean, ",", "")
	clean = strings.ReplaceAll(clean, "USD", "")
	clean = strings.TrimSpace(clean)

	value, err := strconv.ParseFloat(clean, 64)
	if err != nil {
		return 0
	}
	return value
}

func (rs *RecommendationService) determineConfidence(score float64, historyCount int, dataQuality float64) models.ConfidenceLevel {
	if score >= 74 && historyCount >= 4 && dataQuality >= 72 {
		return models.ConfidenceHigh
	} else if score >= 58 && historyCount >= 2 && dataQuality >= 48 {
		return models.ConfidenceMedium
	} else {
		return models.ConfidenceLow
	}
}

func (rs *RecommendationService) evaluateHistoryConsensus(history []models.Stock) float64 {
	if len(history) == 0 {
		return 0
	}

	maxItems := maxHistoryItems
	if len(history) < maxItems {
		maxItems = len(history)
	}

	total := 0.0
	for i := 0; i < maxItems; i++ {
		item := history[i]
		recencyWeight := 1.0 - (float64(i) * 0.1)
		if recencyWeight < 0.3 {
			recencyWeight = 0.3
		}
		fluctuation, _ := rs.evaluatePriceFluctuation(item)
		combo := rs.evaluateActionRatingCombination(item, fluctuation)
		combined := fluctuation.Direction*0.45 + fluctuation.Momentum*0.20 + combo*0.25 + rs.evaluateRatingQuality(item.RatingTo)*0.10
		total += combined * recencyWeight
	}

	divisor := 0.0
	for i := 0; i < maxItems; i++ {
		weight := 1.0 - (float64(i) * 0.1)
		if weight < 0.3 {
			weight = 0.3
		}
		divisor += weight
	}

	avg := total / divisor
	return clamp(avg, -100, 100)
}

func (rs *RecommendationService) sortRecommendations(recs []models.StockRecommendation) []models.StockRecommendation {
	sort.Slice(recs, func(i, j int) bool {
		return recs[i].Score > recs[j].Score
	})
	return recs
}

func normalizeText(value string) string {
	clean := strings.ToLower(strings.TrimSpace(value))
	if clean == "" {
		return ""
	}
	clean = strings.ReplaceAll(clean, "-", " ")
	clean = noiseTextRegex.ReplaceAllString(clean, " ")
	clean = strings.Join(strings.Fields(clean), " ")
	return clean
}

func (rs *RecommendationService) calibrateScore(rawScore float64, dataQuality float64) float64 {
	qualityFactor := 0.78 + (clamp(dataQuality, 0, 100) / 100.0 * 0.22)
	adjusted := rawScore * qualityFactor
	return clamp((adjusted+100.0)/2.0, 0, 100)
}

func clamp(value, minValue, maxValue float64) float64 {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}
