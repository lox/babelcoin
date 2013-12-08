package btce

import (
    "bytes"
    "fmt"
    "encoding/json"
	"net/http"
	"io/ioutil"
	"net/url"
	"time"
	"errors"
	"strconv"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
)

const (
	BaseUrl = "https://btc-e.com/tapi"
	TradeFee = 0.02
)

type GetInfoResponse struct {
	Funds 				map[string]float64	`json:"funds"`
	Rights 				map[string]int 		`json:"rights"`
	TransactionCount 	int					`json:"transaction_count"`
	OpenOrders 			int 				`json:"open_orders"`
	ServerTime 			time.Time
}

type TransHistoryResponse struct {
	Id					int
	Type 				int					`json:"type"`
	Amount 				float64 			`json:"amount"`
	Currency 			string				`json:"currency"`
	Description 		string 				`json:"desc"`
	Status 				int 				`json:"status"`
	Timestamp 			time.Time 			`json:"-"`
}

type TradeHistoryResponse struct {
	Id					int
	Pair 				string				`json:"pair"`
	Type 				string				`json:"type"`
	Amount 				float64 			`json:"amount"`
	Rate 				float64 			`json:"rate"`
	OrderId				int 				`json:"order_id"`
	YourOrder			int 				`json:"is_your_order"`
	Timestamp 			time.Time 			`json:"-"`
}

type TradeResponse struct {
	Received			float64 			`json:"received"`
	Remains 			float64 			`json:"remains"`
	OrderId				int 				`json:"order_id"`
	Funds 				map[string]float64	`json:"funds"`
}

type CancelOrderResponse struct {
	OrderId				int 				`json:"order_id"`
	Funds 				map[string]float64	`json:"funds"`
}

type BtceApi struct {
	url 				string
	key 				string
	secret 				string
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

// parse an api response into an object
func (b *BtceApi) parseResponse(resp *http.Response) (json.RawMessage, error) {
    // read the response
    defer resp.Body.Close()
    bytes, error := ioutil.ReadAll(resp.Body)
    if error != nil {
    	return nil, error
    }

    // unmarshal into a map
    var raw map[string]json.RawMessage
    if error = json.Unmarshal(bytes, &raw); error != nil {
    	return nil, error
    }

    var success int64

    // handle parse errors
    if error = json.Unmarshal(raw["success"], &success); error != nil {
    	return nil, error
    }

    // handle service errors
    if success != 1 {
    	var errorMsg string
    	json.Unmarshal(raw["error"], &errorMsg)
    	return nil, errors.New("Request failed: "+errorMsg)
    }

    return raw["return"], nil
}

// parse a specific column as a timestamp from a raw json message
func (b *BtceApi) parseTime(message json.RawMessage, key string) (time.Time, error) {
	var data map[string]json.RawMessage

	if error := json.Unmarshal(message, &data); error != nil {
		return time.Time{}, error
    }

    var timestamp int64

	if tsError := json.Unmarshal(data[key], &timestamp); tsError != nil {
		return time.Time{}, tsError
    }

	return time.Unix(timestamp, 0), nil
}

// make a call to the btc-e api
func (b *BtceApi) apiCall(method string, params map[string]string) (json.RawMessage, error) {
    client := &http.Client{}
    postData := b.encodePostData(method, params)

    r, error := http.NewRequest("POST", b.url, bytes.NewBufferString(postData))
    if error != nil {
    	return nil, error
    }

    r.Header.Add("Sign", Sign(b.secret, postData))
    r.Header.Add("Key", b.key)
    r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    r.Header.Add("Content-Length", strconv.Itoa(len(postData)))

    resp, error := client.Do(r)
    if error != nil {
    	return nil, error
    }

    return b.parseResponse(resp)
}

// get info about the account
func (b *BtceApi) GetInfo() (GetInfoResponse, error) {
	var message = GetInfoResponse{}

	resp, error := b.apiCall("getInfo", map[string]string{})
	if error != nil {
		return message, error
	}

	if error = json.Unmarshal(resp, &message); error != nil {
		return message, error
    }

    timestamp, tsError := b.parseTime(resp, "server_time")
    if error != tsError {
    	return message, tsError
    }

    message.ServerTime = timestamp
	return message, nil
}

// returns transaction history
func (b *BtceApi) TransHistory(params map[string]string) ([]TransHistoryResponse, error) {
	resp, error := b.apiCall("TransHistory", params)
	if error != nil {
		return nil, error
	}

	var transactions []TransHistoryResponse
	var data map[string]json.RawMessage

	if error = json.Unmarshal(resp, &data); error != nil {
		return nil, error
    }

	for id, row := range data {
		var trans TransHistoryResponse
		if error = json.Unmarshal(row, &trans); error != nil {
			return nil, error
	    }

	    idInt, error := strconv.Atoi(id)
	    if error != nil {
	    	return nil, error
	    }

	    trans.Id = idInt

	    timestamp, tsError := b.parseTime(row, "timestamp")
	    if error != tsError {
	    	return nil, tsError
	    }

	    trans.Timestamp = timestamp
		transactions = append(transactions, trans)
	}

	return transactions, nil
}

// returns trade history
func (b *BtceApi) TradeHistory(params map[string]string) ([]TradeHistoryResponse, error) {
	resp, error := b.apiCall("TradeHistory", params)
	if error != nil {
		return nil, error
	}

	var trades []TradeHistoryResponse
	var data map[string]json.RawMessage

	if error = json.Unmarshal(resp, &data); error != nil {
		return nil, error
    }

	for id, row := range data {
		var trade TradeHistoryResponse
		if error = json.Unmarshal(row, &trade); error != nil {
			return nil, error
	    }

	    idInt, error := strconv.Atoi(id)
	    if error != nil {
	    	return nil, error
	    }

	    trade.Id = idInt

	    timestamp, tsError := b.parseTime(row, "timestamp")
	    if error != tsError {
	    	return nil, tsError
	    }

	    trade.Timestamp = timestamp
		trades = append(trades, trade)
	}

	return trades, nil
}

// performs a trade
func (b *BtceApi) Trade(pair string, t string, rate float64, amount float64) (TradeResponse, error) {
	var message = TradeResponse{}
	resp, error := b.apiCall("Trade", map[string]string{
		"pair": pair,
		"type": t,
		"rate": strconv.FormatFloat(rate, 'f', -1, 64),
		"amount": strconv.FormatFloat(amount, 'f', -1, 64),
	})
	if error != nil {
		return message, error
	}
	if error = json.Unmarshal(resp, &message); error != nil {
		return message, error
    }

	return message, nil
}

// cancels an order
func (b *BtceApi) CancelOrder(tradeId int) (CancelOrderResponse, error) {
	var message = CancelOrderResponse{}
	resp, error := b.apiCall("CancelOrder", map[string]string{
		"order_id": strconv.Itoa(tradeId),
	})
	if error != nil {
		return message, error
	}
	if error = json.Unmarshal(resp, &message); error != nil {
		return message, error
    }

	return message, nil
}
