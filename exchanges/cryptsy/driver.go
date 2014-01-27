/*
The driver for accessing cryptsy
*/
package cryptsy

import (
	"errors"
	"strconv"
	"strings"
	"time"

	b "github.com/lox/babelcoin/core"
	util "github.com/lox/babelcoin/util"
)

type Driver struct {
	exchange string
	config   map[string]interface{}
	markets  map[string]market
	client   *util.JsonRPCClient
}

type market struct {
	Label, MarketId, Created string
}

// creates a new cryptsy driver
func New(exchange string, config map[string]interface{}) b.Exchange {
	if _, ok := config["private_api_url"]; !ok {
		config["private_api_url"] = "https://www.cryptsy.com/api"
	}

	return &Driver{
		exchange: exchange,
		config:   config,
		client: &util.JsonRPCClient{
			config["private_api_url"].(string),
			config["key"].(string),
			config["secret"].(string),
		},
	}
}

func (d *Driver) MarketData(pair b.Pair) (b.MarketData, error) {
	panic("Not implemented")
}

func (d *Driver) TradeHistory(pairs []b.Pair, after time.Time, limit int, channel chan<- b.Trade) error {
	markets, err := d.getMarkets()
	if err != nil {
		return err
	}

	var resp []struct {
		TradeId    string
		DateTime   string
		TradePrice string
		Quantity   string
		Total      string
		OrderType  string `json:"initiate_ordertype"`
	}

	for _, pair := range pairs {
		market, ok := markets[pair.String()]
		if !ok {
			return errors.New("Unknown pair " + pair.String())
		}

		if err := d.client.Call("markettrades", &resp,
			map[string]string{"marketid": market.MarketId}); err != nil {
			return err
		}
		for _, trade := range resp {
			rate, _ := strconv.ParseFloat(trade.TradePrice, 64)
			amount, _ := strconv.ParseFloat(trade.Quantity, 64)

			t, err := time.Parse("2006-01-02 15:04:05", trade.DateTime)
			if err != nil {
				return err
			}

			channel <- b.Trade{
				Pair:      pair,
				Amount:    amount,
				Rate:      rate,
				Exchange:  "cryptsy",
				Timestamp: t,
				Type:      b.TradeType(strings.ToUpper(trade.OrderType)),
			}
		}
	}

	close(channel)
	return nil
}

func (d *Driver) Ticker(pair b.Pair, channel chan<- b.MarketData) error {
	panic("Not implemented")
}

func (d *Driver) Pairs() ([]b.Pair, error) {
	markets, err := d.getMarkets()
	if err != nil {
		return nil, err
	}

	var pairs []b.Pair
	for k, _ := range markets {
		pairs = append(pairs, b.ParsePair(k))
	}
	return pairs, nil
}

func (d *Driver) getMarkets() (map[string]market, error) {
	if d.markets == nil {
		var resp []market
		if err := d.client.Call("getmarkets", &resp, map[string]string{}); err != nil {
			return nil, err
		}

		d.markets = map[string]market{}
		for _, market := range resp {
			d.markets[strings.ToLower(strings.Replace(market.Label, "/", "_", -1))] = market
		}
	}
	return d.markets, nil
}

func (d *Driver) Account() b.ExchangeAccount {
	panic("Not implemented")
}

func init() {
	b.AddExchangeFactory("cryptsy", b.ExchangeFactory(New))
}
