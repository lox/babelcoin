package btce

import (
	"../../core"
	"time"
)

type BtceExchange struct {
	config map[string]interface{}
}

func NewExchange(config map[string]interface{}) (*BtceExchange, error) {
	return &BtceExchange{config: nil}, nil
}

func (b *BtceExchange) MarketData(pair string) (babelcoin.MarketDataService, error) {
	return &BtceMarketDataService{
		NewTickerApi(TickerUrl, []string{pair}), b,
	}, nil
}

func (b *BtceExchange) Symbols() ([]string, error) {
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

type BtceMarketDataService struct {
	ticker   *BtceTickerApi
	exchange *BtceExchange
}

func (b *BtceMarketDataService) Fetch() (babelcoin.MarketData, error) {
	data, err := b.ticker.MarketData()
	if err != nil {
		return nil, err
	}
	return &BtceMarketData{data[0]}, nil
}

func (b *BtceMarketDataService) Feed() (babelcoin.MarketDataFeed, error) {
	duration, ok := b.exchange.config["poll_duration"]
	if !ok {
		duration = time.Duration(10) * time.Second
	}

	return babelcoin.NewMarketDataServicePoller(b, duration.(time.Duration)), nil
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

func (d *BtceMarketData) Volume() float64 {
	return d.data.Volume
}

func (d *BtceMarketData) Updated() time.Time {
	return d.data.Updated.Time
}
