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
	MarketData(symbol string) (MarketDataService, error)
	Symbols() ([]string, error)
}

type MarketDataService interface {
	Fetch() (MarketData, error)
	Feed() (MarketDataFeed, error)
}

type MarketDataFeed interface {
	Channel() chan MarketData
}