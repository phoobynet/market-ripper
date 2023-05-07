package equity

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"github.com/phoobynet/market-ripper/trade"
)

type StreamTradeAdapter = func(t stream.Trade) trade.Trade

func Adapt(t stream.Trade) trade.Trade {
	return trade.Trade{
		Class:     "e",
		Symbol:    t.Symbol,
		Price:     t.Price,
		Size:      float64(t.Size),
		Exchange:  t.Exchange,
		Timestamp: t.Timestamp,
	}
}
