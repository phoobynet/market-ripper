package snapshots

import (
	"context"
	"fmt"
	. "github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/phoobynet/market-ripper/config"
	"github.com/phoobynet/market-ripper/database"
	"github.com/questdb/go-questdb-client"
	"github.com/samber/lo"
	"log"
)

func Load(ctx context.Context, configuration *config.Config, allSymbols []string) {
	log.Println("Loading snapshots...")

	db := database.NewDatabase(ctx, configuration)
	db.DropTable("sip_observer_snapshots")
	defer func(db *database.Database) {
		db.Close()
	}(db)

	sender, err := questdb.NewLineSender(ctx, questdb.WithAddress(fmt.Sprintf("%s:%s", configuration.DBHost, configuration.DBILPPort)))

	defer func(sender *questdb.LineSender) {
		if err != nil {
			log.Fatal(err)
		}
		_ = sender.Close()
	}(sender)

	if err != nil {
		log.Fatal(err)
	}

	client := NewClient(ClientOpts{})

	// You cannot request snapshots using splat, hence we need override the symbols value in this specific case
	var actualSymbols []string

	if len(configuration.Symbols) == 1 && configuration.Symbols[0] == "*" {
		actualSymbols = allSymbols
	} else {
		actualSymbols = configuration.Symbols
	}

	symbolChunks := lo.Chunk(actualSymbols, 500)

	mostRecentSnapshots := make(map[string]*Snapshot)

	for i, symbols := range symbolChunks {
		log.Printf("Loading snapshots...from Alpaca...chunk #%d of %d", i+1, len(symbolChunks))
		snapshotsChunk, err := client.GetSnapshots(symbols, GetSnapshotRequest{
			Feed: "sip",
		})

		if err != nil {
			log.Fatal(err)
		}

		for symbol, snapshot := range snapshotsChunk {
			mostRecentSnapshots[symbol] = snapshot
		}
	}

	log.Println("Loading snapshots...from Alpaca...DONE")

	count := 0

	for symbol, snapshot := range mostRecentSnapshots {
		if snapshot == nil {
			continue
		}

		dailyBar := snapshot.DailyBar

		if dailyBar == nil {
			continue
		}

		prevDailyBar := snapshot.PrevDailyBar

		if prevDailyBar == nil {
			continue
		}

		err = sender.Table("sip_observer_snapshots").
			Symbol("ticker", symbol).
			Float64Column("daily_bar_o", dailyBar.Open).
			Float64Column("daily_bar_h", dailyBar.High).
			Float64Column("daily_bar_l", dailyBar.Low).
			Float64Column("daily_bar_c", dailyBar.Close).
			Float64Column("daily_bar_v", float64(dailyBar.Volume)).
			Float64Column("daily_bar_n", float64(dailyBar.TradeCount)).
			Int64Column("daily_bar_t_micro", dailyBar.Timestamp.UnixMicro()).
			Float64Column("prev_daily_bar_o", prevDailyBar.Open).
			Float64Column("prev_daily_bar_h", prevDailyBar.High).
			Float64Column("prev_daily_bar_l", prevDailyBar.Low).
			Float64Column("prev_daily_bar_c", prevDailyBar.Close).
			Float64Column("prev_daily_bar_v", float64(prevDailyBar.Volume)).
			Float64Column("prev_daily_bar_n", float64(prevDailyBar.TradeCount)).
			Int64Column("prev_daily_bar_t_micro", prevDailyBar.Timestamp.UnixMicro()).
			AtNow(ctx)

		if err != nil {
			log.Fatal(err)
		}

		count++

		if count%1_000 == 0 {
			err = sender.Flush(ctx)

			if err != nil {
				log.Fatal(err)
			}
		}
	}

	err = sender.Flush(ctx)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Loading snapshots...COMPLETED %d loaded\n", count)
}
