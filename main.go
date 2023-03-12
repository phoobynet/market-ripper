package main

import (
	"flag"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/phoobynet/market-ripper/config"
	"github.com/phoobynet/market-ripper/query"
	"github.com/phoobynet/market-ripper/reader"
	"github.com/phoobynet/market-ripper/types"
	"github.com/phoobynet/market-ripper/writer"
	"log"
	"os"
	"os/signal"
	"time"
)

var configurationFile string

func main() {
	config.ValidateEnv()

	flag.StringVar(&configurationFile, "config", "config.toml", "Configuration file")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	configuration := config.Load(configurationFile)
	log.Printf("%s", configuration)

	query.Connect(configuration)
	reader.StartClients()
	writer.StartLineSender(configuration)

	// ASSET loading...
	assetRepository := query.NewAssetRepository()

	if assetRepository.Count() == 0 || assetRepository.IsStale(-24*time.Hour) {
		assetReader := reader.NewAssetReader()
		assets := assetReader.ReadAllActive()

		assetWriter := writer.NewAssetWriter()
		assetWriter.Write(assets)
	}

	// SNAPSHOTS loading...
	snapshotsRepository := query.NewSnapshotRepository(configuration)
	snapshotsRepository.Truncate()

	snapshots := reader.NewSnapshotReader(configuration, assetRepository).Read()
	writer.NewSnapshotWriter(configuration).Write(snapshots)

	barWriter := writer.NewBarWriter(configuration)
	tradeWriter := writer.NewTradeWriter(configuration)
	var streamingTradesChan = make(chan types.Trade, 100_000)
	var streamingBarsChan = make(chan types.Bar, 20_000)

	if configuration.Class == alpaca.Crypto {
		cryptoReader := reader.NewCryptoReader(configuration)
		go func() {
			cryptoReader.Subscribe(streamingTradesChan, streamingBarsChan)

			for {
				select {
				case t := <-streamingTradesChan:
					tradeWriter.Write(t)
				case b := <-streamingBarsChan:
					barWriter.Write(b)
				}
			}
		}()
		<-quit
		cryptoReader.Unsubscribe()
	} else if configuration.Class == alpaca.USEquity {
		equityReader := reader.NewEquityTradeReader(configuration)
		go func() {
			equityReader.Subscribe(streamingTradesChan, streamingBarsChan)

			for {
				select {
				case t := <-streamingTradesChan:
					tradeWriter.Write(t)
				case b := <-streamingBarsChan:
					barWriter.Write(b)
				}
			}
		}()
		<-quit
		equityReader.Unsubscribe()
	}

	query.Disconnect()
	writer.CloseLineSender()
}
