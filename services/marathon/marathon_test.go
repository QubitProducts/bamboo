package marathon

import (
	. "github.com/QubitProducts/bamboo/Godeps/_workspace/src/github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestParseHealthCheckPathTCP(t *testing.T) {
	Convey("#parseHealthCheckPath", t, func() {
		checks := []healthCheck{
			healthCheck{"/", "TCP"},
			healthCheck{"/foobar", "TCP"},
			healthCheck{"", "TCP"},
		}
		Convey("should return no path if all checks are TCP", func() {
			So(parseHealthCheckPath(checks), ShouldEqual, "")
		})
	})
}

func TestParseHealthCheckPathHTTP(t *testing.T) {
	Convey("#parseHealthCheckPath", t, func() {
		checks := []healthCheck{
			healthCheck{"/first", "HTTP"},
			healthCheck{"/", "HTTP"},
			healthCheck{"", "HTTP"},
		}
		Convey("should return the first path if all checks are HTTP", func() {
			So(parseHealthCheckPath(checks), ShouldEqual, "/first")
		})
	})
}

func TestParseHealthCheckPathMixed(t *testing.T) {
	Convey("#parseHealthCheckPath", t, func() {
		checks := []healthCheck{
			healthCheck{"", "TCP"},
			healthCheck{"/path", "HTTP"},
			healthCheck{"/", "HTTP"},
		}
		Convey("should return the first path if some checks are HTTP", func() {
			So(parseHealthCheckPath(checks), ShouldEqual, "/path")
		})
	})
}
