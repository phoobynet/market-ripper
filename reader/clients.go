package reader

import (
	"context"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"log"
)

var stocksClient *stream.StocksClient
var cryptoClient *stream.CryptoClient
var alpacaClient *alpaca.Client
var marketDataClient *marketdata.Client

func StartClients() {
	stocksClient = stream.NewStocksClient(marketdata.SIP)
	err := stocksClient.Connect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}

	cryptoClient = stream.NewCryptoClient(marketdata.US)
	err = cryptoClient.Connect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}

	alpacaClient = alpaca.NewClient(alpaca.ClientOpts{})
	marketDataClient = marketdata.NewClient(marketdata.ClientOpts{})
}
