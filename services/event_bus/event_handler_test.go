package event_bus

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/QubitProducts/bamboo/configuration"
)

// Yanked off http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

func randString(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func orPanic(err error) {
	if err != nil {
		panic(fmt.Sprintf("Error! %s", err))
	}
}

func TestEventHandler(t *testing.T) {
	var config configuration.Configuration

	Convey("#isReloadRequired", t, func() {

		tmpPath := fmt.Sprintf("/tmp/bamboo_irr%v.conf", rand.Int31())
		contents := randString(int(rand.Int31n(1048576)))

		orPanic(ioutil.WriteFile(tmpPath, []byte(contents), 0644))

		Convey("When we test with a file's own contents", func() {

			required, err := isReloadRequired(tmpPath, contents)
			orPanic(err)

			Convey("a reload should not be required", func() {
				So(required, ShouldEqual, false)
			})
		})

		Convey("When we test with changed config", func() {
			required, err := isReloadRequired(tmpPath, contents+"arst")
			orPanic(err)

			Convey("a reload should be required", func() {
				So(required, ShouldEqual, true)
			})
		})

		Convey("When we test against an absent file", func() {
			required, err := isReloadRequired("/tmp/bamboo_irr_nonexistant.conf", contents)
			orPanic(err)

			Convey("a reload should be required", func() {
				So(required, ShouldEqual, true)
			})
		})

		Convey("When we test against an invalid file", func() {
			_, err := isReloadRequired("/dev/null/foobar", contents)

			Convey("we should get an error", func() {
				So(err, ShouldNotEqual, nil)
			})
		})
	})

	Convey("#changeConfig", t, func() {
		Convey("When we change the config with a failing command", func() {
			config.HAProxy.ReloadCommand = "exit 1"
			config.HAProxy.OutputPath = fmt.Sprintf("/tmp/bamboo_irr%v.conf", rand.Int31())
			reloaded, err := changeConfig(&config, "arst")

			Convey("We should get an error", func() {
				So(err, ShouldNotEqual, nil)
			})

			Convey("It should report not reloaded", func() {
				So(reloaded, ShouldEqual, false)
			})
		})

		Convey("When we change the config with a succeeding command", func() {
			config.HAProxy.ReloadCommand = "exit 0"
			config.HAProxy.OutputPath = fmt.Sprintf("/tmp/bamboo_irr%v.conf", rand.Int31())
			reloaded, err := changeConfig(&config, "farst")

			Convey("We should not get an error", func() {
				So(err, ShouldEqual, nil)
			})

			Convey("It should report reloaded", func() {
				So(reloaded, ShouldEqual, true)
			})
		})
	})

	Convey("#validateConfig", t, func() {
		Convey("When we validate the config with a failing command", func() {
			err := validateConfig("exit 1", "arst")

			Convey("The err should be non nil", func() {
				So(err, ShouldNotEqual, nil)
			})
		})

		Convey("When we validate the config with a succeeding command", func() {
			err := validateConfig("exit 0", "arst")

			Convey("The err should be nil", func() {
				So(err, ShouldEqual, nil)
			})
		})

		Convey("When we validate the config with a template", func() {
			err := validateConfig("[ \"x$(cat {{.}})\" = 'xarst' ]", "arst")

			Convey("It should substitute the file", func() {
				So(err, ShouldEqual, nil)
			})
		})
	})
}
