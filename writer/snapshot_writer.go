package writer

import (
	"context"
	"fmt"
	"github.com/phoobynet/market-ripper/config"
	"github.com/phoobynet/market-ripper/types"
	"github.com/questdb/go-questdb-client"
	"log"
)

type SnapshotWriter struct {
	configuration *config.Config
	lineSender    *questdb.LineSender
	tableName     string
}

func NewSnapshotWriter(configuration *config.Config) *SnapshotWriter {
	sender, err := questdb.NewLineSender(context.TODO(), configuration.GetIngressAddress())

	if err != nil {
		log.Fatal(err)
	}

	return &SnapshotWriter{
		lineSender:    sender,
		configuration: configuration,
		tableName:     fmt.Sprintf("%s_snapshots", configuration.Class),
	}
}

func (s *SnapshotWriter) Write(snapshots map[string]*types.Snapshot) {
	ctx := context.Background()

	count := 0

	for symbol, snapshot := range snapshots {
		if snapshot == nil {
			continue
		}

		err := s.lineSender.Table(s.tableName).Symbol("ticker", symbol).
			Float64Column("daily_bar_o", snapshot.DailyOpen).
			Float64Column("daily_bar_h", snapshot.DailyHigh).
			Float64Column("daily_bar_l", snapshot.DailyLow).
			Float64Column("daily_bar_c", snapshot.DailyClose).
			Float64Column("daily_bar_v", snapshot.DailyVolume).
			Int64Column("daily_bar_t", snapshot.DailyTimestamp.UnixMicro()).
			Float64Column("prev_daily_bar_o", snapshot.PreviousOpen).
			Float64Column("prev_daily_bar_h", snapshot.PreviousHigh).
			Float64Column("prev_daily_bar_l", snapshot.PreviousLow).
			Float64Column("prev_daily_bar_c", snapshot.PreviousClose).
			Float64Column("prev_daily_bar_v", snapshot.PreviousVolume).
			Int64Column("prev_daily_bar_t", snapshot.PreviousTimestamp.UnixMicro()).
			AtNow(ctx)

		if err != nil {
			log.Fatal(err)
		}

		count++

		if count%1_000 == 0 {
			err = s.lineSender.Flush(ctx)

			if err != nil {
				log.Fatal(err)
			}
		}
	}

	err := s.lineSender.Flush(ctx)

	if err != nil {
		log.Fatal(err)
	}
}

func (s *SnapshotWriter) Close() {
	_ = s.lineSender.Close()
}
