package query

import (
	"context"
	"fmt"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/phoobynet/market-ripper/config"
	"log"
)

type SnapshotRepository struct {
	configuration *config.Config
	tableName     string
}

func NewSnapshotRepository(configuration *config.Config) *SnapshotRepository {
	var tableName string

	if configuration.Class == alpaca.Crypto {
		tableName = "crypto_snapshots"
	} else {
		tableName = "equity_snapshots"
	}

	_, err := connection.Exec(
		context.TODO(),
		fmt.Sprintf(`
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
				prev_daily_bar_t long
				timestamp timestamp
			)`,
			tableName),
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
		fmt.Sprintf("TRUNCATE TABLE %s", s.tableName),
	)

	if err != nil {
		log.Fatal(err)
	}
}
