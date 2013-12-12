package bitcoincharts

import (
	babel "../util"
)

const (
	MarketUrl = "http://api.bitcoincharts.com/v1/markets.json"
)

type Market struct {
	Symbol         string          `json:"symbol"`
	Currency       string          `json:"currency"`
	Bid            float64         `json:"bid"`
	Ask            float64         `json:"ask"`
	LatestTrade    babel.Timestamp `json:"latest_trade"`
	Open           float64         `json:"open"`
	High           float64         `json:"high"`
	Low            float64         `json:"low"`
	Close          float64         `json:"close"`
	PreviousClose  float64         `json:"previous_close"`
	Volume         float64         `json:"volume"`
	CurrencyVolume float64         `json:"currency_volume"`
}

type MarketsApi struct {
	url     string
	markets []Market
}

func NewMarketsApi(url string) *MarketsApi {
	return &MarketsApi{url: url}
}

func (t *MarketsApi) Markets() (map[string]Market, error) {
	var resp []Market

	error := babel.HttpGetJson(t.url, &resp)
	if error != nil {
		return nil, error
	}

	markets := map[string]Market{}
	for _, market := range resp {
		markets[market.Symbol] = market
	}

	return markets, nil
}
