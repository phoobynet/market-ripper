package equity

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/phoobynet/market-ripper/snapshot/models"
)

// Adapt adapts an equity snapshot to a models.Snapshot
func Adapt(
	symbol string,
	equitySnapshot *marketdata.Snapshot,
) *models.Snapshot {
	if symbol == "" {
		return nil
	}

	if equitySnapshot == nil {
		return nil
	}

	previousDailyBar := equitySnapshot.PrevDailyBar

	if previousDailyBar == nil {
		return nil
	}

	dailyBar := equitySnapshot.DailyBar

	if dailyBar == nil {
		return nil
	}

	return &models.Snapshot{
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
