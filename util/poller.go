package babelcoin

import (
	. "github.com/lox/babelcoin/core"
	"time"
)

// polls Exchange.MarketData periodically, writes data to a channel
func Poller(ex Exchange, p Pair, freq time.Duration, channel chan<- MarketData) error {
	ticker := time.NewTicker(freq)
	data, err := ex.MarketData(p)
	if err != nil {
		return err
	}

	channel <- data
	go func() {
		for _ = range ticker.C {
			data, _ = ex.MarketData(p)
			channel <- data
		}
	}()

	return nil
}
