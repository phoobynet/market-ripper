package bar

type Reader interface {
	Subscribe(bars chan Bar) error
	Unsubscribe() error
}
