/*
The driver for accessing bitcoincharts,
supports History(), Ticker(), Pairs() and MarketData()
*/
package bitcoincharts

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	b "github.com/lox/babelcoin/core"
	util "github.com/lox/babelcoin/util"
)

type Driver struct {
	exchange string
	config   map[string]interface{}
}

// creates a new bitcoincharts driver
func New(exchange string, config map[string]interface{}) b.Exchange {
	if _, ok := config["api_url"]; !ok {
		config["api_url"] = "http://api.bitcoincharts.com/v1"
	}

	return &Driver{
		exchange: exchange,
		config:   config,
	}
}

func (d *Driver) MarketData(pair b.Pair) (b.MarketData, error) {
	var resp []struct {
		Symbol      string        `json:"symbol"`
		Bid         float64       `json:"bid"`
		Ask         float64       `json:"ask"`
		LatestTrade util.UnixTime `json:"latest_trade"`
		Close       float64       `json:"close"`
		Volume      float64       `json:"volume"`
	}

	err := util.HttpGetJson(d.config["api_url"].(string)+"/markets.json", &resp)
	if err != nil {
		return b.MarketData{}, err
	}

	for _, data := range resp {
		if data.Symbol == d.getSymbol(pair) {
			return b.MarketData{
				Pair:    pair,
				Last:    data.Close,
				Buy:     data.Bid,
				Sell:    data.Ask,
				Volume:  data.Volume,
				Updated: data.LatestTrade.Time,
			}, nil
		}
	}

	return b.MarketData{}, errors.New("Unknown pair " + pair.String())
}

func (d *Driver) TradeHistory(pairs []b.Pair, after time.Time, limit int, channel chan<- b.Trade) error {
	tempChannel := make(chan b.Trade, 10)

	for _, pair := range pairs {
		reader, err := d.getHistoryCsv(pair)
		if err != nil {
			return err
		}
		go func() {
			d.readAllCsvTrades(pair, reader, tempChannel)
		}()
	}

	// process trades for time limiting
	go func() {
		for trade := range tempChannel {
			if trade.Timestamp.After(after) {
				channel <- trade
			}
		}
		close(channel)
	}()

	return nil
}

func (d *Driver) Ticker(pair b.Pair, channel chan<- b.MarketData) error {
	duration, ok := d.config["poll_duration"].(time.Duration)
	if !ok {
		duration = time.Duration(5) * time.Second
	}

	ticker := time.NewTicker(duration)
	go func() {
		for _ = range ticker.C {
			data, err := d.MarketData(pair)
			if err != nil {
				panic(err)
			}
			channel <- data
		}
		close(channel)
	}()

	return nil
}

func (d *Driver) Pairs() ([]b.Pair, error) {
	var resp []struct {
		Symbol   string `json:"symbol"`
		Currency string `json:"currency"`
	}

	err := util.HttpGetJson(d.config["api_url"].(string)+"/markets.json", &resp)
	if err != nil {
		return []b.Pair{}, err
	}

	parts := strings.Split(d.exchange, ":")
	var pairs []b.Pair

	if len(parts) != 2 {
		panic("Exchange name must be in bitcoincharts:xxxx format")
	}

	for _, data := range resp {
		if strings.Index(data.Symbol, parts[1]) == 0 {
			pairs = append(pairs, b.Pair{
				b.Symbol("btc"), b.Symbol(strings.ToLower(data.Currency)),
			})
		}
	}

	return pairs, nil
}

func (d *Driver) Account() b.ExchangeAccount {
	panic("Not implemented")
}

func (d *Driver) getSymbol(pair b.Pair) string {
	parts := strings.SplitN(d.exchange, ":", 2)
	return parts[1] + strings.ToUpper(string(pair.Counter))
}

// read all trades in csv format from the bitcoincharts api
func (d *Driver) readAllCsvTrades(pair b.Pair, reader io.Reader, channel chan<- b.Trade) error {
	csv := csv.NewReader(reader)
	for {
		fields, err := csv.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		timestamp, _ := strconv.ParseInt(fields[0], 10, 64)
		rate, _ := strconv.ParseFloat(fields[1], 64)
		amount, _ := strconv.ParseFloat(fields[2], 64)
		channel <- b.Trade{
			Pair:      pair,
			Timestamp: time.Unix(timestamp, 0),
			Rate:      rate,
			Amount:    amount,
		}
	}
	close(channel)
	return nil
}

// gets a reader for the csv data for full history for a pair
func (d *Driver) getHistoryCsv(pair b.Pair) (io.Reader, error) {

	// if provided, use a cache dir for the csvs
	if cache := os.Getenv("BTCCHARTS_CACHE"); cache != "" {
		filename := fmt.Sprintf("%s/%s.csv", cache, d.getSymbol(pair))
		log.Printf("Using cached bitcoincharts history file %s", filename)
		return os.Open(filename)
	}

	url := fmt.Sprintf("%s/csv/%s.csv", d.config["api_url"], d.getSymbol(pair))
	log.Printf("Downloading full history from %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func init() {
	b.AddExchangeFactory("bitcoincharts", b.ExchangeFactory(New))
}
