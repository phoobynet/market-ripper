package clients

import (
	"context"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
)

type Clients struct {
	stocksStream *stream.StocksClient
	cryptoStream *stream.CryptoClient
	alpaca       *alpaca.Client
	marketData   *marketdata.Client
}

func NewClients(ctx context.Context) (*Clients, error) {
	var c *Clients

	c.stocksStream = stream.NewStocksClient(marketdata.SIP)
	err := c.stocksStream.Connect(ctx)

	if err != nil {
		return nil, err
	}

	c.cryptoStream = stream.NewCryptoClient(marketdata.US)
	err = c.cryptoStream.Connect(ctx)

	if err != nil {
		return nil, err
	}

	c.alpaca = alpaca.NewClient(alpaca.ClientOpts{})
	c.marketData = marketdata.NewClient(marketdata.ClientOpts{})

	return c, nil
}

func (c *Clients) StocksStream() *stream.StocksClient {
	return c.stocksStream
}

func (c *Clients) CryptoStream() *stream.CryptoClient {
	return c.cryptoStream
}

func (c *Clients) Alpaca() *alpaca.Client {
	return c.alpaca
}

func (c *Clients) MarketData() *marketdata.Client {
	return c.marketData
}
