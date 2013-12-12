package main

import (
    "github.com/docopt/docopt.go"
    //"github.com/davecgh/go-spew/spew"
	"./factory"
	"time"
	"log"
)

func main() {
usage := `Babelcoin. An interface to cryptocoin exchanges

Usage:
  babelcoin ticker <pair> [--interval=<duration>]
  babelcoin currencies <pair> [--interval=<duration>]
  babelcoin trade <pair> (BID|ASK) <amount> [price] [--timeout=<duration>]
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

    if _,ticker := args["ticker"]; ticker {
    	Ticker(args)
	}
}

func Ticker(args map[string]interface{}) {
	pair := args["<pair>"].(string)
	exchange, err := factory.NewExchange(pair)
	if err != nil {
		panic(err)
	}

	duration, err := time.ParseDuration(args["--interval"].(string))
	if err != nil {
		panic(err)
	}

	market, err := exchange.MarketData()
	if err != nil {
		panic(err)
	}

	// anonymous function to pull and format market stats
	tick := func() {
		data, err := market.Fetch()
		if err != nil {
			panic(err)
		}
		log.Printf("Last: %0.6f Ask: %.6f Bid %.6f",
			data.Last(), data.Ask(), data.Bid())
	}

	tick()
	doneChan := make(chan bool)
	ticker := time.NewTicker(duration)
    go func() {
        for _ = range ticker.C {
        	tick()
        }
    }()

    // block until timer is done
    <- doneChan
}

func Currencies(args map[string]interface{}) {
	exchange, err := factory.NewExchange(args["<exchange>"].(string))
	if err != nil {
		panic(err)
	}
}