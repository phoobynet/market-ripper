package types

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"time"
)

type Trade struct {
	Class     string    `json:"class"`
	Symbol    string    `json:"S"`
	Price     float64   `json:"p"`
	Size      float64   `json:"s"`
	TakerSide string    `json:"takerSide"`
	Exchange  string    `json:"e"`
	Timestamp time.Time `json:"t"`
}

func FromCryptoTrade(t stream.CryptoTrade) Trade {
	return Trade{
		Class:     "c",
		Symbol:    t.Symbol,
		Price:     t.Price,
		Size:      t.Size,
		TakerSide: t.TakerSide,
		Exchange:  t.Exchange,
		Timestamp: t.Timestamp,
	}
}

func FromEquityTrade(t stream.Trade) Trade {
	return Trade{
		Class:     "e",
		Symbol:    t.Symbol,
		Price:     t.Price,
		Size:      float64(t.Size),
		Exchange:  t.Exchange,
		Timestamp: t.Timestamp,
	}
}
