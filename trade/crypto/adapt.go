package equity

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	. "github.com/phoobynet/market-ripper/trade"
)

type StreamTradeAdapter = func(t stream.CryptoTrade) Trade

func Adapt(t stream.CryptoTrade) Trade {
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
