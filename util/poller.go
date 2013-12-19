package babelcoin

import (
	. "../core"
	"time"
)

func Poller(freq time.Duration, fetch func() MarketData) (chan MarketData, chan bool, error) {
	channel := make(chan MarketData, 10)
	quit := make(chan bool)
	ticker := time.NewTicker(freq)
	channel <- fetch()

	go func() {
		for _ = range ticker.C {
			select {
			case <-quit:
				break
			default:
				channel <- fetch()
			}
		}
	}()

	return channel, quit, nil
}
