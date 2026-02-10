package models

import (
	"time"
)

type Stock struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
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
	Stock      Stock   `json:"stock"`
	Score      float64 `json:"score"`
	Reason     string  `json:"reason"`
	Confidence string  `json:"confidence"` // "high", "medium", "low"
}
