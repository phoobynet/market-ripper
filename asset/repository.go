package asset

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

// Repository is a repository for assets.
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new Repository instance.
func NewRepository(db *gorm.DB) (*Repository, error) {
	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) count() (int64, error) {
	var count int64
	err := r.db.Model(&alpaca.Asset{}).Count(&count).Error

	return count, err
}

// IsEmpty checks if the repository is empty.
func (r *Repository) IsEmpty() (bool, error) {
	count, err := r.count()

	if err != nil {
		return false, err
	}

	return count == 0, nil
}

// GetSymbolByClass gets symbols by class, either "us_equity" or "crypto".
func (r *Repository) GetSymbolByClass(class alpaca.AssetClass) ([]string, error) {
	var symbols []string

	err := r.db.
		Model(&alpaca.Asset{}).
		Where("class = ?", class).
		Order("symbol").
		Pluck("symbol", &symbols).
		Error

	if err != nil {
		return nil, err
	}

	return symbols, nil
}

// Insert inserts assets into the repository in bulk.
func (r *Repository) Insert(assets []alpaca.Asset) error {
	assetChunks := lo.Chunk(assets, 100)

	var err error

	for _, assetChunk := range assetChunks {
		err = r.db.Create(assetChunk).Error

		if err != nil {
			break
		}
	}

	return err
}
