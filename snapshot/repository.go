package snapshot

import (
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) (*Repository, error) {
	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) Insert(snapshots []Snapshot) error {
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

func (r *Repository) update(snapshots map[string]Snapshot) error {
	// TODO: Is it possible to update all records in one query?
	var err error

	for _, s := range lo.Values(snapshots) {
		err = r.db.Updates(s).Error

		if err != nil {
			return err
		}
	}

	return nil
}
