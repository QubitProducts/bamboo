package service

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"

	"log"
	"time"

	"github.com/QubitProducts/bamboo/configuration"
	"github.com/samuel/go-zookeeper/zk"
)

var zkTimeout = time.Second
var zkConf = configuration.Zookeeper{
	Host:           "localhost:2181",
	Path:           "/test-bamboo",
	ReportingDelay: 1,
}

func cleanZK(conn *zk.Conn) {
	deleteRecursive(conn, zkConf.Path)
}

func deleteRecursive(conn *zk.Conn, path string) {
	children, _, err := conn.Children(path)
	if err != nil {
		log.Printf("failed to get children %s: %s", path, err)
		return
	}
	for _, child := range children {
		deleteRecursive(conn, path+"/"+child)
	}
	err = conn.Delete(path, -1)
	if err != nil {
		log.Printf("failed to delete %s: %s", path, err)
	}
}

func loadToZK(conn *zk.Conn, data [][2]string) {
	for _, entry := range data {
		k := entry[0]
		v := entry[1]
		_, err := conn.Create(k, []byte(v), 0, defaultACL())
		orPanic(err)
	}
}

func TestZKStorage(t *testing.T) {
	conn, _, err := zk.Connect(zkConf.ConnectionString(), zkTimeout)
	orPanic(err)

	Convey("#NewZKStorage", t, func() {
		cleanZK(conn)
		s, err := NewZKStorage(conn, zkConf)

		So(err, ShouldBeNil)

		Convey("it should implement the Storage interface", func() {
			_, ok := interface{}(s).(Storage)
			So(ok, ShouldBeTrue)
		})

		Convey("it should have created the base path", func() {
			exists, _, err := conn.Exists(zkConf.Path)

			So(err, ShouldBeNil)
			So(exists, ShouldBeTrue)
		})
	})

	Convey("#ZKStorage.All", t, func() {
		s, err := NewZKStorage(conn, zkConf)
		So(err, ShouldBeNil)

		Convey("when I get all in an empty zookeeper", func() {
			cleanZK(conn)
			entries, err := s.All()

			So(err, ShouldBeNil)

			Convey("there should be no entries", func() {
				So(len(entries), ShouldEqual, 0)
			})
		})

		Convey("when I get all in a legacy/v1 zookeeper", func() {
			cleanZK(conn)
			loadToZK(conn, [][2]string{
				[2]string{zkConf.Path, ""},
				[2]string{zkConf.Path + "/test", "hdr(host) -i foo"},
				[2]string{zkConf.Path + "/test2", "fozbaz"},
			})

			entries, err := s.All()

			So(err, ShouldBeNil)

			Convey("there should be the correct number of entries", func() {
				So(len(entries), ShouldEqual, 2)
			})

			Convey("there should be an entry with the id test", func() {
				var test Service
				found := false
				for _, i := range entries {
					if i.Id == "test" {
						test = i
						found = true
						break
					}
				}
				So(found, ShouldBeTrue)

				Convey("which should have an acl 'hdr(host) -i foo'", func() {
					So(test.Acl, ShouldEqual, "hdr(host) -i foo")
					So(test.Config["Acl"], ShouldEqual, "hdr(host) -i foo")
				})
			})
		})

		Convey("when I get all in a mixed v1/v2 zookeeper", func() {
			cleanZK(conn)
			loadToZK(conn, [][2]string{
				[2]string{zkConf.Path, ""},
				[2]string{zkConf.Path + "/test", `{"version": "2", "config": {"Acl": "foo", "arb": "barb"}}`},
				[2]string{zkConf.Path + "/test2", "fozbaz"},
			})

			entries, err := s.All()

			So(err, ShouldBeNil)

			Convey("there should be the correct number of entries", func() {
				So(len(entries), ShouldEqual, 2)
			})

			Convey("there should be an entry with the id test", func() {
				var test Service
				found := false
				for _, i := range entries {
					if i.Id == "test" {
						test = i
						found = true
						break
					}
				}
				So(found, ShouldBeTrue)

				Convey("which should have an acl 'foo'", func() {
					So(test.Acl, ShouldEqual, "foo")
					So(test.Config["Acl"], ShouldEqual, "foo")
				})

				Convey("which should have a config entry 'arb'", func() {
					So(test.Config["arb"], ShouldEqual, "barb")
				})
			})
		})
	})

	Convey("#ZKStorage.Upsert", t, func() {
		s, err := NewZKStorage(conn, zkConf)
		So(err, ShouldBeNil)

		testService := Service{
			Id: "test",
			Config: map[string]string{
				"Acl":   "foo",
				"barst": "carst",
			},
		}

		Convey("when I insert into an empty key", func() {
			cleanZK(conn)
			err := s.Upsert(testService)
			So(err, ShouldBeNil)

			entries, err := s.All()
			So(err, ShouldBeNil)

			Convey("there should be 1 entry", func() {
				So(len(entries), ShouldEqual, 1)

				readEntry := entries[0]

				Convey("which should have an Acl 'foo'", func() {
					So(readEntry.Acl, ShouldEqual, "foo")
					So(readEntry.Config["Acl"], ShouldEqual, "foo")
				})

				Convey("which should have a config entry 'barst'", func() {
					So(readEntry.Config["barst"], ShouldEqual, "carst")
				})
			})
		})

		Convey("when I insert into an existing key", func() {
			cleanZK(conn)
			loadToZK(conn, [][2]string{
				[2]string{zkConf.Path, ""},
				[2]string{zkConf.Path + "/test", "fozbaz"},
			})

			err := s.Upsert(testService)
			So(err, ShouldBeNil)

			entries, err := s.All()
			So(err, ShouldBeNil)

			Convey("there should be 1 entry", func() {
				So(len(entries), ShouldEqual, 1)

				readEntry := entries[0]

				Convey("which should have an Acl 'foo'", func() {
					So(readEntry.Acl, ShouldEqual, "foo")
					So(readEntry.Config["Acl"], ShouldEqual, "foo")
				})

				Convey("which should have a config entry 'barst'", func() {
					So(readEntry.Config["barst"], ShouldEqual, "carst")
				})
			})
		})
	})

	Convey("#ZKStorage.Delete", t, func() {
		s, err := NewZKStorage(conn, zkConf)
		So(err, ShouldBeNil)

		Convey("when I delete the only service", func() {
			cleanZK(conn)
			loadToZK(conn, [][2]string{
				[2]string{zkConf.Path, ""},
				[2]string{zkConf.Path + "/test", "fozbaz"},
			})

			err := s.Delete("test")
			So(err, ShouldBeNil)

			Convey("there should be zero entries", func() {
				entries, err := s.All()
				So(err, ShouldBeNil)

				So(len(entries), ShouldEqual, 0)
			})
		})

		Convey("when I delete an non-existant service", func() {
			cleanZK(conn)

			err := s.Delete("test")

			Convey("it should error", func() {
				So(err, ShouldNotBeNil)
			})
			So(err, ShouldNotBeNil)

			Convey("there should be zero entries", func() {
				entries, err := s.All()
				So(err, ShouldBeNil)

				So(len(entries), ShouldEqual, 0)
			})
		})
	})
}
