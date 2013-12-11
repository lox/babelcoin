package btce

import(
	"net/http"
	"net/http/httptest"
	"testing"
	"io"
	//"github.com/davecgh/go-spew/spew"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTicker(t *testing.T) {
	Convey("Subject: BTC-e Ticker API", t, func() {

		json := `{"success":0,"error":"llamas be trippin"}`

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, json)
		}))

		Convey(`Ticker should return candles`, func() {
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

			candles, error := ticker.Candles()

			So(error, ShouldBeNil)
			So(len(candles) , ShouldEqual, 1)
			So(candles[0].Pair, ShouldEqual, "btc_usd")
		})
	})
}

