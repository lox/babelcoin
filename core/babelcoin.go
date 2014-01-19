/*
Babelcoin provides a generic interface to cryptocurrency exchanges.

The interfaces defined in this file are what each driver provides,
beyond the native interface that is defined by the underlying exchange
API.

This interface is subject to change at this stage.
*/
package babelcoin

import (
	"fmt"
	"strings"
	"time"
)

// a symbol represents a currency type, e.g btc, usd, ltc
type Symbol string

const (
	BTC Symbol = "btc"
	LTC Symbol = "ltc"
	FTC Symbol = "ftc"
	USD Symbol = "usd"
	AUD Symbol = "aud"
	EUR Symbol = "eur"
)

// a pair represents the trading between two symbols, e.g btc/usd
type Pair struct {
	Base, Counter Symbol
}

var (
	BTC_USD Pair = Pair{BTC, USD}
	LTC_USD Pair = Pair{LTC, USD}
	LTC_BTC Pair = Pair{LTC, BTC}
)

// the state of trading for a given pair on an exchange
type MarketData struct {
	Pair                    Pair
	Buy, Sell, Last, Volume float64
	Updated                 time.Time
}

// the type of a trade, either buy or sell
type TradeType string

const (
	Buy  TradeType = "buy"
	Sell TradeType = "sell"
)

// a single trade that has been executed on a market
type Trade struct {
	Id           string
	Pair         Pair
	Amount, Rate float64
	Timestamp    time.Time
	Type         TradeType
}

// an order on an exchange, either ours or other peoples
type Order struct {
	Id                         string
	Pair                       Pair
	Type                       TradeType
	Timestamp                  time.Time
	Amount, Remains, Rate, Fee float64
}

// an operation against an account
type Transaction struct {
	Id        string
	Symbol    Symbol
	Timestamp time.Time
	Amount    float64
}

// the order book showing asks and bids
type OrderBook struct {
	Asks, Bids []struct {
		Price, Amount float64
	}
}

// the interface to an exchange
type Exchange interface {
	// the users balance for the provided symbol, an empty
	// slice should result in all balances being returned
	Balance(symbols []Symbol) (map[Symbol]float64, error)

	// returns the current market state
	MarketData(pair Pair) (MarketData, error)

	// get a live feed of market data
	Ticker(pair Pair, channel chan<- MarketData) error

	// returns the pairs that are supported on the exchange
	Pairs() ([]Pair, error)

	// executes a trade on the exchange, either as a limit order if a rate
	// is provided, or a market order if -1 is provided as rate. If amount is -1 then
	// the entire balance the user has is used
	Trade(t TradeType, pair Pair, amount float64, rate float64) (Order, error)

	// cancels an order that was previously placed
	CancelOrder(order Order) error

	// returns historical trades for the exchange for the provided timeframe
	History(pair Pair, after time.Time, channel chan<- Trade) error

	// returns the users orders
	Orders(limit int) ([]Order, error)

	// returns the users transactions
	Transactions(limit int) ([]Transaction, error)

	// returns an order book of a given depth for the given pair
	// accepts a limit to limit to the top N orders
	OrderBook(pair Pair, limit int) (OrderBook, error)
}

// a function for creating an Exchange
type ExchangeFactory func(key string, config map[string]interface{}) Exchange

// returns a pair in the form btc_usd as a string
func (p *Pair) String() string {
	return string(p.Base + "_" + p.Counter)
}

// parses a pair in the form of btc_usd
func ParsePair(pair string) Pair {
	parts := strings.SplitN(pair, "_", 2)
	return Pair{Symbol(strings.ToLower(parts[0])), Symbol(strings.ToLower(parts[1]))}
}

// returns a string version of a trade
func (t *Trade) String() string {
	return fmt.Sprintf("%s %s %.5f@%.5f on %s",
		strings.ToUpper(string(t.Type)), t.Pair.String(),
		t.Amount, t.Rate, t.Timestamp.Format(time.Stamp))
}
