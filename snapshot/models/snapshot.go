package models

import (
	"time"
)

// Snapshot is a snapshot of a stock or crypto
type Snapshot struct {
	Symbol            string
	PreviousOpen      float64
	PreviousHigh      float64
	PreviousLow       float64
	PreviousClose     float64
	PreviousVolume    float64
	PreviousTimestamp time.Time
	DailyOpen         float64
	DailyHigh         float64
	DailyLow          float64
	DailyClose        float64
	DailyVolume       float64
	DailyTimestamp    time.Time
}
