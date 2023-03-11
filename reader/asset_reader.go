package reader

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"log"
)

type AssetReader struct {
}

func NewAssetReader() *AssetReader {
	return &AssetReader{}
}

func (t *AssetReader) ReadAllActive() []alpaca.Asset {
	assets, err := alpacaClient.GetAssets(alpaca.GetAssetsRequest{Status: "active"})

	if err != nil {
		log.Fatal(err)
	}

	return assets
}
