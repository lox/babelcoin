package btce

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	// "net/http/httputil"
	"net/url"
	"testing"
	// "github.com/davecgh/go-spew/spew"
	. "github.com/smartystreets/goconvey/convey"
)

func TestApiSpec(t *testing.T) {
	Convey("Subject: BTC-e Private API", t, func() {

		json := `{"success":0,"error":"llamas be trippin"}`
		var handler func(r *http.Request, body []byte)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			if r.Header.Get("Key") == "" {
				http.Error(w, `{"success":0,"error":"Missing key header"}`,
					http.StatusInternalServerError)
			} else if r.Header.Get("Sign") == "" {
				http.Error(w, `{"success":0,"error":"Missing sign header"}`,
					http.StatusInternalServerError)
			} else if r.Header.Get("Sign") != Sign("credentials", string(body)) {
				http.Error(w, `{"success":0,"error":"Invalid signature"}`,
					http.StatusInternalServerError)
			} else {
				if handler != nil {
					handler(r, body)
				}
				io.WriteString(w, json)
			}
		}))

		Convey(`GetInfo with valid credentials should work`, func() {
			btce := NewBtceApi(server.URL, "valid", "credentials")
			json = `{"success":1,"return":{
				"transaction_count": 0,
				"open_order": 0,
				"server_time": 1386741355
			}}`

			_, error := btce.GetInfo()
			So(error, ShouldBeNil)
		})

		Convey(`GetInfo with invalid credentials shouldn't work`, func() {
			btce := NewBtceApi(server.URL, "invalid", "llamas")
			_, error := btce.GetInfo()

			So(error, ShouldNotBeNil)
		})

		Convey(`GetInfo should return funds`, func() {
			btce := NewBtceApi(server.URL, "valid", "credentials")
			json = `{"success":1,"return":{
				"funds": {"usd": 101},
				"transaction_count": 0,
				"open_order": 0,
				"server_time": 1342448420
			}}`

			info, error := btce.GetInfo()

			So(error, ShouldBeNil)
			So(info.Funds["usd"], ShouldEqual, 101)
		})

		Convey(`TransHistory should return transactions`, func() {
			btce := NewBtceApi(server.URL, "valid", "credentials")
			json = `{ "success":1,
				"return":{
					"1081672":{
						"type":1,
						"amount":1.00000000,
						"currency":"BTC",
						"desc":"BTC Payment",
						"status":2,
						"timestamp":1342448420
					}}}`

			transactions, error := btce.TransHistory(map[string]string{})
			So(error, ShouldBeNil)
			So(len(transactions), ShouldEqual, 1)
			So(transactions[0].Id, ShouldEqual, 1081672)
			So(transactions[0].Amount, ShouldEqual, 1)
			So(transactions[0].Description, ShouldEqual, "BTC Payment")
			So(transactions[0].Timestamp.Unix(), ShouldEqual, 1342448420)
		})

		Convey(`TradeHistory should return transactions`, func() {
			btce := NewBtceApi(server.URL, "valid", "credentials")
			json = `{"success":1,
				"return":{
					"166830":{
						"pair":"btc_usd",
						"type":"sell",
						"amount":1,
						"rate":1,
						"order_id":343148,
						"is_your_order":1,
						"timestamp":1342445793
					}}}`

			trades, error := btce.TradeHistory(map[string]string{})
			So(error, ShouldBeNil)
			So(len(trades), ShouldEqual, 1)
		})

		Convey(`ActiveOrders should return orders`, func() {
			btce := NewBtceApi(server.URL, "valid", "credentials")
			json = `{"success":1,
				"return":{
					"343152":{
						"pair":"btc_usd",
						"type":"sell",
						"amount":1.00000000,
						"rate":3.00000000,
						"timestamp_created":1342448420,
						"status":0
					}}}`

			trades, error := btce.ActiveOrders("btc_usd")

			So(error, ShouldBeNil)
			So(len(trades), ShouldEqual, 1)
		})

		Convey(`Trade should return funds`, func() {
			handler = func(r *http.Request, body []byte) {
				values, err := url.ParseQuery(string(body))
				So(err, ShouldBeNil)
				So(values.Get("rate"), ShouldEqual, "1")
				So(values.Get("amount"), ShouldEqual, "10.5")
				So(values.Get("pair"), ShouldEqual, "btc_usd")
				So(values.Get("type"), ShouldEqual, "buy")
			}

			btce := NewBtceApi(server.URL, "valid", "credentials")
			json = `{ "success":1,
				"return":{
					"received":0.1,
					"remains":0,
					"order_id":10,
					"funds":{
						"usd":325,
						"btc":2.498,
						"sc":121.998,
						"ltc":0,
						"ruc":0,
						"nmc":0
				}}}`

			trade, error := btce.Trade("btc_usd", "buy", 1.0, 10.5)

			So(error, ShouldBeNil)
			So(trade.OrderId, ShouldEqual, 10)
			So(trade.Received, ShouldEqual, 0.1)
			So(trade.Remains, ShouldEqual, 0)
		})

		Convey(`CancelOrder should return funds`, func() {
			handler = func(r *http.Request, body []byte) {}

			btce := NewBtceApi(server.URL, "valid", "credentials")
			json = `{ "success":1,
				"return":{
					"order_id":343154,
					"funds":{
						"usd":325,
						"btc":24.998,
						"sc":121.998,
						"ltc":0,
						"ruc":0,
						"nmc":0
					}}}`

			result, error := btce.CancelOrder(343154)

			So(error, ShouldBeNil)
			So(result.OrderId, ShouldEqual, 343154)
		})
	})
}
