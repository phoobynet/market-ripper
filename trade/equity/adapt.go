package equity

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	. "github.com/phoobynet/market-ripper/trade"
)

type StreamTradeAdapter = func(t stream.Trade) Trade

func Adapt(t stream.Trade) Trade {
	return Trade{
		Class:     "e",
		Symbol:    t.Symbol,
		Price:     t.Price,
		Size:      float64(t.Size),
		Exchange:  t.Exchange,
		Timestamp: t.Timestamp,
	}
}
