package babelcoin

import (
	"time"
)

const (
	MarketPrice = -1
)

type MarketData interface {
	Ask() float64
	Bid() float64
	Last() float64
	Volume() float64
	Updated() time.Time
}

type Trade interface {
	Amount() float64
	Rate() float64
}

type Order interface {
	Fee() (float64, error)
	Execute() (chan Trade, error)
	Cancel() error
}

type Exchange interface {
	Ticker(symbol string) (chan MarketData, chan bool, error)
	Symbols() ([]string, error)
	Buy(symbol string, amount float64, rate float64) Order
	Sell(symbol string, amount float64, rate float64) Order
}
