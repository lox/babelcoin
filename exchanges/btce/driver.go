/*
The driver for accessing btce
*/
package btce

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	b "github.com/lox/babelcoin/core"
	util "github.com/lox/babelcoin/util"
)

type Driver struct {
	config     map[string]interface{}
	publicApi  string
	privateApi string
}

func New(exchange string, config map[string]interface{}) *Driver {
	driver := &Driver{config: config}

	if url, ok := config["public_api_url"]; !ok {
		driver.publicApi = "https://btc-e.com/api/3"
	} else {
		driver.publicApi = url.(string)
	}

	if url, ok := config["private_api_url"]; !ok {
		driver.privateApi = "https://btc-e.com/tapi"
	} else {
		driver.privateApi = url.(string)
	}

	return driver
}

func (d *Driver) MarketData(pair b.Pair) (b.MarketData, error) {
	var resp map[string]struct {
		Vol, Last, Buy, Sell float64
		Updated              int64 `json:"updated"`
	}

	if err := util.HttpGetJson(d.publicApi+"/ticker/"+pairString(pair), &resp); err != nil {
		return b.MarketData{}, publicApiError(err)
	}

	t := resp[pairString(pair)]
	return b.MarketData{pair, t.Buy, t.Sell, t.Last, t.Vol, time.Unix(t.Updated, 0)}, nil
}

func (d *Driver) Balance(symbols []b.Symbol) (map[b.Symbol]float64, error) {
	var resp struct {
		Funds map[string]float64 `json:"funds"`
	}
	if err := d.privateApiCall("getInfo", &resp, map[string]string{}); err != nil {
		return map[b.Symbol]float64{}, err
	}

	balances := map[b.Symbol]float64{}
	for symbol, amount := range resp.Funds {
		if len(symbols) == 0 || containsSymbol(b.Symbol(symbol), symbols) {
			balances[b.Symbol(symbol)] = amount
		}
	}

	return balances, nil
}

func (d *Driver) Ticker(pair b.Pair, channel chan<- b.MarketData) error {
	duration, ok := d.config["poll_duration"].(time.Duration)
	if !ok {
		duration = time.Duration(5) * time.Second
	}

	return util.Poller(d, pair, duration, channel)
}

func (d *Driver) Pairs() ([]b.Pair, error) {
	var resp struct {
		Pairs map[string]struct{}
	}

	if err := util.HttpGetJson(d.publicApi+"/info", &resp); err != nil {
		return []b.Pair{}, publicApiError(err)
	}

	pairs := []b.Pair{}
	for k := range resp.Pairs {
		parts := strings.SplitN(k, "_", 2)
		pairs = append(pairs, b.Pair{b.Symbol(parts[0]), b.Symbol(parts[1])})
	}

	return pairs, nil
}

func (d *Driver) Trade(t b.TradeType, pair b.Pair, amount float64, rate float64) (b.Order, error) {
	panic("Not implemented")
}

func (d *Driver) CancelOrder(order b.Order) error {
	panic("Not implemented")
}

func (d *Driver) History(pair b.Pair, after time.Time, channel chan<- b.Trade) error {
	var resp map[string][]struct {
		Type           b.TradeType
		Price, Amount  float64
		Tid, Timestamp int64
	}

	// the since param doesn't seem to work any more
	url := fmt.Sprintf("%s/trades/%s?limit=%d&since=%d",
		d.publicApi, pairString(pair), 2000, after.Unix())

	if err := util.HttpGetJson(url, &resp); err != nil {
		return publicApiError(err)
	}

	for _, t := range resp[pairString(pair)] {
		channel <- b.Trade{
			strconv.FormatInt(t.Tid, 10), pair, t.Amount,
			t.Price, time.Unix(t.Timestamp, 0), t.Type,
		}
	}

	close(channel)
	return nil
}

func (d *Driver) Orders(limit int) ([]b.Order, error) {
	panic("Not implemented")
}

func (d *Driver) Transactions(limit int) ([]b.Transaction, error) {
	panic("Not implemented")
}

func (d *Driver) OrderBook(pair b.Pair, limit int) (b.OrderBook, error) {
	panic("Not implemented")
}

/*
func (b *Driver) History(symbol string) (chan babel.Trade, error) {
	//_ = NewTradesApi(b.config["trades_url"].(string), []string{symbol}, 2000)
	channel := make(chan babel.Trade, 1)
	close(channel)
	return channel, nil
}

func (b *Driver) Trade(t babel.TradeType, symbol string, amount float64, price float64) babel.Order {
	return &OrderAdaptor{driver: b, symbol: symbol, t: "buy", amount: amount, rate: price}
}

func (b *Driver) Sell(symbol string, amount float64, price float64) babel.Order {
	return &OrderAdaptor{driver: b, symbol: symbol, t: "sell", amount: amount, rate: price}
}

type OrderAdaptor struct {
	driver  *Driver
	symbol  string
	t       string
	amount  float64
	rate    float64
	orderId int
}

func (o *OrderAdaptor) Execute() (chan babel.Trade, error) {
	channel := make(chan babel.Trade, 10)

	api := NewDriver(o.driver.config["private_url"].(string),
		o.driver.config["key"].(string), o.driver.config["secret"].(string))

	unixtime := time.Now().Unix()
	resp, err := api.Trade(o.symbol, o.t, o.rate, o.amount)
	if err != nil {
		close(channel)
		return channel, err
	}

	o.orderId = resp.OrderId
	// spew.Dump(resp)

	if resp.Remains > 0 {
		// start := time.Now()
		// remains := resp.Remains
		ticker := time.NewTicker(15 * time.Second)
		go func() {
			for _ = range ticker.C {
				spew.Println("Tick")
				trades, _ := api.TradeHistory(map[string]string{
					"pair":  o.symbol,
					"since": string(unixtime),
				})

				spew.Printf("%d trades\n", len(trades))
				for trade := range trades {
					spew.Dump(trade)
				}
			}
		}()
	} else {
		spew.Dump("No remains!")
		channel <- &TradeAdaptor{o.amount, o.rate}
		close(channel)
	}

	return channel, nil
}

func (o *OrderAdaptor) Pair() babel.Pair {
	return ""
}

func (o *OrderAdaptor) Amount() float64 {
	return 0
}

func (o *OrderAdaptor) Remains() float64 {
	return 0
}

func (o *OrderAdaptor) Timestamp() time.Time {
	return time.Now()
}

func (o *OrderAdaptor) Type() babel.TradeType {
	return babel.Buy
}

func (o *OrderAdaptor) Fee() (float64, error) {
	api := NewInfoApi(o.driver.config["info_url"].(string))
	pairs, err := api.Pairs()

	if err != nil {
		return 0.0, err
	}

	pair, ok := pairs[o.symbol]
	if !ok {
		return 0.0, err
	}

	return pair.Fee, nil
}

func (o *OrderAdaptor) Cancel() error {
	api := NewDriver(o.driver.config["private_url"].(string),
		o.driver.config["key"].(string), o.driver.config["secret"].(string))

	_, err := api.CancelOrder(o.orderId)
	if err != nil {
		return err
	}

	return nil
}

type TradeAdaptor struct {
	amount float64
	rate   float64
}

func (t *TradeAdaptor) Pair() babel.Pair {
	return ""
}

func (t *TradeAdaptor) Type() babel.TradeType {
	return babel.Buy
}

func (t *TradeAdaptor) Timestamp() time.Time {
	return time.Now()
}

func (t *TradeAdaptor) Amount() float64 {
	return t.amount
}

func (t *TradeAdaptor) Rate() float64 {
	return t.rate
}
*/

// attempts to extract a message from an api error
func publicApiError(err *util.HttpError) error {
	var er struct {
		Success int
		Error   string
	}

	if err.ResponseBody == nil {
		return err
	} else if err2 := json.Unmarshal(err.ResponseBody, &er); err2 != nil {
		return err
	}

	return errors.New("API Error: " + er.Error)
}

// returns a pair in the form ltc_usd
func pairString(pair b.Pair) string {
	return string(pair.Base) + "_" + string(pair.Counter)
}

// checks if a Symbol is in a slice of Symbols
func containsSymbol(a b.Symbol, list []b.Symbol) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
