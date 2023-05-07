package asset

import "github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"

type Fetch struct {
	client *alpaca.Client
}

func NewFetch(client *alpaca.Client) *Fetch {
	return &Fetch{client: client}
}

func (f *Fetch) Fetch() ([]alpaca.Asset, error) {
	assets, err := f.client.GetAssets(alpaca.GetAssetsRequest{
		Status: "active",
	})

	if err != nil {
		return nil, err
	}

	return assets, nil
}
