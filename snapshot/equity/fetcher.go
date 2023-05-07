package equity

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/phoobynet/market-ripper/asset"
	"github.com/phoobynet/market-ripper/snapshot/models"
	"github.com/samber/lo"
)

const chunkSize = 50

// Fetcher is an equity snapshot fetcher
type Fetcher struct {
	assetRepository  *asset.Repository
	marketDataClient *marketdata.Client
}

// NewFetcher creates a new equity snapshot fetcher
func NewFetcher(assetRepository *asset.Repository, marketDataClient *marketdata.Client) *Fetcher {
	return &Fetcher{
		assetRepository:  assetRepository,
		marketDataClient: marketDataClient,
	}
}

// Fetch fetches equity snapshots for the given symbols and converts them to models.Snapshot
func (f *Fetcher) Fetch(symbols []string) (map[string]models.Snapshot, error) {
	snapshots := make(map[string]models.Snapshot)

	symbolChunks := lo.Chunk(symbols, chunkSize)

	for _, symbolsSubset := range symbolChunks {
		equitySnapshots, err := f.marketDataClient.GetSnapshots(symbolsSubset, marketdata.GetSnapshotRequest{})

		if err != nil {
			return nil, err
		}

		for symbol, equitySnapshot := range equitySnapshots {
			snapshots[symbol] = *Adapt(symbol, equitySnapshot)
		}
	}

	return snapshots, nil
}
