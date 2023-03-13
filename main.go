package main

import (
	"flag"
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

	flag.StringVar(
		&configurationFile,
		"config",
		"config.toml",
		"Configuration file",
	)

	quit := make(
		chan os.Signal,
		1,
	)
	signal.Notify(
		quit,
		os.Interrupt,
	)

	configuration := config.Load(configurationFile)
	log.Printf(
		"%s",
		configuration,
	)

	query.Connect(configuration)
	defer query.Disconnect()

	reader.StartClients()

	// ASSET loading...
	assetRepository := query.NewAssetRepository()

	if assetRepository.Count() == 0 || assetRepository.IsStale(-24*time.Hour) {
		assetReader := reader.NewAssetReader()
		assets := assetReader.GetActive()

		assetWriter := writer.NewAssetWriter(configuration)
		defer assetWriter.Close()
		assetWriter.Write(assets)
	}

	// SNAPSHOTS loading...
	snapshotsRepository := query.NewSnapshotRepository(configuration)
	snapshotsRepository.Truncate()

	snapshotReader := reader.NewSnapshotReader(
		configuration,
		assetRepository,
	)
	snapshots := snapshotReader.Read()
	snapshotWriter := writer.NewSnapshotWriter(configuration)
	defer snapshotWriter.Close()
	snapshotWriter.Write(snapshots)

	// TRADES and BARS readers...
	barWriter := writer.NewBarWriter(configuration)
	defer barWriter.Close()

	tradeWriter := writer.NewTradeWriter(configuration)
	defer tradeWriter.Close()

	// Channels and tickers...
	var snapshotRefreshTimer = time.NewTicker(10 * time.Minute)
	defer snapshotRefreshTimer.Stop()

	var streamingTradesChan = make(
		chan types.Trade,
		100_000,
	)

	var streamingBarsChan = make(
		chan types.Bar,
		20_000,
	)

	// Class reader
	classReader := reader.CreateClassReader(configuration)
	defer classReader.Unsubscribe()
	go func() {
		classReader.Subscribe(
			streamingTradesChan,
			streamingBarsChan,
		)

		for {
			select {
			case t := <-streamingTradesChan:
				tradeWriter.Write(t)
			case b := <-streamingBarsChan:
				barWriter.Write(b)
			case <-snapshotRefreshTimer.C:
				go func() {
					log.Println("Refreshing snapshots...")
					snapshotsRepository.Truncate()
					snapshots = snapshotReader.Read()
					snapshotWriter.Write(snapshots)
				}()
			}
		}
	}()
	<-quit
}
