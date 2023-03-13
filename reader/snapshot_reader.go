package reader

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/phoobynet/market-ripper/config"
	"github.com/phoobynet/market-ripper/query"
	"github.com/phoobynet/market-ripper/types"
	"github.com/samber/lo"
	"log"
)

type SnapshotReader struct {
	configuration   *config.Config
	assetRepository *query.AssetRepository
}

func NewSnapshotReader(
	configuration *config.Config,
	assetRepository *query.AssetRepository,
) *SnapshotReader {
	return &SnapshotReader{
		configuration:   configuration,
		assetRepository: assetRepository,
	}
}

func (t *SnapshotReader) Read() map[string]*types.Snapshot {
	var actualSymbols []string

	if len(t.configuration.Symbols) == 1 && t.configuration.Symbols[0] == "*" {
		actualSymbols = t.assetRepository.GetSymbolByClass(t.configuration.Class)
	} else {
		actualSymbols = t.configuration.Symbols
	}

	symbolChunks := lo.Chunk(
		actualSymbols,
		500,
	)

	// TODO: This code sucks! Refactor it!
	mostRecentSnapshots := make(map[string]*types.Snapshot)

	for i, symbols := range symbolChunks {
		log.Printf(
			"Loading snapshots...from Alpaca...chunk #%d of %d",
			i+1,
			len(symbolChunks),
		)

		if t.configuration.Class == alpaca.Crypto {
			snapshotsChunk, err := marketDataClient.GetCryptoSnapshots(
				symbols,
				marketdata.GetCryptoSnapshotRequest{},
			)

			if err != nil {
				log.Fatal(err)
			}

			for symbol, snapshot := range snapshotsChunk {
				mostRecentSnapshots[symbol] = types.FromSnapshot(
					symbol,
					snapshot,
				)
			}
		} else {
			snapshotsChunk, err := marketDataClient.GetSnapshots(
				symbols,
				marketdata.GetSnapshotRequest{
					Feed: "sip",
				},
			)

			if err != nil {
				log.Fatal(err)
			}

			for symbol, snapshot := range snapshotsChunk {
				mostRecentSnapshots[symbol] = types.FromSnapshot(
					symbol,
					snapshot,
				)
			}
		}

	}

	return mostRecentSnapshots
}
