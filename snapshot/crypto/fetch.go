package crypto

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/phoobynet/market-ripper/asset"
	"github.com/phoobynet/market-ripper/snapshot"
	"github.com/samber/lo"
)

const chunkSize = 50

type Fetch struct {
	assetRepository  *asset.Repository
	marketDataClient *marketdata.Client
}

func NewFetch(assetRepository *asset.Repository, marketDataClient *marketdata.Client) *Fetch {
	return &Fetch{
		assetRepository:  assetRepository,
		marketDataClient: marketDataClient,
	}
}

func (f *Fetch) Fetch(symbols []string) (map[string]snapshot.Snapshot, error) {
	snapshots := make(map[string]snapshot.Snapshot)

	symbolChunks := lo.Chunk(symbols, chunkSize)

	for _, symbolSubset := range symbolChunks {
		cryptoSnapshots, err := f.marketDataClient.GetCryptoSnapshots(symbolSubset, marketdata.GetCryptoSnapshotRequest{})

		if err != nil {
			return nil, err
		}

		for symbol, cryptoSnapshot := range cryptoSnapshots {
			s, err := Adapt(symbol, cryptoSnapshot)

			if err != nil {
				return nil, err
			}

			snapshots[symbol] = *s
		}
	}

	return snapshots, nil
}