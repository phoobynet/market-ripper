package equity

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"github.com/phoobynet/market-ripper/trade"
)

type StreamTradeAdapter = func(t stream.CryptoTrade) trade.Trade

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
