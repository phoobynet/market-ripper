package query

import (
	"context"
	"fmt"
	"github.com/phoobynet/market-ripper/config"
	"github.com/phoobynet/market-ripper/types"
	"log"
	"time"
)

type SnapshotRepository struct {
	configuration *config.Config
	tableName     string
}

func NewSnapshotRepository(configuration *config.Config) *SnapshotRepository {
	tableName := fmt.Sprintf(
		"%s_snapshots",
		configuration.Class,
	)

	_, err := connection.Exec(
		context.TODO(),
		fmt.Sprintf(
			`
			CREATE TABLE IF NOT EXISTS %s(
				ticker symbol,
				daily_bar_o float,
				daily_bar_h float,
				daily_bar_l float,
				daily_bar_c float,
				daily_bar_v float,
				daily_bar_t long,
				prev_daily_bar_o float,
				prev_daily_bar_h float,
				prev_daily_bar_l float,
				prev_daily_bar_c float,
				prev_daily_bar_v float,
				prev_daily_bar_t long,
				timestamp timestamp
			)`,
			tableName,
		),
	)

	if err != nil {
		log.Fatal(err)
	}

	return &SnapshotRepository{
		configuration: configuration,
		tableName:     tableName,
	}
}

func (s *SnapshotRepository) Truncate() {
	_, err := connection.Exec(
		context.TODO(),
		fmt.Sprintf(
			"TRUNCATE TABLE %s",
			s.tableName,
		),
	)

	if err != nil {
		log.Fatal(err)
	}
}

func (s *SnapshotRepository) Update(snapshots map[string]*types.Snapshot) {
	updateSQL := fmt.Sprintf(
		`UPDATE %s 
				SET
				daily_bar_o = $1, 
				daily_bar_h = $2, 
				daily_bar_l = $3, 
				daily_bar_c = $4, 
				daily_bar_v = $5, 
				daily_bar_t = $6,
				prev_daily_bar_o = $7, 
				prev_daily_bar_h = $8, 
				prev_daily_bar_l = $9, 
				prev_daily_bar_c = $10, 
				prev_daily_bar_v = $11, 
				prev_daily_bar_t = $12,
				timestamp = $13
				WHERE ticker = $14`,
		s.tableName,
	)

	for symbol, snapshot := range snapshots {
		_, err := connection.Exec(
			context.TODO(),
			updateSQL,
			snapshot.DailyOpen,
			snapshot.DailyHigh,
			snapshot.DailyLow,
			snapshot.DailyClose,
			snapshot.DailyVolume,
			snapshot.DailyTimestamp.UnixMicro(),
			snapshot.PreviousOpen,
			snapshot.PreviousHigh,
			snapshot.PreviousLow,
			snapshot.PreviousClose,
			snapshot.PreviousVolume,
			snapshot.PreviousTimestamp.UnixMicro(),
			time.Now(),
			symbol,
		)

		if err != nil {
			log.Fatal(err)
		}
	}
}
