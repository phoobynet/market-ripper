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
	"gorm.io/gorm"
)

func Prepare(
	configuration *config.Config,
	db *gorm.DB,
	assetRepository *asset.Repository,
	marketDataClient *marketdata.Client) (*Repository, *Updater, error) {
	err := db.AutoMigrate(&models.Snapshot{})

	if err != nil {
		return nil, nil, err
	}

	repository, err := newRepository(db)

	if err != nil {
		return nil, nil, err
	}

	err = repository.truncate()

	if err != nil {
		return nil, nil, err
	}

	var snapshots map[string]models.Snapshot

	symbolsChunks := lo.Chunk(configuration.Symbols, chunkSize)

	var fetcher Fetcher

	if configuration.Class == alpaca.USEquity {
		fetcher = equity.NewFetcher(assetRepository, marketDataClient)
	} else {
		fetcher = crypto.NewFetcher(assetRepository, marketDataClient)
	}

	for _, symbolsChunk := range symbolsChunks {
		snapshotsInChunk, err := fetcher.Fetch(symbolsChunk)

		if err != nil {
			return nil, nil, err
		}

		snapshots = snapshotsInChunk
	}

	err = repository.Insert(lo.Values(snapshots))

	if err != nil {
		return nil, nil, err
	}

	updater, err := NewUpdater(configuration, assetRepository, marketDataClient, repository)

	if err != nil {
		return nil, nil, err
	}

	return repository, updater, nil
}
