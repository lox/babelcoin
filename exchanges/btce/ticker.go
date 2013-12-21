package btce

import (
	"encoding/json"
	"errors"
	"strings"
	//"github.com/davecgh/go-spew/spew"
	babel "../../util"
)

const (
	TickerUrl = "https://btc-e.com/api/3/ticker/"
)

type MarketData struct {
	Pair          string          `json:"-"`
	High          float64         `json:"high"`
	Low           float64         `json:"low"`
	Average       float64         `json:"avg"`
	Volume        float64         `json:"vol"`
	VolumeCurrent float64         `json:"vol_cur"`
	Last          float64         `json:"last"`
	Buy           float64         `json:"buy"`
	Sell          float64         `json:"sell"`
	Updated       babel.Timestamp `json:"updated"`
}

type BtceTickerApi struct {
	url        string
	currencies []string
}

func NewTickerApi(url string, currencies []string) *BtceTickerApi {
	return &BtceTickerApi{url, currencies}
}

func (t *BtceTickerApi) MarketData() ([]MarketData, error) {
	var resp map[string]MarketData

	err := babel.HttpGetJson(t.url+strings.Join(t.currencies, "-"), &resp)
	if err != nil {
		var errorResp = &struct {
			Success int
			Error   string
		}{}

		// check if we got an error encoded in json
		err2 := json.Unmarshal(err.ResponseBody, &errorResp)
		if err2 == nil {
			return nil, errors.New(errorResp.Error)
		}

		return nil, err2
	}

	var data []MarketData

	for pair, row := range resp {
		row.Pair = pair
		data = append(data, row)
	}

	return data, nil
}
