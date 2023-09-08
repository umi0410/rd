package domain

import (
	"time"
)

type Alias struct {
	ID              uint           `json:"id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       time.Time      `json:"deleted_at"`
	AliasGroup      string         `json:"alias_group"`
	Name            string         `json:"name"`
	Destination     string         `json:"destination"`
	RecentHitCounts map[string]int `json:"recent_hit_counts"`
}
