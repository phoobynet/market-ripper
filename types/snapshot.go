package types

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"log"
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

func FromSnapshot(symbol string, snapshot interface{}) *Snapshot {
	switch s := snapshot.(type) {
	case marketdata.CryptoSnapshot:
		return FromCryptoSnapshot(symbol, &s)
	case marketdata.Snapshot:
		return FromEquitySnapshot(symbol, &s)
	default:
		log.Fatalf("unknown snapshot type: %T", snapshot)
		return nil
	}
}

func FromCryptoSnapshot(symbol string, snapshot *marketdata.CryptoSnapshot) *Snapshot {
	if symbol == "" {
		return nil
	}

	if snapshot == nil {
		return nil
	}

	previousDailyBar := snapshot.PrevDailyBar

	if previousDailyBar == nil {
		return nil
	}

	dailyBar := snapshot.DailyBar

	if dailyBar == nil {
		return nil
	}

	return &Snapshot{
		Symbol:            symbol,
		PreviousOpen:      previousDailyBar.Open,
		PreviousHigh:      previousDailyBar.High,
		PreviousLow:       previousDailyBar.Low,
		PreviousClose:     previousDailyBar.Close,
		PreviousVolume:    previousDailyBar.Volume,
		PreviousTimestamp: previousDailyBar.Timestamp,
		DailyOpen:         dailyBar.Open,
		DailyHigh:         dailyBar.High,
		DailyLow:          dailyBar.Low,
		DailyClose:        dailyBar.Close,
		DailyVolume:       dailyBar.Volume,
		DailyTimestamp:    dailyBar.Timestamp,
	}
}

func FromEquitySnapshot(symbol string, snapshot *marketdata.Snapshot) *Snapshot {
	if symbol == "" {
		return nil
	}

	if snapshot == nil {
		return nil
	}

	previousDailyBar := snapshot.PrevDailyBar

	if previousDailyBar == nil {
		return nil
	}

	dailyBar := snapshot.DailyBar

	if dailyBar == nil {
		return nil
	}

	return &Snapshot{
		Symbol:            symbol,
		PreviousOpen:      previousDailyBar.Open,
		PreviousHigh:      previousDailyBar.High,
		PreviousLow:       previousDailyBar.Low,
		PreviousClose:     previousDailyBar.Close,
		PreviousVolume:    float64(previousDailyBar.Volume),
		PreviousTimestamp: previousDailyBar.Timestamp,
		DailyOpen:         dailyBar.Open,
		DailyHigh:         dailyBar.High,
		DailyLow:          dailyBar.Low,
		DailyClose:        dailyBar.Close,
		DailyVolume:       float64(dailyBar.Volume),
		DailyTimestamp:    dailyBar.Timestamp,
	}
}
