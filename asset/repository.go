package asset

import (
	"database/sql"
	"fmt"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"log"
	"time"
)

type Repository struct {
	connection *sql.DB
}

func NewRepository(connection *sql.DB) (*Repository, error) {
	_, err := connection.Exec(
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
		return nil, err
	}

	return &Repository{
		connection: connection,
	}, nil
}

func (r *Repository) Count() int {
	var count int
	err := r.connection.QueryRow("SELECT COUNT(*) FROM assets").Scan(&count)

	if err != nil {
		log.Fatal(err)
	}

	return count
}

func (r *Repository) LastUpdated() time.Time {
	var lastUpdated time.Time
	err := r.connection.QueryRow("SELECT MAX(timestamp) FROM assets").Scan(&lastUpdated)

	if err != nil {
		log.Fatal(err)
	}

	return lastUpdated
}

func (r *Repository) IsStale(maxAge time.Duration) bool {
	lastUpdated := r.LastUpdated()

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

func (r *Repository) Truncate() {
	_, err := r.connection.Exec("TRUNCATE TABLE assets")
	if err != nil {
		log.Fatal(err)
	}
}

func (r *Repository) GetSymbolByClass(class alpaca.AssetClass) ([]string, error) {
	var symbols []string
	cursor, err := r.connection.Query(fmt.Sprintf("select ticker from assets where class = '%s' order by ticker", class))

	if err != nil {
		log.Fatal(err)
	}

	var symbol string

	for cursor.Next() {
		err := cursor.Scan(&symbol)
		if err != nil {
			return nil, err
		}

		symbols = append(symbols, symbol)
	}

	return symbols, nil
}
