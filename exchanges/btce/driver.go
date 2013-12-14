package btce

import (
	//"github.com/davecgh/go-spew/spew"
	babel "../../core"
	util "../../util"
	"time"
)

type Driver struct {
	config map[string]interface{}
}

func NewDriver(config map[string]interface{}) (*Driver, error) {
	return &Driver{config}, nil
}

func (b *Driver) Symbols() ([]string, error) {
	api := NewInfoApi(InfoUrl)
	pairs, error := api.Pairs()
	if error != nil {
		return []string{}, error
	}

	var symbols []string
	for symbol := range pairs {
		symbols = append(symbols, symbol)
	}

	return symbols, nil
}

func (b *Driver) Ticker(symbol string) (chan babel.MarketData, chan bool, error) {
	api := NewTickerApi(TickerUrl, []string{symbol})

	duration, ok := b.config["poll_duration"].(time.Duration)
	if !ok {
		duration = time.Duration(5) * time.Second
	}

	channel, quit, err := util.Poller(duration, func() babel.MarketData{
		data, err := api.MarketData()
		if err != nil {
			panic(err)
		}
		return &MarketDataAdaptor{data[0]}
	})

	if err != nil {
		return channel, quit, err
	}

	return channel, quit, nil
}

type MarketDataAdaptor struct {
	data MarketData
}

func (d *MarketDataAdaptor) Ask() float64 {
	return d.data.Sell
}

func (d *MarketDataAdaptor) Bid() float64 {
	return d.data.Buy
}

func (d *MarketDataAdaptor) Last() float64 {
	return d.data.Last
}

func (d *MarketDataAdaptor) Volume() float64 {
	return d.data.Volume
}

func (d *MarketDataAdaptor) Updated() time.Time {
	return d.data.Updated.Time
}
