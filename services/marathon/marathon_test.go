package marathon

import (
	. "github.com/QubitProducts/bamboo/Godeps/_workspace/src/github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetMesosDnsId_Simple(t *testing.T) {
	Convey("#getMesosDnsId", t, func() {
		Convey("should return simple appname", func() {
			So(getMesosDnsId("appname"), ShouldEqual, "appname")
		})

		Convey("should return simple appname if slash prefixed", func() {
			So(getMesosDnsId("/appname"), ShouldEqual, "appname")
		})

		Convey("should return groups reverse-added to appname", func() {
			So(getMesosDnsId("/group/appname"), ShouldEqual, "appname-group")
		})

		Convey("should return groups reverse-added to appname but no blanks", func() {
			So(getMesosDnsId("//group/again//appname/"), ShouldEqual, "appname-again-group")
		})
	})
}

func TestParseHealthCheckPathTCP(t *testing.T) {
	Convey("#parseHealthCheckPath", t, func() {
		checks := []marathonHealthCheck{
			marathonHealthCheck{"/", "TCP", 0},
			marathonHealthCheck{"/foobar", "TCP", 0},
			marathonHealthCheck{"", "TCP", 0},
		}
		Convey("should return no path if all checks are TCP", func() {
			So(parseHealthCheckPath(checks), ShouldEqual, "")
		})
	})
}

func TestParseHealthCheckPathHTTP(t *testing.T) {
	Convey("#parseHealthCheckPath", t, func() {
		checks := []marathonHealthCheck{
			marathonHealthCheck{"/first", "HTTP", 0},
			marathonHealthCheck{"/", "HTTP", 0},
			marathonHealthCheck{"", "HTTP", 0},
		}
		Convey("should return the first path if all checks are HTTP", func() {
			So(parseHealthCheckPath(checks), ShouldEqual, "/first")
		})
	})
}

func TestParseHealthCheckPathMixed(t *testing.T) {
	Convey("#parseHealthCheckPath", t, func() {
		checks := []marathonHealthCheck{
			marathonHealthCheck{"", "TCP", 0},
			marathonHealthCheck{"/path", "HTTP", 0},
			marathonHealthCheck{"/", "HTTP", 0},
		}
		Convey("should return the first path if some checks are HTTP", func() {
			So(parseHealthCheckPath(checks), ShouldEqual, "/path")
		})
	})
}
