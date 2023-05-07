package bar

import "github.com/phoobynet/market-ripper/bar/models"

type Reader interface {
	Subscribe(bars chan models.Bar) error
	Unsubscribe() error
}
