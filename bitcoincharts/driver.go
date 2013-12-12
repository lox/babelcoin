package bitcoincharts

import (
	"../core"
	"time"
	"errors"
)

type Exchange struct {
}

func NewExchange() *Exchange {
	return &Exchange{}
}

func (b *Exchange) MarketData(symbol string) (babelcoin.MarketDataService, error) {
	return &MarketDataService{symbol}, nil
}

func (b *Exchange) Symbols() ([]string, error) {
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

type MarketDataService struct {
	symbol string
}

func (b *MarketDataService) Fetch() (babelcoin.MarketData, error) {
	api := NewMarketsApi(MarketUrl)
	markets, err := api.Markets()
	if err != nil {
		return nil, err
	}

	market, ok := markets[b.symbol]
	if !ok {
		return nil, errors.New("No market data found for "+b.symbol)
	}

	return &BitcoinChartsMarketData{market}, nil
}

type BitcoinChartsMarketData struct {
	data Market
}

func (d *BitcoinChartsMarketData) Ask() float64 {
	return d.data.Ask
}

func (d *BitcoinChartsMarketData) Bid() float64 {
	return d.data.Bid
}

func (d *BitcoinChartsMarketData) Last() float64 {
	return d.data.Close
}

func (d *BitcoinChartsMarketData) Volume() float64 {
	return d.data.Volume
}

func (d *BitcoinChartsMarketData) Updated() time.Time {
	return d.data.LatestTrade.Time
}