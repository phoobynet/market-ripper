package equity

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"github.com/phoobynet/market-ripper/bar/models"
)

type StreamBarAdapter = func(b stream.Bar) models.Bar

func Adapt(b stream.Bar) models.Bar {
	return models.Bar{
		Class:      "e",
		Symbol:     b.Symbol,
		Open:       b.Open,
		High:       b.High,
		Low:        b.Low,
		Close:      b.Close,
		Volume:     float64(b.Volume),
		VWAP:       b.VWAP,
		TradeCount: float64(b.TradeCount),
		Timestamp:  b.Timestamp,
	}
}
