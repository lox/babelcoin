package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/docopt/docopt.go"
	"github.com/lox/babelcoin/core"
	"github.com/lox/babelcoin/exchanges/bitcoincharts"
	"github.com/lox/babelcoin/exchanges/btce"
)

func main() {
	usage := `Babelcoin. An interface to cryptocoin exchanges

Usage:
  babelcoin ticker <exchange> <pair> [--interval=<duration>]
  babelcoin history <exchange> <pair>
  babelcoin (buy|sell) <exchange> <pair> <amount> <rate> [--timeout=<duration>]
  babelcoin pairs <exchange>
  babelcoin balances <exchange>
  babelcoin -h | --help
  babelcoin --version

Options:
  -h --help     			Show this screen.
  --version     			Show version.
  -i --interval=<duration>  Time interval to use [default: 30s].
  --timeout=<duration>  	A timeout to cancel the order by if not filled.`

	args, err := docopt.Parse(usage, nil, true, "Babelcoin", false)
	if err != nil {
		panic(err)
	}

	if ticker := args["ticker"]; ticker.(bool) {
		Ticker(args)
	} else if pairs := args["pairs"]; pairs.(bool) {
		Pairs(args)
	} else if buy := args["buy"]; buy.(bool) {
		//Trade(args, "buy")
	} else if sell := args["sell"]; sell.(bool) {
		//Trade(args, "sell")
	} else if history := args["history"]; history.(bool) {
		History(args)
	} else if balances := args["balances"]; balances.(bool) {
		Balances(args)
	}
}

func Ticker(args map[string]interface{}) {
	duration, err := time.ParseDuration(args["--interval"].(string))
	if err != nil {
		panic(err)
	}

	exchange, err := NewExchange(args["<exchange>"].(string), map[string]interface{}{
		"poll_duration": duration,
	})
	if err != nil {
		panic(err)
	}

	channel := make(chan babelcoin.MarketData, 10)
	err = exchange.Ticker(babelcoin.ParsePair(args["<pair>"].(string)), channel)
	if err != nil {
		panic(err)
	}

	for data := range channel {
		log.Printf("Last: %0.6f Sell: %.6f Buy %.6f Volume: %.2f",
			data.Last, data.Sell, data.Buy, data.Volume)
	}
}

func History(args map[string]interface{}) {
	exchange, err := NewExchange(args["<exchange>"].(string), map[string]interface{}{})
	if err != nil {
		panic(err)
	}

	// get history for up to 2 months ago
	after := time.Now().AddDate(-1, -2, 0)
	channel := make(chan babelcoin.Trade)

	if err := exchange.History(babelcoin.ParsePair(args["<pair>"].(string)), after, channel); err != nil {
		panic(err)
	}

	log.Printf("Loading history after %s", after)
	for trade := range channel {
		fmt.Printf("%s %s %.4f @ %.4f\n",
			trade.Timestamp.Format("2006-01-02T15:04:05"), trade.Type, trade.Amount, trade.Rate)
	}
}

func Pairs(args map[string]interface{}) {
	exchange, err := NewExchange(args["<exchange>"].(string), map[string]interface{}{})
	if err != nil {
		panic(err)
	}

	pairs, err := exchange.Pairs()
	if err != nil {
		panic(err)
	}

	for _, pair := range pairs {
		fmt.Printf("%s\n", pair.String())
	}
}

func Balances(args map[string]interface{}) {
	exchange, err := NewExchange(args["<exchange>"].(string), map[string]interface{}{})
	if err != nil {
		panic(err)
	}

	balances, err := exchange.Balance([]babelcoin.Symbol{})
	if err != nil {
		panic(err)
	}

	for symbol, amount := range balances {
		fmt.Printf("%s => %.8f\n", symbol, amount)
	}
}

/*
func Symbols(args map[string]interface{}) {
	exchange, err := factory.NewExchange(args["<exchange>"].(string))
	if err != nil {
		panic(err)
	}

	symbols, err := exchange.Symbols()
	if err != nil {
		panic(err)
	}

	for _, symbol := range symbols {
		fmt.Printf("%s\n", symbol)
	}
}

func Trade(args map[string]interface{}, t string) {
	symbol := babelcoin.ParseSymbol(args["<symbol>"].(string))
	exchange, err := factory.NewExchange(symbol.Exchange())
	if err != nil {
		panic(err)
	}

	pair, err := symbol.Pair()
	if err != nil {
		panic(err)
	}

	if _, ok := args["<amount>"]; !ok {
		panic("Must provide a rate")
	}

	amount, err := strconv.ParseFloat(args["<amount>"].(string), 64)
	if err != nil {
		panic(err)
	}

	rate, err := strconv.ParseFloat(args["<rate>"].(string), 64)
	if err != nil {
		panic(err)
	}

	timeout, err := time.ParseDuration(args["--timeout"].(string))
	if err != nil {
		panic(err)
	}

	var order babelcoin.Order

	if t == "buy" {
		log.Printf("Buying %.4f %s @ %.4f", amount, pair, rate)
		order = exchange.Buy(pair, amount, rate)
	} else if t == "sell" {
		log.Printf("Selling %.4f %s @ %.4f", amount, pair, rate)
		order = exchange.Sell(pair, amount, rate)
	}

	trades, err := order.Execute()
	if err != nil {
		panic(err)
	}

poll:
	for {
		select {
		case trade, ok := <-trades:
			if !ok {
				log.Printf("Order completed")
				break poll
			} else {
				spew.Dump(trade)
			}
		case <-time.After(timeout):
			log.Printf("Order timed out, cancelling")
			close(trades)
			err := order.Cancel()
			if err != nil {
				panic(err)
			}
			log.Printf("Order cancelled")
			break poll
		}
	}
}

*/

// parse the name of an exchange and return an instance
func NewExchange(exchange string, config map[string]interface{}) (babelcoin.Exchange, error) {
	parts := strings.SplitN(exchange, ":", 2)

	switch parts[0] {
	case "bitcoincharts":
		return bitcoincharts.New(exchange, config), nil
	case "btce":
		config["key"] = os.Getenv("BTCE_KEY")
		config["secret"] = os.Getenv("BTCE_SECRET")
		return btce.New(exchange, config), nil
	}

	return nil, errors.New("Unknown exchange " + exchange)
}
