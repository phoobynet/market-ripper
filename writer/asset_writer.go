package writer

import (
	"context"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/samber/lo"
	"log"
)

type AssetWriter struct {
}

func NewAssetWriter() *AssetWriter {
	return &AssetWriter{}
}

func (a *AssetWriter) Write(asset []alpaca.Asset) {
	assetChunks := lo.Chunk(asset, 1_000)

	ctx := context.TODO()

	for _, assets := range assetChunks {
		for _, asset := range assets {
			var class string

			if asset.Class == "us_equity" {
				class = "e"
			} else {
				class = "c"
			}

			err := lineSender.Table("assets").
				Symbol("ticker", asset.Symbol).
				StringColumn("class", class).
				StringColumn("name", asset.Name).
				StringColumn("exchange", asset.Exchange).
				AtNow(ctx)

			if err != nil {
				log.Fatal(err)
			}
		}

		err := lineSender.Flush(ctx)

		if err != nil {
			log.Fatal(err)
		}
	}
}
