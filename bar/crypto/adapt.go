package crypto

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"github.com/phoobynet/market-ripper/bar/models"
)

// Adapt adapts a stream.CryptoBar to a models.Bar
func Adapt(b stream.CryptoBar) models.Bar {
	return models.Bar{
		Class:      "c",
		Symbol:     b.Symbol,
		Open:       b.Open,
		High:       b.High,
		Low:        b.Low,
		Close:      b.Close,
		Volume:     b.Volume,
		VWAP:       b.VWAP,
		TradeCount: float64(b.TradeCount),
		Timestamp:  b.Timestamp,
	}
}
