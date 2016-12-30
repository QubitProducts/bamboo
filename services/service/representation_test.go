package service

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"

	. "github.com/smartystreets/goconvey/convey"
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

func TestV1ServiceRepr(t *testing.T) {
	Convey("#ParseV1ServiceRepr", t, func() {
		Convey("when we parse an arbitrary string", func() {
			str := randString(128)
			path := randString(8)

			repr, err := ParseV1ServiceRepr([]byte(str), path)

			Convey("it should not error", func() {
				So(err, ShouldBeNil)
			})

			Convey("it should be a V1ServiceRepr", func() {
				So(repr, ShouldHaveSameTypeAs, &V1ServiceRepr{})
			})

			Convey("when we create a service from the repr", func() {
				service := repr.Service()

				Convey("it should have the same id as the passed path", func() {
					So(service.Id, ShouldEqual, path)
				})

				Convey("it should have the same ACL as the passed body", func() {
					So(service.Acl, ShouldEqual, str)
				})

				Convey("it should have no other config entries", func() {
					acl, ok := service.Config["Acl"]
					So(ok, ShouldBeTrue)
					So(acl, ShouldNotBeBlank)
					So(len(service.Config), ShouldEqual, 1)
				})
			})

		})
	})
}

func TestV2ServiceRepr(t *testing.T) {
	Convey("#ParseV1ServiceRepr", t, func() {
		Convey("when we parse an invalid string", func() {
			str := "}{" + randString(128)
			path := randString(8)

			_, err := ParseV2ServiceRepr([]byte(str), path)

			Convey("it should error", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("when we parse a json blob with an incorrect version", func() {
			path := randString(8)
			body := []byte(`{"version": "3", "config": {}}`)
			_, err := ParseV2ServiceRepr(body, path)

			Convey("it should error", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("when we parse a json blob with a correct version and no config", func() {
			path := randString(8)
			body := []byte(`{"version": "2", "config": {}}`)
			repr, err := ParseV2ServiceRepr(body, path)

			Convey("it should not error", func() {
				So(err, ShouldBeNil)
			})

			Convey("when we create a service from the repr", func() {
				service := repr.Service()

				Convey("it should have the same id as the passed path", func() {
					So(service.Id, ShouldEqual, path)
				})

				Convey("it should have no ACL", func() {
					So(service.Acl, ShouldBeBlank)
				})

				Convey("it should have no config value", func() {
					So(len(service.Config), ShouldEqual, 0)
				})
			})
		})

		Convey("when we parse a json blob with a correct version and a complete config", func() {
			path := randString(8)
			body := []byte(`{"version": "2", "config": {"Acl": "foo", "arb": "barb"}}`)
			repr, err := ParseV2ServiceRepr(body, path)

			Convey("it should not error", func() {
				So(err, ShouldBeNil)
			})

			Convey("when we create a service from the repr", func() {
				service := repr.Service()

				Convey("it should have the same id as the passed path", func() {
					So(service.Id, ShouldEqual, path)
				})

				Convey("it should have an ACL", func() {
					So(service.Acl, ShouldEqual, "foo")
				})

				Convey("it should have an additional config value", func() {
					So(len(service.Config), ShouldEqual, 2)
					So(service.Config["arb"], ShouldEqual, "barb")
				})
			})
		})
	})

	Convey("#MakeV2ServiceRepr", t, func() {
		Convey("MakeV2ServiceRepr and V2ServiceRepr.Service should compose to identity", func() {
			property := func(s Service) bool {
				// The generator may create services with different .Acl and .Conf["Acl"]
				s.Config["Acl"] = s.Acl
				s2 := MakeV2ServiceRepr(s).Service()
				return reflect.DeepEqual(s, s2)
			}

			err := quick.Check(property, nil)
			So(err, ShouldBeNil)
		})
	})
}

func TestParseServiceRepr(t *testing.T) {
	Convey("#ParseServiceRepr", t, func() {
		Convey("when we parse an invalid json string", func() {
			path := randString(8)
			body := []byte("}{" + randString(128))

			repr, err := ParseServiceRepr([]byte(body), path)

			Convey("it should not error", func() {
				So(err, ShouldBeNil)
			})

			Convey("it should be a V1ServiceRepr", func() {
				So(repr, ShouldHaveSameTypeAs, &V1ServiceRepr{})
			})

			Convey("when we create a service from the repr", func() {
				service := repr.Service()

				Convey("it should ave the same id as the passed path", func() {
					So(service.Id, ShouldEqual, path)
				})

				Convey("its acl should be the body", func() {
					So(service.Acl, ShouldEqual, string(body))
				})

				Convey("it should have no other config values", func() {
					So(len(service.Config), ShouldEqual, 1)
				})
			})
		})

		Convey("when we parse a v2 repr", func() {
			path := randString(8)
			body := []byte(`{"version": "2", "config": {"Acl": "foo", "arb": "barb"}}`)

			repr, err := ParseServiceRepr(body, path)

			Convey("it should not error", func() {
				So(err, ShouldBeNil)
			})

			Convey("it should be a V2ServiceRepr", func() {
				So(repr, ShouldHaveSameTypeAs, &V2ServiceRepr{})
			})

			Convey("when we create a service from the repr", func() {
				service := repr.Service()

				Convey("it should have the same id as the passed path", func() {
					So(service.Id, ShouldEqual, path)
				})

				Convey("its acl should be foo", func() {
					So(service.Acl, ShouldEqual, "foo")
				})

				Convey("it should have one other config value", func() {
					So(len(service.Config), ShouldEqual, 2)
					So(service.Config["arb"], ShouldEqual, "barb")
				})
			})
		})
	})
}
