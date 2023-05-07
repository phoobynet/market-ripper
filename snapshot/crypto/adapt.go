package crypto

import (
	"errors"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/phoobynet/market-ripper/snapshot/models"
)

// Adapt adapts a crypto snapshot to a models.Snapshot
func Adapt(
	symbol string,
	cryptoSnapshot marketdata.CryptoSnapshot,
) (*models.Snapshot, error) {
	if symbol == "" {
		return nil, errors.New("symbol cannot be empty")
	}

	previousDailyBar := cryptoSnapshot.PrevDailyBar

	if previousDailyBar == nil {
		return nil, errors.New("unexpected prevDailyBar was empty")
	}

	dailyBar := cryptoSnapshot.DailyBar

	if dailyBar == nil {
		return nil, errors.New("unexpected dailyBar was empty")
	}

	return &models.Snapshot{
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
	}, nil
}
