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
	"github.com/phoobynet/market-ripper/bar/models"
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

func fatalOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

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

	fatalOnErr(err)

	log.Printf(
		"%s",
		configuration,
	)

	alpacaClient := alpaca.NewClient(alpaca.ClientOpts{})
	marketDataClient := marketdata.NewClient(marketdata.ClientOpts{})

	// Connect to equities streaming client
	stocksClientContext, stocksClientCancel := context.WithCancel(context.Background())
	defer stocksClientCancel()
	stocksClient := stream.NewStocksClient(marketdata.SIP)
	err = stocksClient.Connect(stocksClientContext)
	fatalOnErr(err)
	log.Println("Connected to stocks client")

	// Connect to crypto streaming client
	cryptoClientContext, cryptoClientCancel := context.WithCancel(context.Background())
	defer cryptoClientCancel()
	cryptoClient := stream.NewCryptoClient("us")
	err = cryptoClient.Connect(cryptoClientContext)
	fatalOnErr(err)
	log.Println("Connected to crypto client")

	// Connect to QuestDB database using the Postgres protocol
	pgConnection, err := database.Connect(configuration)
	fatalOnErr(err)

	defer func(pgConnection *gorm.DB) {
		db, err := pgConnection.DB()

		if err == nil {
			_ = db.Close()
		}
	}(pgConnection)

	assetRepository, err := asset.Prepare(pgConnection, alpacaClient)
	fatalOnErr(err)

	_, snapshotUpdater, err := snapshot.Prepare(configuration, pgConnection, assetRepository, marketDataClient)
	fatalOnErr(err)

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
		chan models.Bar,
		20_000,
	)

	var barReader bar.Reader
	var tradeReader trade.Reader

	if configuration.Class == "us_equity" {
		barReader, err = barEquity.NewReader(configuration, stocksClient)
		fatalOnErr(err)

		tradeReader, err = tradeEquity.NewReader(configuration, stocksClient)
		fatalOnErr(err)
	} else {
		barReader, err = barCrypto.NewReader(configuration, cryptoClient)
		fatalOnErr(err)

		tradeReader, err = tradeCrypto.NewReader(configuration, cryptoClient)
		fatalOnErr(err)
	}

	err = barReader.Subscribe(barsChannel)
	fatalOnErr(err)

	err = tradeReader.Subscribe(tradesChannel)
	fatalOnErr(err)

	defer func(barReader bar.Reader) {
		_ = barReader.Unsubscribe()
	}(barReader)

	defer func(tradeReader trade.Reader) {
		_ = tradeReader.Unsubscribe()
	}(tradeReader)

	barWriter := bar.NewWriter(configuration)
	defer barWriter.Close()

	tradeWriter := trade.NewWriter(configuration)
	defer tradeWriter.Close()

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
