package btce

import (
	"fmt"
	"strings"
    babel "../util"
)

const (
	TradesUrl = "https://btc-e.com/api/3/trades/"
)

type Trade struct {
	Pair				string 				`json:"-"`
	Type 				string				`json:"type"`
	Price 				float64 			`json:"price"`
	Amount				float64 			`json:"amount"`
	TransactionId		int64 				`json:"tid"`
	Timestamp 			babel.Timestamp	    `json:"Timestamp"`
}

type BtceTradesApi struct {
	url 				string
	currencies 			[]string
	limit				int
}

func NewTradesApi(url string, currencies []string, limit int) *BtceTradesApi {
	if limit > 2000 {
		limit = 2000
	}
	return &BtceTradesApi{url, currencies, limit}
}

func (t *BtceTradesApi) Trades() ([]Trade, error) {
	var resp map[string][]Trade

	error := babel.HttpGetJson(fmt.Sprintf("%s%s?limit=%d",
		t.url, strings.Join(t.currencies, "-"), t.limit), &resp)

    if error != nil {
    	return nil, error
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