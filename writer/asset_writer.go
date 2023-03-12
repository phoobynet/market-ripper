package writer

import (
	"context"
	"fmt"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/phoobynet/market-ripper/config"
	"github.com/questdb/go-questdb-client"
	"github.com/samber/lo"
	"log"
	"time"
)

type AssetWriter struct {
	lineSender *questdb.LineSender
}

func NewAssetWriter(configuration *config.Config) *AssetWriter {
	sender, err := questdb.NewLineSender(context.TODO(), questdb.WithAddress(fmt.Sprintf("%s:%s", configuration.DBHost, configuration.DBILPPort)))

	if err != nil {
		log.Fatal(err)
	}

	return &AssetWriter{
		lineSender: sender,
	}
}

func (a *AssetWriter) Write(asset []alpaca.Asset) {
	assetChunks := lo.Chunk(asset, 1_000)

	ctx := context.TODO()

	for _, assets := range assetChunks {
		for _, asset := range assets {
			var class string

			if asset.Class == alpaca.USEquity {
				class = "us_equity"
			} else {
				class = "crypto"
			}

			err := a.lineSender.Table("assets").
				Symbol("ticker", asset.Symbol).
				StringColumn("class", class).
				StringColumn("name", asset.Name).
				StringColumn("exchange", asset.Exchange).
				TimestampColumn("timestamp", time.Now().UnixMicro()).
				AtNow(ctx)

			if err != nil {
				log.Fatal(err)
			}
		}

		err := a.lineSender.Flush(ctx)

		if err != nil {
			log.Fatal(err)
		}
	}
}

func (a *AssetWriter) Close() {
	a.lineSender.Close()
}
