package babelcoin

import (
	"time"
	. "github.com/lox/babelcoin/core"
)

// polls Exchange.MarketData periodically, writes data to a channel
func MarketDataPoller(ex Exchange, p Pair, freq time.Duration, channel chan<- MarketData) error {
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

// polls Exchange.History periodically, writes new trades to channel
func HistoryPoller(ex Exchange, p Pair, freq time.Duration, channel chan<- Trade) error {
	ticker := time.NewTicker(freq)

	go func() {
		set := map[string]time.Time{}

		for _ = range ticker.C {
			trades := make(chan Trade)

			go func() {
				// get history for a day ago
				if err := ex.History(p, time.Now().AddDate(0, 0, -1), trades); err != nil {
					close(trades)
				}
			}()

			for trade := range trades {
				if _, exists := set[trade.Id]; !exists {
					channel <- trade
					set[trade.Id] = trade.Timestamp
				}
			}

			// purge old set entries
			for tid, ts := range set {
				if ts.Before(time.Now().Add(-(time.Minute * 3600))) {
					delete(set, tid)
				}
			}
		}
	}()

	return nil
}
