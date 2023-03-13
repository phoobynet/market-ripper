package reader

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/phoobynet/market-ripper/config"
	"github.com/phoobynet/market-ripper/types"
	"log"
)

type ClassReader interface {
	Subscribe(
		streamingTradesChan chan types.Trade,
		streamingBarsChan chan types.Bar,
	)

	Unsubscribe()
}

func CreateClassReader(configuration *config.Config) ClassReader {
	if configuration.Class == alpaca.USEquity {
		return NewEquityTradeReader(configuration)
	} else if configuration.Class == alpaca.Crypto {
		return NewCryptoReader(configuration)
	} else {
		log.Fatalln("Invalid class: " + configuration.Class)
	}

	return nil
}
