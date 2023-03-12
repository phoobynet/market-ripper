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
	defer query.Disconnect()

	reader.StartClients()

	// ASSET loading...
	assetRepository := query.NewAssetRepository()

	if assetRepository.Count() == 0 || assetRepository.IsStale(-24*time.Hour) {
		assetReader := reader.NewAssetReader()
		assets := assetReader.ReadAllActive()

		assetWriter := writer.NewAssetWriter(configuration)
		defer assetWriter.Close()
		assetWriter.Write(assets)
	}

	// SNAPSHOTS loading...
	snapshotsRepository := query.NewSnapshotRepository(configuration)
	snapshotsRepository.Truncate()

	snapshots := reader.NewSnapshotReader(configuration, assetRepository).Read()
	snapshotWriter := writer.NewSnapshotWriter(configuration)
	defer snapshotWriter.Close()
	snapshotWriter.Write(snapshots)

	barWriter := writer.NewBarWriter(configuration)
	defer barWriter.Close()

	tradeWriter := writer.NewTradeWriter(configuration)
	defer tradeWriter.Close()

	var streamingTradesChan = make(chan types.Trade, 100_000)
	var streamingBarsChan = make(chan types.Bar, 20_000)

	if configuration.Class == alpaca.Crypto {
		cryptoReader := reader.NewCryptoReader(configuration)
		defer cryptoReader.Unsubscribe()
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
	} else if configuration.Class == alpaca.USEquity {
		equityReader := reader.NewEquityTradeReader(configuration)
		equityReader.Unsubscribe()
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
	}
}
