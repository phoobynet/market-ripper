package snapshot

import (
	"time"
)

type Snapshot struct {
	Symbol            string    `json:"S"`
	PreviousOpen      float64   `json:"pdo"`
	PreviousHigh      float64   `json:"pdh"`
	PreviousLow       float64   `json:"pdl"`
	PreviousClose     float64   `json:"pdc"`
	PreviousVolume    float64   `json:"pdv"`
	PreviousTimestamp time.Time `json:"pdt"`
	DailyOpen         float64   `json:"do"`
	DailyHigh         float64   `json:"dh"`
	DailyLow          float64   `json:"dl"`
	DailyClose        float64   `json:"dc"`
	DailyVolume       float64   `json:"dv"`
	DailyTimestamp    time.Time `json:"dt"`
}
