package crypto

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"github.com/phoobynet/market-ripper/bar/models"
	"github.com/phoobynet/market-ripper/config"
)

// Reader is a bar reader for crypto bars
type Reader struct {
	configuration *config.Config
	client        *stream.CryptoClient
}

// NewReader creates a new crypto bar reader
func NewReader(configuration *config.Config, client *stream.CryptoClient) (*Reader, error) {
	reader := &Reader{
		configuration: configuration,
		client:        client,
	}

	return reader, nil
}

// Subscribe subscribes to crypto bars
func (r *Reader) Subscribe(bars chan models.Bar) error {
	return r.client.SubscribeToBars(
		func(b stream.CryptoBar) {
			bars <- Adapt(b)
		},
		r.configuration.Symbols...,
	)
}

// Unsubscribe unsubscribes from crypto bars
func (r *Reader) Unsubscribe() error {
	return r.client.UnsubscribeFromBars(r.configuration.Symbols...)
}
