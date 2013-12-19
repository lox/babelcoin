package bitcoincharts

import (
	"testing"
	//"github.com/davecgh/go-spew/spew"
	babel "../../core"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDriverSpec(t *testing.T) {
	Convey("Subject: BitcoinCharts Driver", t, func() {

		Convey(`Creating a driver should work`, func() {
			var driver babel.Exchange
			driver = NewDriver(map[string]interface{}{})

			So(driver, ShouldNotBeNil)
		})
	})
}
