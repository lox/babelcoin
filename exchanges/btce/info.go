package btce

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	//"github.com/davecgh/go-spew/spew"
	babel "../../util"
)

const (
	InfoUrl = "https://btc-e.com/api/3/info/"
)

type InfoResponse struct {
	Pairs      map[string]CurrencyPair `json:"Pairs"`
	ServerTime babel.Timestamp         `json:"Server_time"`
}

type CurrencyPair struct {
	DecimalPlaces int     `json:"Decimal_places"`
	MinPrice      float64 `json:"Min_price"`
	MaxPrice      float64 `json:"Max_price"`
	MinAmount     float64 `json:"Min_amount"`
	Hidden        float64 `json:"Hidden"`
	Fee           float64 `json:"Fee"`
}

type BtceInfoApi struct {
	url string
}

func NewInfoApi(url string) *BtceInfoApi {
	return &BtceInfoApi{url}
}

func (ticker *BtceInfoApi) Pairs() (map[string]CurrencyPair, error) {
	var data InfoResponse

	resp, error := http.Get(ticker.url)
	if error != nil {
		return data.Pairs, error
	}

	// read the response
	defer resp.Body.Close()
	bytes, error := ioutil.ReadAll(resp.Body)
	if error != nil {
		return data.Pairs, error
	}

	if error = json.Unmarshal(bytes, &data); error != nil {
		return data.Pairs, error
	}

	return data.Pairs, nil
}
