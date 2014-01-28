/*
The driver for accessing cryptsy
*/
package cryptsy

import (
	"log"
	"strconv"
	"strings"
	"time"

	b "github.com/lox/babelcoin/core"
	util "github.com/lox/babelcoin/util"
)

type Driver struct {
	exchange       string
	config         map[string]interface{}
	markets        map[b.Pair]market
	client         *util.JsonRPCClient
	serverLocation *time.Location
}

type market struct {
	Label, MarketId, Created string
}

// creates a new cryptsy driver
func New(exchange string, config map[string]interface{}) b.Exchange {
	if _, ok := config["private_api_url"]; !ok {
		config["private_api_url"] = "https://www.cryptsy.com/api"
	}

	if _, ok := config["key"]; !ok {
		panic("Missing key for cryptsy")
	}

	if _, ok := config["secret"]; !ok {
		panic("Missing secret for cryptsy")
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
	markets, err := d.getMarketsByPairs(pairs)
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

	for pair, market := range markets {
		log.Printf("Querying pair %v on cryptsy", pair.String())
		if err := d.client.Call("markettrades", &resp,
			map[string]string{"marketid": market.MarketId}); err != nil {
			return err
		}

		for _, trade := range resp {
			rate, _ := strconv.ParseFloat(trade.TradePrice, 64)
			amount, _ := strconv.ParseFloat(trade.Quantity, 64)

			t, err := time.ParseInLocation("2006-01-02 15:04:05", trade.DateTime, d.location())
			if err != nil {
				log.Printf("Failed to parse trade time %s: %v", trade.DateTime, err)
				return err
			}

			channel <- b.Trade{
				Id:        trade.TradeId,
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
		pairs = append(pairs, k)
	}
	return pairs, nil
}

func (d *Driver) getMarketsByPairs(pairs []b.Pair) (map[b.Pair]market, error) {
	markets, err := d.getMarkets()
	if err != nil {
		return nil, err
	}

	// prune markets not in the passed array
	for key, _ := range markets {
		if !b.ContainsPair(key, pairs) {
			log.Printf("%v isn't in the pairs array", key)
			delete(markets, key)
		}
	}

	return markets, nil
}

func (d *Driver) getMarkets() (map[b.Pair]market, error) {
	if d.markets == nil {
		var resp []market
		if err := d.client.Call("getmarkets", &resp, map[string]string{}); err != nil {
			return nil, err
		}

		d.markets = map[b.Pair]market{}
		for _, market := range resp {
			pair := b.ParsePair(strings.Replace(market.Label, "/", "_", -1))
			d.markets[pair] = market
		}
	}
	return d.markets, nil
}

func (d *Driver) Account() b.ExchangeAccount {
	panic("Not implemented")
}

func (d *Driver) location() *time.Location {
	if d != nil {
		// urgh https://cryptsy.freshdesk.com/support/discussions/topics/30997
		location, err := time.LoadLocation("EST5EDT")
		if err != nil {
			panic(err)
		}
		d.serverLocation = location
	}
	return d.serverLocation
}

func init() {
	b.AddExchangeFactory("cryptsy", b.ExchangeFactory(New))
}
