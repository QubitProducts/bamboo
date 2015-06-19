package marathon

import (
	. "github.com/QubitProducts/bamboo/Godeps/_workspace/src/github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestParseHealthCheckPathTCP(t *testing.T) {
	Convey("#parseHealthCheckPath", t, func() {
		checks := []HealthCheck{
			HealthCheck{"/", "TCP"},
			HealthCheck{"/foobar", "TCP"},
			HealthCheck{"", "TCP"},
		}
		Convey("should return no path if all checks are TCP", func() {
			So(parseHealthCheckPath(checks), ShouldEqual, "")
		})
	})
}

func TestParseHealthCheckPathHTTP(t *testing.T) {
	Convey("#parseHealthCheckPath", t, func() {
		checks := []HealthCheck{
			HealthCheck{"/first", "HTTP"},
			HealthCheck{"/", "HTTP"},
			HealthCheck{"", "HTTP"},
		}
		Convey("should return the first path if all checks are HTTP", func() {
			So(parseHealthCheckPath(checks), ShouldEqual, "/first")
		})
	})
}

func TestParseHealthCheckPathMixed(t *testing.T) {
	Convey("#parseHealthCheckPath", t, func() {
		checks := []HealthCheck{
			HealthCheck{"", "TCP"},
			HealthCheck{"/path", "HTTP"},
			HealthCheck{"/", "HTTP"},
		}
		Convey("should return the first path if some checks are HTTP", func() {
			So(parseHealthCheckPath(checks), ShouldEqual, "/path")
		})
	})
}
