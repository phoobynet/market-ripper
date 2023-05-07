package snapshot

import "github.com/phoobynet/market-ripper/snapshot/models"

type Fetcher interface {
	Fetch(symbols []string) (map[string]models.Snapshot, error)
}
