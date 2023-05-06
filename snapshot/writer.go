package snapshot

import (
	"context"
	"fmt"
	"github.com/phoobynet/market-ripper/config"
	"github.com/questdb/go-questdb-client"
	"log"
	"time"
)

type Writer struct {
	configuration *config.Config
	lineSender    *questdb.LineSender
	tableName     string
}

func NewWriter(configuration *config.Config) (*Writer, error) {
	sender, err := questdb.NewLineSender(
		context.TODO(),
		configuration.GetIngressAddress(),
	)

	if err != nil {
		return nil, err
	}

	return &Writer{
		lineSender:    sender,
		configuration: configuration,
		tableName: fmt.Sprintf(
			"%s_snapshots",
			configuration.Class,
		),
	}, nil
}

func (w *Writer) Write(snapshots map[string]*Snapshot) {
	ctx := context.Background()

	count := 0

	for symbol, s := range snapshots {
		if s == nil {
			continue
		}

		err := w.lineSender.Table(w.tableName).Symbol(
			"ticker",
			symbol,
		).
			Float64Column(
				"daily_bar_o",
				s.DailyOpen,
			).
			Float64Column(
				"daily_bar_h",
				s.DailyHigh,
			).
			Float64Column(
				"daily_bar_l",
				s.DailyLow,
			).
			Float64Column(
				"daily_bar_c",
				s.DailyClose,
			).
			Float64Column(
				"daily_bar_v",
				s.DailyVolume,
			).
			Int64Column(
				"daily_bar_t",
				s.DailyTimestamp.UnixMicro(),
			).
			Float64Column(
				"prev_daily_bar_o",
				s.PreviousOpen,
			).
			Float64Column(
				"prev_daily_bar_h",
				s.PreviousHigh,
			).
			Float64Column(
				"prev_daily_bar_l",
				s.PreviousLow,
			).
			Float64Column(
				"prev_daily_bar_c",
				s.PreviousClose,
			).
			Float64Column(
				"prev_daily_bar_v",
				s.PreviousVolume,
			).
			Int64Column(
				"prev_daily_bar_t",
				s.PreviousTimestamp.UnixMicro(),
			).
			TimestampColumn(
				"timestamp",
				time.Now().UnixMicro(),
			).
			AtNow(ctx)

		if err != nil {
			log.Fatal(err)
		}

		count++

		if count%1_000 == 0 {
			err = w.lineSender.Flush(ctx)

			if err != nil {
				log.Fatal(err)
			}
		}
	}

	err := w.lineSender.Flush(ctx)

	if err != nil {
		log.Fatal(err)
	}
}

func (w *Writer) Close() {
	_ = w.lineSender.Close()
}
