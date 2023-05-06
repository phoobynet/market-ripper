package equity

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"github.com/phoobynet/market-ripper/config"
	"github.com/phoobynet/market-ripper/trade"
)

type Reader struct {
	configuration *config.Config
	client        *stream.StocksClient
}

func NewReader(configuration *config.Config, client *stream.StocksClient) (*Reader, error) {
	reader := &Reader{
		configuration: configuration,
		client:        client,
	}

	return reader, nil
}

func (r *Reader) Subscribe(
	out chan trade.Trade,
) error {
	return r.client.SubscribeToTrades(
		func(t stream.Trade) {
			out <- Adapt(t)
		},
		r.configuration.Symbols...,
	)
}

func (r *Reader) Unsubscribe() error {
	return r.client.UnsubscribeFromTrades(r.configuration.Symbols...)
}
