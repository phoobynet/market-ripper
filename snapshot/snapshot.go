package snapshot

import (
	"gorm.io/gorm"
	"time"
)

type Snapshot struct {
	gorm.Model
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
