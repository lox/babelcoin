package babelcoin

import (
	"time"
	. "github.com/lox/babelcoin/core"
)

// polls Exchange.MarketData periodically, writes data to a channel
func MarketDataPoller(ex Exchange, pair Pair, freq time.Duration, channel chan<- MarketData) error {
	ticker := time.NewTicker(freq)
	data, err := ex.MarketData(pair)
	if err != nil {
		return err
	}

	channel <- data
	go func() {
		for _ = range ticker.C {
			data, _ = ex.MarketData(pair)
			channel <- data
		}
	}()

	return nil
}

// polls Exchange.History periodically, trades to channel. no de-duping occurs
func HistoryPoller(ex Exchange, pairs []Pair, freq time.Duration, channel chan<- Trade) error {
	ticker := time.NewTicker(freq)

	go func() {
		after := time.Now().AddDate(0, 0, -3) // 3 days ago
		limit := 2000

		for _ = range ticker.C {
			var trades = make(chan Trade)

			go func() {
				if err := ex.TradeHistory(pairs, after, limit, trades); err != nil {
					close(trades)
				}
			}()

			for trade := range trades {
				channel <- trade
			}

			after = time.Now().Add(-(time.Minute * 15))
			limit = 100
		}
	}()

	return nil
}
