package btce

import (
	babel "../../util"
	"fmt"
	"strings"
)

const (
	TradesUrl = "https://btc-e.com/api/3/trades/"
)

type Trade struct {
	Pair          string          `json:"-"`
	Type          string          `json:"type"`
	Price         float64         `json:"price"`
	Amount        float64         `json:"amount"`
	TransactionId int64           `json:"tid"`
	Timestamp     babel.Timestamp `json:"Timestamp"`
}

type BtceTradesApi struct {
	url        string
	currencies []string
	limit      int
}

func NewTradesApi(url string, currencies []string, limit int) *BtceTradesApi {
	if limit > 2000 {
		limit = 2000
	}
	return &BtceTradesApi{url, currencies, limit}
}

func (t *BtceTradesApi) Trades() ([]Trade, error) {
	var resp map[string][]Trade

	err := babel.HttpGetJson(fmt.Sprintf("%s%s?limit=%d",
		t.url, strings.Join(t.currencies, "-"), t.limit), &resp)

	if err != nil {
		return nil, err
	}

	var trades []Trade

	for pair, row := range resp {
		for _, trade := range row {
			trade.Pair = pair
			trades = append(trades, trade)
		}
	}

	return trades, nil
}
