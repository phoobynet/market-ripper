package main

import (
	"context"
	"flag"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"github.com/phoobynet/market-ripper/config"
	"github.com/phoobynet/market-ripper/query"
	"github.com/phoobynet/market-ripper/reader"
	"github.com/phoobynet/market-ripper/snapshots"
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

	assetRepository := query.NewAssetRepository()

	if assetRepository.Count() == 0 || assetRepository.IsStale(-24*time.Hour) {
		assetReader := reader.NewAssetReader()
		assets := assetReader.ReadAllActive()

		assetWriter := writer.NewAssetWriter()
		assetWriter.Write(assets)
	}

	snapshots.Load(context.TODO(), configuration, allSymbols)

	cryptoReaderCtx, readerCancel := context.WithCancel(context.Background())
	cryptoReader := reader.NewCryptoReader(cryptoReaderCtx, configuration)

	writerCtx, writerCancel := context.WithCancel(context.Background())

	sipWriter, err := writer(writerCtx, configuration)

	if err != nil {
		log.Fatal(err)
	}

	var streamingTradesChan = make(chan stream.Trade, 100_000)
	var streamingBarsChan = make(chan stream.Bar, 20_000)

	go func() {
		err := sipReader.Observe(streamingTradesChan, streamingBarsChan)

		if err != nil {
			log.Fatal(err)
		} else {
			log.Println("SIP observer started")
		}

		for {
			select {
			case t := <-streamingTradesChan:
				sipWriter.WriteTrade(t)
			case b := <-streamingBarsChan:
				sipWriter.WriteBar(b)
			case <-readerCtx.Done():
				log.Println("Shutting down SIP Observer reader")
				_ = sipReader.Disconnect()
				return
			case <-writerCtx.Done():
				log.Println("Shutting down SIP Observer writer")
				_ = sipWriter.Close()
				return
			}
		}
	}()

	<-quit
	readerCancel()
	writerCancel()
}
