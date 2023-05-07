package trade

// Reader is a trade reader interface for either crypto or equity trades
type Reader interface {
	Subscribe(bars chan Trade) error
	Unsubscribe() error
}
