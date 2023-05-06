package crypto

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"github.com/phoobynet/market-ripper/bar"
)

type StreamBarAdapter = func(b stream.CryptoBar) bar.Bar

func Adapt(b stream.CryptoBar) bar.Bar {
	return bar.Bar{
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
