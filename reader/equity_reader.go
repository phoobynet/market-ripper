package reader

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"github.com/phoobynet/market-ripper/config"
	"github.com/phoobynet/market-ripper/types"
	"log"
)

type EquityTradeReader struct {
	configuration *config.Config
}

func NewEquityTradeReader(configuration *config.Config) *EquityTradeReader {
	return &EquityTradeReader{
		configuration: configuration,
	}
}

func (t *EquityTradeReader) Subscribe(
	streamingTradesChan chan types.Trade,
	streamingBarsChan chan types.Bar,
) {
	err := stocksClient.SubscribeToTrades(
		func(t stream.Trade) {
			streamingTradesChan <- types.FromEquityTrade(t)
		},
		t.configuration.Symbols...,
	)

	if err != nil {
		log.Fatal(err)
	}

	err = stocksClient.SubscribeToBars(
		func(b stream.Bar) {
			streamingBarsChan <- types.FromEquityBar(b)
		},
		t.configuration.Symbols...,
	)

	if err != nil {
		log.Fatal(err)
	}
}

func (t *EquityTradeReader) Unsubscribe() {
	err := stocksClient.UnsubscribeFromTrades(t.configuration.Symbols...)

	if err != nil {
		log.Fatal(err)
	}

	err = stocksClient.UnsubscribeFromBars(t.configuration.Symbols...)

	if err != nil {
		log.Fatal(err)
	}
}
