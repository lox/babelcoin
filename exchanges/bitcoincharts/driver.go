package bitcoincharts

import (
	//"github.com/davecgh/go-spew/spew"
	babel "../../core"
	util "../../util"
	//"errors"
	"time"
)

type Driver struct {
	config map[string]interface{}
}

func NewDriver(config map[string]interface{}) *Driver {
	return &Driver{config}
}

func (b *Driver) Ticker(symbol string) (chan babel.MarketData, chan bool, error) {
	api := NewMarketsApi(MarketUrl)

	duration, ok := b.config["poll_duration"].(time.Duration)
	if !ok {
		duration = time.Duration(5) * time.Second
	}

	channel, quit, err := util.Poller(duration, func() babel.MarketData{
		markets, err := api.Markets()
		if err != nil {
			panic(err)
		}
		return &MarketDataAdaptor{markets[symbol]}
	})

	if err != nil {
		return channel, quit, err
	}

	return channel, quit, nil
}

func (b *Driver) Symbols() ([]string, error) {
	api := NewMarketsApi(MarketUrl)
	markets, err := api.Markets()
	if err != nil {
		return []string{}, err
	}

	var symbols []string

	for _, market := range markets {
		symbols = append(symbols, market.Symbol)
	}

	return symbols, nil
}

type MarketDataAdaptor struct {
	data Market
}

func (d *MarketDataAdaptor) Ask() float64 {
	return d.data.Ask
}

func (d *MarketDataAdaptor) Bid() float64 {
	return d.data.Bid
}

func (d *MarketDataAdaptor) Last() float64 {
	return d.data.Close
}

func (d *MarketDataAdaptor) Volume() float64 {
	return d.data.Volume
}

func (d *MarketDataAdaptor) Updated() time.Time {
	return d.data.LatestTrade.Time
}


