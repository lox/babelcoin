package btce

import (
	"strings"
    //"github.com/davecgh/go-spew/spew"
    babel "github.com/lox/babelcoin/core"
    util "github.com/lox/babelcoin/util"
)

const (
	TickerUrl = "https://btc-e.com/api/3/ticker/"
)

type Candle struct {
	Pair 				string				`json:"-"`
	High 				float64				`json:"high"`
	Low 				float64 			`json:"low"`
	Average				float64 			`json:"avg"`
	Volume 				float64				`json:"vol"`
	VolumeCurrent 		float64				`json:"vol_cur"`
	Last 				float64 			`json:"last"`
	Buy 				float64 			`json:"buy"`
	Sell 				float64 			`json:"sell"`
	Updated 			babel.Timestamp	    `json:"updated"`
}

type BtceTickerApi struct {
	url 				string
	currencies 			[]string
}

func NewTickerApi(url string, currencies []string) *BtceTickerApi {
	return &BtceTickerApi{url, currencies}
}

func (t *BtceTickerApi) Candles() ([]Candle, error) {
	var resp map[string]Candle

	error := util.HttpGetJson(t.url + strings.Join(t.currencies, "-"), &resp)
    if error != nil {
    	return nil, error
    }

    var candles []Candle

    for pair, candle := range resp {
    	candle.Pair = pair
    	candles = append(candles, candle)
    }

	return candles, nil
}