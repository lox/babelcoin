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

// get a feed of orders
feed := exchange.OrderFeed()

// feed is just a go channel
for order := range feed {
	fmt.Printf("Order: Amount: %.2f Price: %.2f\n", order.Amount, order.Price)
}

// place a limit order
limitOrder := exchange.NewLimitOrder(100.0, 11.0)
limitOrder.Execute()

// place a market order
marketOrder := exchange.NewLimitOrder(100.0, 11.0)
marketOrder.Execute()
```










