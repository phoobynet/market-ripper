package snapshot

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/phoobynet/market-ripper/asset"
	"github.com/phoobynet/market-ripper/config"
	"github.com/phoobynet/market-ripper/snapshot/crypto"
	"github.com/phoobynet/market-ripper/snapshot/equity"
	"github.com/phoobynet/market-ripper/snapshot/models"
	"github.com/samber/lo"
)

// Updater updates the snapshot table with the latest market data
type Updater struct {
	config             *config.Config
	cryptoFetcher      *crypto.Fetcher
	equityFetcher      *equity.Fetcher
	snapshotRepository *Repository
}

const chunkSize = 100

func NewUpdater(
	config *config.Config,
	assetRepository *asset.Repository,
	marketDataClient *marketdata.Client,
	snapshotRepository *Repository) (*Updater, error) {
	return &Updater{
		config:             config,
		cryptoFetcher:      crypto.NewFetcher(assetRepository, marketDataClient),
		equityFetcher:      equity.NewFetcher(assetRepository, marketDataClient),
		snapshotRepository: snapshotRepository,
	}, nil
}

func (u *Updater) Update() error {
	var snapshots map[string]models.Snapshot
	symbolsChunks := lo.Chunk(u.config.Symbols, chunkSize)

	var fetcher Fetcher

	if u.config.Class == alpaca.USEquity {
		fetcher = u.equityFetcher
	} else {
		fetcher = u.cryptoFetcher
	}

	for _, symbolsChunk := range symbolsChunks {
		snapshotsInChunk, err := fetcher.Fetch(symbolsChunk)

		if err != nil {
			return err
		}

		snapshots = snapshotsInChunk
	}

	return u.snapshotRepository.upsert(snapshots)
}
