package trade

type Reader interface {
	Subscribe(bars chan Trade) error
	Unsubscribe() error
}
