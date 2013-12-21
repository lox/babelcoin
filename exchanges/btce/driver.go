package btce

import (
	"github.com/davecgh/go-spew/spew"
	//"github.com/davecgh/go-spew/spew"
	babel "../../core"
	util "../../util"
	"time"
)

type Driver struct {
	config map[string]interface{}
}

func NewDriver(config map[string]interface{}) (*Driver, error) {
	if _, ok := config["info_url"]; !ok {
		config["info_url"] = InfoUrl
	}

	if _, ok := config["ticker_url"]; !ok {
		config["ticker_url"] = TickerUrl
	}

	if _, ok := config["private_url"]; !ok {
		config["private_url"] = PrivateApiUrl
	}

	return &Driver{config}, nil
}

func (b *Driver) Symbols() ([]string, error) {
	api := NewInfoApi(b.config["info_url"].(string))
	pairs, err := api.Pairs()
	if err != nil {
		return []string{}, err
	}

	var symbols []string
	for symbol := range pairs {
		symbols = append(symbols, symbol)
	}

	return symbols, nil
}

func (b *Driver) Ticker(symbol string) (chan babel.MarketData, chan bool, error) {
	api := NewTickerApi(b.config["ticker_url"].(string), []string{symbol})

	duration, ok := b.config["poll_duration"].(time.Duration)
	if !ok {
		duration = time.Duration(5) * time.Second
	}

	channel, quit, err := util.Poller(duration, func() babel.MarketData {
		data, err := api.MarketData()
		if err != nil {
			panic(err)
		}
		return &MarketDataAdaptor{data[0]}
	})

	if err != nil {
		return channel, quit, err
	}

	return channel, quit, nil
}

func (b *Driver) Buy(symbol string, amount float64, price float64) babel.Order {
	return &OrderAdaptor{driver: b, symbol: symbol, t: "buy", amount: amount, rate: price}
}

func (b *Driver) Sell(symbol string, amount float64, price float64) babel.Order {
	return &OrderAdaptor{driver: b, symbol: symbol, t: "sell", amount: amount, rate: price}
}

type MarketDataAdaptor struct {
	data MarketData
}

func (d *MarketDataAdaptor) Ask() float64 {
	return d.data.Sell
}

func (d *MarketDataAdaptor) Bid() float64 {
	return d.data.Buy
}

func (d *MarketDataAdaptor) Last() float64 {
	return d.data.Last
}

func (d *MarketDataAdaptor) Volume() float64 {
	return d.data.Volume
}

func (d *MarketDataAdaptor) Updated() time.Time {
	return d.data.Updated.Time
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

	api := NewBtceApi(o.driver.config["private_url"].(string),
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
	api := NewBtceApi(o.driver.config["private_url"].(string),
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

func (t *TradeAdaptor) Amount() float64 {
	return t.amount
}

func (t *TradeAdaptor) Rate() float64 {
	return t.rate
}
