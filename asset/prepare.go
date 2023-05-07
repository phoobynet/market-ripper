package asset

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"gorm.io/gorm"
)

// Prepare prepares the asset repository.  If the repository is empty, it will be populated
// TODO: Need to have some sort of stale check
func Prepare(db *gorm.DB, client *alpaca.Client) (*Repository, error) {
	err := db.AutoMigrate(&alpaca.Asset{})

	if err != nil {
		return nil, err
	}

	repository, err := NewRepository(db)

	if err != nil {
		return nil, err
	}

	if isEmpty, err := repository.IsEmpty(); err != nil {
		return nil, err
	} else if isEmpty {
		assets, err := NewFetcher(client).Fetch()

		if err != nil {
			return nil, err
		}

		err = repository.Insert(assets)

		if err != nil {
			return nil, err
		}
	}

	return repository, nil
}
