package types

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"time"
)

type Bar struct {
	Class      string    `json:"class"`
	Symbol     string    `json:"S"`
	Open       float64   `json:"o"`
	High       float64   `json:"h"`
	Low        float64   `json:"l"`
	Close      float64   `json:"c"`
	Volume     float64   `json:"v"`
	VWAP       float64   `json:"vw"`
	TradeCount float64   `json:"n"`
	Timestamp  time.Time `json:"t"`
}

func FromCryptoBar(b stream.CryptoBar) Bar {
	return Bar{
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

func FromEquityBar(b stream.Bar) Bar {
	return Bar{
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
