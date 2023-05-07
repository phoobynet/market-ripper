package snapshot

import "github.com/phoobynet/market-ripper/snapshot/models"

// Fetcher is an interface for fetching snapshots
type Fetcher interface {
	Fetch(symbols []string) (map[string]models.Snapshot, error)
}
