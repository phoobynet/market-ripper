package bar

import "github.com/phoobynet/market-ripper/bar/models"

// Reader is a bar reader that wraps stream clients for crypto and equity bars
type Reader interface {
	Subscribe(bars chan models.Bar) error
	Unsubscribe() error
}
