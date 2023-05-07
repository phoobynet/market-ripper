package models

import (
	"time"
)

// Bar is the model for a bar which could be either a crypto bar or an equity bar
type Bar struct {
	Class      string    `json:"class"`
	Symbol     string    `json:"S"`
	Open       float64   `json:"o"`
	High       float64   `json:"h"`
	Low        float64   `json:"l"`
	Close      float64   `json:"c"`
	Volume     float64   `json:"v"`
	VWAP       float64   `json:"vw"`
	TradeCount float64   `json:"n"`
	Timestamp  time.Time `json:"t"`
}
