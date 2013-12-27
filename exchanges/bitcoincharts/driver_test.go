package bitcoincharts

import (
	"testing"
	//"github.com/davecgh/go-spew/spew"
	babel "github.com/lox/babelcoin/core"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDriverSpec(t *testing.T) {
	Convey("Subject: BitcoinCharts Driver", t, func() {

		Convey(`Creating a driver should work`, func() {
			var driver babel.Exchange
			driver = New("bitcoincharts-mtgox", map[string]interface{}{})

			So(driver, ShouldNotBeNil)
		})
	})
}
