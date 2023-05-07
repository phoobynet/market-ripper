package main

import (
	"context"
	"flag"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"github.com/phoobynet/market-ripper/asset"
	"github.com/phoobynet/market-ripper/bar"
	barCrypto "github.com/phoobynet/market-ripper/bar/crypto"
	barEquity "github.com/phoobynet/market-ripper/bar/equity"
	"github.com/phoobynet/market-ripper/config"
	"github.com/phoobynet/market-ripper/database"
	"github.com/phoobynet/market-ripper/snapshot"
	"github.com/phoobynet/market-ripper/trade"
	tradeCrypto "github.com/phoobynet/market-ripper/trade/crypto"
	tradeEquity "github.com/phoobynet/market-ripper/trade/equity"
	"gorm.io/gorm"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	config.ValidateEnv()

	var configurationFile string

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

	configuration, err := config.Load(configurationFile)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf(
		"%s",
		configuration,
	)

	alpacaClient := alpaca.NewClient(alpaca.ClientOpts{})
	marketDataClient := marketdata.NewClient(marketdata.ClientOpts{})

	stocksClientContext, stocksClientCancel := context.WithCancel(context.Background())
	defer stocksClientCancel()
	stocksClient := stream.NewStocksClient(marketdata.SIP)

	err = stocksClient.Connect(stocksClientContext)

	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Connected to stocks client")
	}

	cryptoClientContext, cryptoClientCancel := context.WithCancel(context.Background())
	defer cryptoClientCancel()
	cryptoClient := stream.NewCryptoClient("us")
	err = cryptoClient.Connect(cryptoClientContext)

	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Connected to crypto client")
	}

	pgConnection, err := database.Connect(configuration)

	if err != nil {
		log.Fatal(err)
	}

	defer func(pgConnection *gorm.DB) {
		db, err := pgConnection.DB()

		if err == nil {
			_ = db.Close()
		}
	}(pgConnection)

	assetRepository, err := asset.Prepare(pgConnection, alpacaClient)

	if err != nil {
		log.Fatal(err)
	}

	snapshotRepository, err := snapshot.NewRepository(pgConnection)

	if err != nil {
		log.Fatal(err)
	}

	snapshotUpdater, err := snapshot.NewUpdater(configuration, assetRepository, marketDataClient, snapshotRepository)

	barWriter := bar.NewWriter(configuration)
	defer barWriter.Close()

	tradeWriter := trade.NewWriter(configuration)
	defer tradeWriter.Close()

	var snapshotRefreshTimer *time.Ticker

	if configuration.SnapshotRefreshIntervalMins > 0 {
		interval := time.Duration(configuration.SnapshotRefreshIntervalMins)
		snapshotRefreshTimer = time.NewTicker(interval * time.Minute)
	} else {
		snapshotRefreshTimer = time.NewTicker(24 * time.Hour)
	}

	defer snapshotRefreshTimer.Stop()

	var tradesChannel = make(
		chan trade.Trade,
		100_000,
	)

	var barsChannel = make(
		chan bar.Bar,
		20_000,
	)

	var barReader bar.Reader
	var tradeReader trade.Reader

	if configuration.Class == "us_equity" {
		barReader, err = barEquity.NewReader(configuration, stocksClient)
		if err != nil {
			log.Fatal(err)
		}

		tradeReader, err = tradeEquity.NewReader(configuration, stocksClient)

		if err != nil {
			log.Fatal(err)
		}
	} else {
		barReader, err = barCrypto.NewReader(configuration, cryptoClient)
		if err != nil {
			log.Fatal(err)
		}

		tradeReader, err = tradeCrypto.NewReader(configuration, cryptoClient)

		if err != nil {
			log.Fatal(err)
		}
	}

	err = barReader.Subscribe(barsChannel)

	if err != nil {
		log.Fatal(err)
	}

	err = tradeReader.Subscribe(tradesChannel)

	defer func(barReader bar.Reader) {
		err := barReader.Unsubscribe()

		if err != nil {
			log.Fatal(err)
		}
	}(barReader)

	defer func(tradeReader trade.Reader) {
		err := tradeReader.Unsubscribe()

		if err != nil {
			log.Fatal(err)
		}
	}(tradeReader)

	go func() {
		for {
			select {
			case t := <-tradesChannel:
				tradeWriter.Write(t)
			case b := <-barsChannel:
				barWriter.Write(b)
			case <-snapshotRefreshTimer.C:
				go func() {
					err := snapshotUpdater.Update()

					if err != nil {
						log.Printf("Error updating snapshots: %s", err)
					}
				}()
			}
		}
	}()
	<-quit
}
