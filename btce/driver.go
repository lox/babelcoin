package btce

import (
	//"github.com/davecgh/go-spew/spew"
	"../core"
	"time"
)

type BtceExchange struct {
	pair string
	config map[string]string
}

func NewExchange(pair string, config map[string]string) (*BtceExchange, error) {
	return &BtceExchange{pair, config}, nil
}

func (b *BtceExchange) MarketData() (babelcoin.MarketDataService, error) {
	return &BtceMarketDataService{
		NewTickerApi(TickerUrl, []string{b.pair}),
	}, nil
}

type BtceMarketDataService struct {
	ticker *BtceTickerApi
}

func (b *BtceMarketDataService) Fetch() (babelcoin.MarketData, error) {
	data, err := b.ticker.MarketData()
	if err != nil {
		return nil, err
	}
	return &BtceMarketData{data[0]}, nil
}

type BtceMarketData struct {
	data MarketData
}

func (d *BtceMarketData) Ask() float64 {
	return d.data.Sell
}

func (d *BtceMarketData) Bid() float64 {
	return d.data.Buy
}

func (d *BtceMarketData) Last() float64 {
	return d.data.Last
}

func (d *BtceMarketData) Updated() time.Time {
	return d.data.Updated.Time
}