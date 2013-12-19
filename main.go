package main

import (
	"./core"
	"./factory"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	// "github.com/davecgh/go-spew/spew"
	"github.com/docopt/docopt.go"
	"log"
	"strconv"
	"time"
)

func main() {
	usage := `Babelcoin. An interface to cryptocoin exchanges

Usage:
  babelcoin ticker <symbol> [--interval=<duration>]
  babelcoin (buy|sell) <symbol> <amount> <rate> [--timeout=<duration>]
  babelcoin symbols <exchange>
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

	// spew.Dump(args)
	//

	if ticker := args["ticker"]; ticker.(bool) {
		Ticker(args)
	} else if symbols := args["symbols"]; symbols.(bool) {
		Symbols(args)
	} else if buy := args["buy"]; buy.(bool) {
		Trade(args, "buy")
	} else if sell := args["sell"]; sell.(bool) {
		Trade(args, "sell")
	}
}

func Ticker(args map[string]interface{}) {
	symbol := babelcoin.ParseSymbol(args["<symbol>"].(string))
	pair, err := symbol.Pair()
	if err != nil {
		panic(err)
	}

	duration, err := time.ParseDuration(args["--interval"].(string))
	if err != nil {
		panic(err)
	}

	exchange, err := factory.NewExchangeWithConfig(symbol.Exchange(), map[string]interface{}{
		"poll_duration": duration,
	})
	if err != nil {
		panic(err)
	}

	ticker, _, err := exchange.Ticker(pair)
	if err != nil {
		panic(err)
	}

	for data := range ticker {
		log.Printf("Last: %0.6f Ask: %.6f Bid %.6f Volume: %.2f",
			data.Last(), data.Ask(), data.Bid(), data.Volume())
	}
}

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
