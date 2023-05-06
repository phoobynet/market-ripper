package asset

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
)

type Reader struct {
	client *alpaca.Client
}

func NewAssetReader(client *alpaca.Client) (*Reader, error) {
	reader := &Reader{
		client: client,
	}

	return reader, nil
}

// GetActive - Retrieves all active assets from Alpaca (us_equity and crypto)
func (r *Reader) GetActive() ([]alpaca.Asset, error) {
	assets, err := r.client.GetAssets(alpaca.GetAssetsRequest{Status: "active"})

	return assets, err
}
