package equity

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"github.com/phoobynet/market-ripper/bar/models"
	"github.com/phoobynet/market-ripper/config"
)

// Reader is a bar reader for equity bars
type Reader struct {
	configuration *config.Config
	client        *stream.StocksClient
}

// NewReader creates a new equity bar reader
func NewReader(configuration *config.Config, client *stream.StocksClient) (*Reader, error) {
	reader := &Reader{
		configuration: configuration,
		client:        client,
	}

	return reader, nil
}

// Subscribe subscribes to equity bars
func (r *Reader) Subscribe(bars chan models.Bar) error {
	return r.client.SubscribeToBars(
		func(equityBar stream.Bar) {
			bars <- Adapt(equityBar)
		},
		r.configuration.Symbols...,
	)
}

// Unsubscribe unsubscribes from equity bars
func (r *Reader) Unsubscribe() error {
	return r.client.UnsubscribeFromBars(r.configuration.Symbols...)
}
