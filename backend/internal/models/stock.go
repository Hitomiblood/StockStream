package models

import (
	"encoding/json"
	"time"
)

// ConfidenceLevel representa el nivel de confianza de una recomendación
type ConfidenceLevel int

const (
	ConfidenceLow ConfidenceLevel = iota
	ConfidenceMedium
	ConfidenceHigh
)

// String implementa el método Stringer para ConfidenceLevel
func (c ConfidenceLevel) String() string {
	switch c {
	case ConfidenceLow:
		return "low"
	case ConfidenceMedium:
		return "medium"
	case ConfidenceHigh:
		return "high"
	default:
		return "unknown"
	}
}

// MarshalJSON implementa la interfaz json.Marshaler
func (c ConfidenceLevel) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

// UnmarshalJSON implementa la interfaz json.Unmarshaler
func (c *ConfidenceLevel) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	switch s {
	case "low":
		*c = ConfidenceLow
	case "medium":
		*c = ConfidenceMedium
	case "high":
		*c = ConfidenceHigh
	default:
		*c = ConfidenceLow // valor por defecto
	}
	return nil
}

type Stock struct {
	ID         uint64    `gorm:"primaryKey" json:"id,string"`
	Ticker     string    `gorm:"index;not null" json:"ticker"`
	TargetFrom string    `json:"target_from"`
	TargetTo   string    `json:"target_to"`
	Company    string    `json:"company"`
	Action     string    `json:"action"`
	Brokerage  string    `json:"brokerage"`
	RatingFrom string    `json:"rating_from"`
	RatingTo   string    `json:"rating_to"`
	Time       time.Time `gorm:"index" json:"time"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type APIResponse struct {
	Items    []Stock `json:"items"`
	NextPage string  `json:"next_page"`
}

type StockRecommendation struct {
	Stock      Stock           `json:"stock"`
	Score      float64         `json:"score"`
	Reason     string          `json:"reason"`
	Confidence ConfidenceLevel `json:"confidence"`
}
