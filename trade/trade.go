package trade

import (
	"time"
)

type Trade struct {
	Class     string    `json:"class"`
	Symbol    string    `json:"S"`
	Price     float64   `json:"p"`
	Size      float64   `json:"s"`
	TakerSide string    `json:"takerSide"`
	Exchange  string    `json:"e"`
	Timestamp time.Time `json:"t"`
}
