package babelcoin

import (
	"time"
)

type MarketData interface {
	Ask() float64
	Bid() float64
	Last() float64
	Volume() float64
	Updated() time.Time
}

type Exchange interface {
	Ticker(symbol string) (chan MarketData, chan bool, error)
	Symbols() ([]string, error)
}