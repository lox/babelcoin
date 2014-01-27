BabelCoin
=========

A generic interface to access and interoperate between various cryptocurrency exchanges.

Installing
----------

```bash
go get -u github.com/lox/babelcoin
```

Usage
-----

```go
import "github.com/lox/babelcoin"
import "fmt"

exchange := babelcoin.NewExchange("btce", "usd_ltc")

// get the latest market data
ticker, _, _ := exchange.Ticker()

// ticker is just a go channel
for data := range ticker {
	fmt.Printf("Last: %.6f\n", data.Last())
}

// place a limit bid order
limitOrder := exchange.NewBid(100.0, 11.0, -1)
trades, err := limitOrder.Execute()
```

Status
-------------------

In the middle of refactoring API. Trade/MarketData polling is working for BTC-e.