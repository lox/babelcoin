package btce

import (
	//"github.com/davecgh/go-spew/spew"
	babel "../../util"
)

const (
	InfoUrl = "https://btc-e.com/api/3/info/"
)

type InfoResponse struct {
	Pairs      map[string]CurrencyPair `json:"pairs"`
	ServerTime babel.Timestamp         `json:"server_time"`
}

type CurrencyPair struct {
	DecimalPlaces int     `json:"decimal_places"`
	MinPrice      float64 `json:"min_price"`
	MaxPrice      float64 `json:"max_price"`
	MinAmount     float64 `json:"min_amount"`
	Hidden        float64 `json:"hidden"`
	Fee           float64 `json:"fee"`
}

type BtceInfoApi struct {
	url string
}

func NewInfoApi(url string) *BtceInfoApi {
	return &BtceInfoApi{url}
}

func (ticker *BtceInfoApi) Pairs() (map[string]CurrencyPair, error) {
	var resp InfoResponse

	if err := babel.HttpGetJson(ticker.url, &resp); err != nil {
		return nil, err
	}

	return resp.Pairs, nil
}
