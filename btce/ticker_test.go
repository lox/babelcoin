package btce

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	//"github.com/davecgh/go-spew/spew"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTicker(t *testing.T) {
	Convey("Subject: BTC-e Ticker API", t, func() {

		json := `{"success":0,"error":"llamas be trippin"}`

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, json)
		}))

		Convey(`Ticker should return data`, func() {
			ticker := NewTickerApi(server.URL+"/", []string{"btc_usd"})
			json = `{ "btc_usd": {
					"High": 109.88,
					"Low": 91.14,
					"Avg": 100.51,
					"Vol": 1632898.2249,
					"Vol_cur": 16541.51969,
					"Last": 101.773,
					"Buy": 101.9,
					"Sell": 101.773,
					"Updated": 1370816308
				}
			}`

			data, error := ticker.MarketData()

			So(error, ShouldBeNil)
			So(len(data), ShouldEqual, 1)
			So(data[0].Pair, ShouldEqual, "btc_usd")
		})

		Convey(`Ticker should fail with invalid pairs`, func() {
			ticker := NewTickerApi(server.URL+"/", []string{"aaa_bbb"})
			json = `{"success":0, "error":"Invalid pair name: aaa_bbb"}`

			data, error := ticker.MarketData()

			So(error, ShouldNotBeNil)
			So(error.Error(), ShouldEqual, "Invalid pair name: aaa_bbb")
			So(len(data), ShouldEqual, 0)
		})

	})
}
