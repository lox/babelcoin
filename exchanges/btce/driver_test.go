package btce

import (
	babel "../../core"
	"github.com/davecgh/go-spew/spew"
	"io"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDriverSpec(t *testing.T) {
	Convey("Subject: BTC-e Driver", t, func() {

		json := make(chan string, 25)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			select {
			case body := <-json:
				spew.Printf("%s\n", body)
				io.WriteString(w, body)
			case <-time.After(1000):
				panic("timed out waiting for json")
			}
		}))

		Convey(`Creating a limit order should work`, func() {
			var driver babel.Exchange
			driver, err := NewDriver(map[string]interface{}{
				"key":         "correct",
				"secret":      "credentials",
				"info_url":    server.URL,
				"ticker_url":  server.URL,
				"private_url": server.URL,
			})

			So(err, ShouldBeNil)
			So(driver, ShouldNotBeNil)

			order := driver.Buy("btc_usd", 100, 0.590)

			Convey(`Order should have fees`, func() {
				json <- `{
					"server_time": 1370814956,
					"pairs": {
						"btc_usd": {
							"decimal_places": 3,
							"min_price": 0.1,
							"max_price": 400,
							"min_amount": 0.01,
							"hidden": 0,
							"fee": 0.0234
					}}}`

				So(order, ShouldNotBeNil)

				fee, err := order.Fee()

				So(err, ShouldBeNil)
				So(fee, ShouldEqual, 0.0234)
			})

			Convey(`Order execution should return a channel with trades`, func() {
				json <- `{
						"success":1,
						"return":{
							"received":0.1,
							"remains":0,
							"order_id":10024,
							"funds":{
								"usd":325,
								"btc":2.498
						}}}`

				json <- `{
						"success":1,
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

				trades, err := order.Execute()
				results := []babel.Trade{}

				for trade := range trades {
					results = append(results, trade)
				}

				So(err, ShouldBeNil)
				So(len(results), ShouldEqual, 1)
			})

		})
	})
}
