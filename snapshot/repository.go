package snapshot

import (
	"github.com/phoobynet/market-ripper/snapshot/models"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func newRepository(db *gorm.DB) (*Repository, error) {
	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) Insert(snapshots []models.Snapshot) error {
	snapshotChunks := lo.Chunk(snapshots, 100)

	var err error

	for _, snapshotChunk := range snapshotChunks {
		err = r.db.Create(snapshotChunk).Error

		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) truncate() error {
	return r.db.Exec("TRUNCATE TABLE snapshots").Error
}

func (r *Repository) upsert(snapshots map[string]models.Snapshot) error {
	var err error

	for _, s := range lo.Values(snapshots) {
		err = r.db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(s).Error

		if err != nil {
			return err
		}
	}

	return nil
}
