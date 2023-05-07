package equity

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"github.com/phoobynet/market-ripper/trade"
)

// Adapt adapts a crypto trade to a trade.Trade
func Adapt(t stream.CryptoTrade) trade.Trade {
	return trade.Trade{
		Class:     "c",
		Symbol:    t.Symbol,
		Price:     t.Price,
		Size:      t.Size,
		TakerSide: t.TakerSide,
		Exchange:  t.Exchange,
		Timestamp: t.Timestamp,
	}
}
