package snapshot

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/phoobynet/market-ripper/asset"
	"github.com/phoobynet/market-ripper/config"
	"gorm.io/gorm"
)

func Prepare(config *config.Config, db *gorm.DB, assetRepository *asset.Repository) (*Repository, error) {
	repository, err := NewRepository(db)

	if err != nil {
		return nil, err
	}

	err = repository.truncate()

	if err != nil {
		return nil, err
	}

	if config.Class == alpaca.USEquity {
		symbols, err := assetRepository.GetSymbolByClass(alpaca.USEquity)

		if err != nil {
			return nil, err
		}

		err = repository.Insert(symbols)

		if err != nil {
			return nil, err
		}
	}
}
