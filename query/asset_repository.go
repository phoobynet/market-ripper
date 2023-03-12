package query

import (
	"context"
	"fmt"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"log"
	"time"
)

type AssetRepository struct{}

func NewAssetRepository() *AssetRepository {
	_, err := connection.Exec(
		context.TODO(),
		`
			CREATE TABLE IF NOT EXISTS assets(
				ticker symbol,
				name string,
				exchange string,
				class string,
				timestamp timestamp
			)`,
	)

	if err != nil {
		log.Fatal(err)
	}

	return &AssetRepository{}
}

func (t *AssetRepository) Count() int {
	var count int
	err := connection.QueryRow(context.TODO(), "SELECT COUNT(*) FROM assets").Scan(&count)

	if err != nil {
		log.Fatal(err)
	}

	return count
}

func (t *AssetRepository) LastUpdated() time.Time {
	var lastUpdated time.Time
	err := connection.QueryRow(context.TODO(), "SELECT MAX(timestamp) FROM assets").Scan(&lastUpdated)

	if err != nil {
		log.Fatal(err)
	}

	return lastUpdated
}

func (t *AssetRepository) IsStale(maxAge time.Duration) bool {
	lastUpdated := t.LastUpdated()

	now := time.Now()

	staleAt := now.Add(maxAge)

	if staleAt.After(now) {
		log.Fatal("Incorrect use of IsStale maxAge - duration should be negative")
	}

	if lastUpdated.Before(staleAt) {
		return true
	}

	return false
}

func (t *AssetRepository) Truncate() {
	_, err := connection.Exec(context.TODO(), "TRUNCATE TABLE assets")
	if err != nil {
		log.Fatal(err)
	}
}

func (t *AssetRepository) GetSymbolByClass(class alpaca.AssetClass) []string {
	var symbols []string
	cursor, err := connection.Query(context.TODO(), fmt.Sprintf("select ticker from assets where class = '%s' order by ticker", class))

	if err != nil {
		log.Fatal(err)
	}

	var symbol string

	for cursor.Next() {
		err := cursor.Scan(&symbol)
		if err != nil {
			log.Fatal(err)
		}

		symbols = append(symbols, symbol)
	}

	return symbols
}
