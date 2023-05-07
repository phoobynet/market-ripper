package equity

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"github.com/phoobynet/market-ripper/config"
	"github.com/phoobynet/market-ripper/trade"
)

// Reader is an crypto trade reader
type Reader struct {
	configuration *config.Config
	client        *stream.CryptoClient
}

// NewReader creates a new crypto trade reader
func NewReader(configuration *config.Config, client *stream.CryptoClient) (*Reader, error) {
	reader := &Reader{
		configuration: configuration,
		client:        client,
	}

	return reader, nil
}

// Subscribe subscribes to crypto trades
func (r *Reader) Subscribe(
	out chan trade.Trade,
) error {
	return r.client.SubscribeToTrades(
		func(t stream.CryptoTrade) {
			out <- Adapt(t)
		},
		r.configuration.Symbols...,
	)
}

// Unsubscribe unsubscribes from crypto trades
func (r *Reader) Unsubscribe() error {
	return r.client.UnsubscribeFromTrades(r.configuration.Symbols...)
}
