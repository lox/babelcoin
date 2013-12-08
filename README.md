CoinX
=========

A generic interface to access and interoperate between various cryptocurrency exchanges.

Installing
----------


Usage
-----

import "github.com/lox/coinx"
import "fmt"

exchange := coinx.NewExchange("btce", "usd/ltc")

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










