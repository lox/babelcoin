package btce

import (
	"fmt"
    //"github.com/davecgh/go-spew/spew"
    babelcoin "../util"
)

const (
	DepthUrl = "https://btc-e.com/api/3/depth/"
)

type Order struct {
	Price 				float64
	Amount 				float64
}

type OrderBook struct {
	Asks 				[]Order
	Bids 				[]Order
}

type BtceDepthApi struct {
	url 				string
	currency 			string
	limit				int
}

func NewDepthApi(url string, currency string, limit int) *BtceDepthApi {
	if limit > 2000 {
		limit = 2000
	}
	return &BtceDepthApi{url, currency, limit}
}

func (d *BtceDepthApi) Orders() (OrderBook, error) {
	var resp map[string]map[string][][]float64

	error := babelcoin.HttpGetJson(fmt.Sprintf("%s?limit=%d",d.url + d.currency, d.limit), &resp)
    if error != nil {
    	return OrderBook{}, error
    }

    var book OrderBook

    for _, order := range resp[d.currency]["asks"] {
    	book.Asks = append(book.Asks, Order{order[0], order[1]})
    }

    for _, order := range resp[d.currency]["bids"] {
    	book.Bids = append(book.Bids, Order{order[0], order[1]})
    }

    return book, nil
}
