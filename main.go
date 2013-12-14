package main

import (
	"./core"
	"./factory"
	"fmt"
	"github.com/docopt/docopt.go"
	"log"
	"time"
)

func main() {
	usage := `Babelcoin. An interface to cryptocoin exchanges

Usage:
  babelcoin ticker <symbol> [--interval=<duration>]
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

	if ticker := args["ticker"]; ticker.(bool) {
		Ticker(args)
	} else if symbols := args["symbols"]; symbols.(bool) {
		Symbols(args)
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
