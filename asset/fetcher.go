package asset

import "github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"

// Fetcher is a service that fetches assets.
type Fetcher struct {
	client *alpaca.Client
}

// NewFetcher creates a new Fetcher instance.
func NewFetcher(client *alpaca.Client) *Fetcher {
	return &Fetcher{client: client}
}

// Fetch fetches "active" assets.
func (f *Fetcher) Fetch() ([]alpaca.Asset, error) {
	assets, err := f.client.GetAssets(alpaca.GetAssetsRequest{
		Status: "active",
	})

	if err != nil {
		return nil, err
	}

	return assets, nil
}
