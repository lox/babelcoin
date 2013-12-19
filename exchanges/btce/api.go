package btce

import (
	babel "../../util"
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	// "net/http/httputil"
	"net/url"
	"strconv"
	"time"
)

const (
	PrivateApiUrl = "https://btc-e.com/tapi"
)

type GetInfoResponse struct {
	Funds            map[string]float64 `json:"funds"`
	Rights           map[string]int     `json:"rights"`
	TransactionCount int                `json:"transaction_count"`
	OpenOrders       int                `json:"open_orders"`
	ServerTime       babel.Timestamp    `json:"server_time"`
}

type TransHistoryResponse struct {
	Id          int             `json:"-"`
	Type        int             `json:"type"`
	Amount      float64         `json:"amount"`
	Currency    string          `json:"currency"`
	Description string          `json:"desc"`
	Status      int             `json:"status"`
	Timestamp   babel.Timestamp `json:"timestamp"`
}

type TradeHistoryResponse struct {
	Id        int             `json:"-"`
	Pair      string          `json:"pair"`
	Type      string          `json:"type"`
	Amount    float64         `json:"amount"`
	Rate      float64         `json:"rate"`
	OrderId   int             `json:"order_id"`
	YourOrder int             `json:"is_your_order"`
	Timestamp babel.Timestamp `json:"timestamp"`
}

type TradeResponse struct {
	Received float64            `json:"received"`
	Remains  float64            `json:"remains"`
	OrderId  int                `json:"order_id"`
	Funds    map[string]float64 `json:"funds"`
}

type ActiveOrdersResponse struct {
	Id        int             `json:"-"`
	Pair      string          `json:"pair"`
	Type      string          `json:"type"`
	Amount    float64         `json:"amount"`
	Rate      float64         `json:"rate"`
	OrderId   int             `json:"-"`
	Status    int             `json:"status"`
	Timestamp babel.Timestamp `json:"timestamp_created"`
}

type CancelOrderResponse struct {
	OrderId int                `json:"order_id"`
	Funds   map[string]float64 `json:"funds"`
}

type BtceApi struct {
	url    string
	key    string
	secret string
}

func NewBtceApi(url string, key string, secret string) *BtceApi {
	return &BtceApi{url, key, secret}
}

// generate hmac-sha512 hash, hex encoded
func Sign(secret string, payload string) string {
	mac := hmac.New(sha512.New, []byte(secret))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}

// returns a url encoded string to be signed an send
func (b *BtceApi) encodePostData(method string, params map[string]string) string {
	// for some reason, btce barfs on larger nonces
	nonce := time.Now().UnixNano() / 1000000000
	result := fmt.Sprintf("method=%s&nonce=%d", method, nonce)

	// params are unordered, but after method and nonce
	if len(params) > 0 {
		v := url.Values{}
		for key := range params {
			v.Add(key, params[key])
		}
		result = result + "&" + v.Encode()
	}

	return result
}

// marshal an api response into an object
func (b *BtceApi) marshalResponse(resp *http.Response, v interface{}) error {
	// read the response
	defer resp.Body.Close()
	bytes, error := ioutil.ReadAll(resp.Body)
	if error != nil {
		return error
	}

	// sometimes, btc-e returns a non-json error
	if string(bytes) == "invalid POST data" {
		return errors.New("Request failed: invalid post data")
	}

	data := &struct {
		Success int
		Return  json.RawMessage
		Error   string
	}{}

	if error = json.Unmarshal(bytes, &data); error != nil {
		return error
	}

	if data.Success != 1 {
		return errors.New("Request failed: " + data.Error)
	}

	if error = json.Unmarshal(data.Return, &v); error != nil {
		return error
	}

	return nil
}

// make a call to the btc-e api, marshal into v
func (b *BtceApi) apiCall(method string, v interface{}, params map[string]string) error {
	client := &http.Client{}
	postData := b.encodePostData(method, params)

	r, error := http.NewRequest("POST", b.url, bytes.NewBufferString(postData))
	if error != nil {
		return error
	}

	r.Header.Add("Sign", Sign(b.secret, postData))
	r.Header.Add("Key", b.key)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(postData)))

	// bytes, _ := httputil.DumpRequest(r, true)
	// spew.Printf("%s", bytes)

	resp, error := client.Do(r)
	if error != nil {
		return error
	}

	return b.marshalResponse(resp, v)
}

// get info about the account
func (b *BtceApi) GetInfo() (GetInfoResponse, error) {
	var resp = GetInfoResponse{}

	if err := b.apiCall("getInfo", &resp, map[string]string{}); err != nil {
		return resp, err
	}

	return resp, nil
}

// returns transaction history
func (b *BtceApi) TransHistory(params map[string]string) ([]TransHistoryResponse, error) {
	var resp = map[string]TransHistoryResponse{}

	if err := b.apiCall("TransHistory", &resp, params); err != nil {
		return nil, err
	}

	var transactions []TransHistoryResponse

	for id, trans := range resp {
		idInt, error := strconv.Atoi(id)
		if error != nil {
			return nil, error
		}

		trans.Id = idInt
		transactions = append(transactions, trans)
	}

	return transactions, nil
}

// returns trade history
func (b *BtceApi) TradeHistory(params map[string]string) ([]TradeHistoryResponse, error) {
	resp := map[string]TradeHistoryResponse{}

	error := b.apiCall("TradeHistory", &resp, params)
	if error != nil {
		return nil, error
	}

	trades := []TradeHistoryResponse{}

	for id, trade := range resp {
		idInt, error := strconv.Atoi(id)
		if error != nil {
			return nil, error
		}

		trade.Id = idInt
		trades = append(trades, trade)
	}

	return trades, nil
}

// return active orders
func (b *BtceApi) ActiveOrders(pair string) ([]ActiveOrdersResponse, error) {
	var resp map[string]ActiveOrdersResponse

	error := b.apiCall("ActiveOrders", &resp, map[string]string{
		"pair": pair,
	})
	if error != nil {
		return nil, error
	}

	var orders []ActiveOrdersResponse

	for id, order := range resp {
		idInt, error := strconv.Atoi(id)
		if error != nil {
			return nil, error
		}

		order.OrderId = idInt
		orders = append(orders, order)
	}

	return orders, nil
}

// performs a trade
func (b *BtceApi) Trade(pair string, t string, rate float64, amount float64) (TradeResponse, error) {
	if t != "buy" && t != "sell" {
		return TradeResponse{}, errors.New("t must be either buy or sell")
	}

	var resp = TradeResponse{}
	error := b.apiCall("Trade", &resp, map[string]string{
		"pair":   pair,
		"type":   t,
		"rate":   strconv.FormatFloat(rate, 'f', -1, 64),
		"amount": strconv.FormatFloat(amount, 'f', -1, 64),
	})
	if error != nil {
		return resp, error
	}

	return resp, nil
}

// cancels an order
func (b *BtceApi) CancelOrder(orderId int) (CancelOrderResponse, error) {
	var resp = CancelOrderResponse{}

	error := b.apiCall("CancelOrder", &resp, map[string]string{
		"order_id": strconv.Itoa(orderId),
	})
	if error != nil {
		return resp, error
	}

	return resp, nil
}
