package reader

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"github.com/phoobynet/market-ripper/config"
	"github.com/phoobynet/market-ripper/types"
	"log"
)

type CryptoReader struct {
	configuration *config.Config
}

func NewCryptoReader(configuration *config.Config) *CryptoReader {
	return &CryptoReader{
		configuration: configuration,
	}
}

func (t *CryptoReader) Subscribe(streamingTradesChan chan types.Trade, streamingBarsChan chan types.Bar) {
	err := cryptoClient.SubscribeToTrades(func(t stream.CryptoTrade) {
		streamingTradesChan <- types.FromCryptoTrade(t)
	}, t.configuration.Symbols...)

	if err != nil {
		log.Fatal(err)
	}

	err = cryptoClient.SubscribeToBars(func(b stream.CryptoBar) {
		streamingBarsChan <- types.FromCryptoBar(b)
	}, t.configuration.Symbols...)

	if err != nil {
		log.Fatal(err)
	}
}

func (t *CryptoReader) Unsubscribe() {
	err := cryptoClient.UnsubscribeFromTrades(t.configuration.Symbols...)

	if err != nil {
		log.Fatal(err)
	}

	err = cryptoClient.UnsubscribeFromBars(t.configuration.Symbols...)

	if err != nil {
		log.Fatal(err)
	}
}
