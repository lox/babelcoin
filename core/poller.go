package babelcoin

import (
	"time"
)

type MarketDataServicePoller struct {
	service MarketDataService
	frequency time.Duration
}

func NewMarketDataServicePoller(service MarketDataService, d time.Duration) MarketDataFeed {
	return &MarketDataServicePoller{service, d} 
}

func (p *MarketDataServicePoller) Channel() chan MarketData {
	channel := make(chan MarketData)
	ticker := time.NewTicker(p.frequency)

    go func() {
		data, _ := p.service.Fetch()
		channel <- data

        for _ = range ticker.C {
			data, _ := p.service.Fetch()
			channel <- data
        }
    }()

	return channel
}
