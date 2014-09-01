package qzk

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/samuel/go-zookeeper/zk"

	c "github.com/QubitProducts/bamboo/configuration"
)

var logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

func pollZooKeeper(conn *zk.Conn, path string, evts chan zk.Event, quit chan bool) {

	children, _, err := conn.Children(path)
	if err != nil {
		fmt.Printf("%+v", err)
		panic(err)
	}

	watcherControl := make([]chan<- bool, len(children)+2)
	watcherControl[0] = sinkSelfEvents(conn, path, evts)
	watcherControl[1] = sinkChildEvents(conn, path, evts)

	for i, child := range children {
		p := path + "/" + child
		watcherControl[i+2] = sinkSelfEvents(conn, p, evts)
	}

	<-quit

	for _, ch := range watcherControl {
		ch <- true
	}
}

func debounce(ch chan zk.Event, delay time.Duration) chan zk.Event {
	debounced := make(chan zk.Event)
	var t time.Timer
	var latest zk.Event
	go func() {
		latest = <-ch
		logger.Println("Got event. Delaying post")
		t = *time.NewTimer(delay)
		for {
			select {
			case ev := <-ch:
				logger.Println("Got event. Delaying post")
				latest = ev
				t.Reset(delay)
			case _ = <-t.C:
				logger.Println("No further debouncing. Posting event")
				debounced <- latest
			}
		}
	}()

	return debounced
}

func sinkSelfEvents(conn *zk.Conn, path string, sink chan<- zk.Event) chan<- bool {
	control := make(chan bool)
	go func() {
		_, _, selfCh, err := conn.GetW(path)
		if err != nil {
			logger.Panicf("failed to set listener on path: %s", err.Error())
		}
		for {
			select {
			case _ = <-control:
				break
			case ev := <-selfCh:
				sink <- ev
				_, _, selfCh, err = conn.GetW(path)
				if err != nil {
					logger.Printf("failed to set listener on path: %s\n", err.Error())
				}
			}
		}
	}()

	return control
}

func sinkChildEvents(conn *zk.Conn, path string, sink chan<- zk.Event) chan<- bool {
	control := make(chan bool)
	go func() {
		_, _, selfCh, err := conn.ChildrenW(path)
		if err != nil {
			logger.Panic("failed to set listener on path")
		}
		for {
			select {
			case _ = <-control:
				break
			case ev := <-selfCh:
				sink <- ev
				_, _, selfCh, err = conn.ChildrenW(path)
				if err != nil {
					logger.Panic("failed to set listener on path")
				}
			}
		}
	}()

	return control
}
func ListenToZooKeeper(config c.Zookeeper, deb bool) (chan zk.Event, chan bool) {
	c, _, err := zk.Connect(config.ConnectionString(), time.Second)

	if err != nil {
		panic(err)
	}

	return ListenToConn(c, config.Path, deb)
}

func ListenToConn(c *zk.Conn, path string, deb bool) (chan zk.Event, chan bool) {

	quit := make(chan bool)
	evts := make(chan zk.Event)

	go pollZooKeeper(c, path, evts, quit)

	if deb {
		evts = debounce(evts, 100*time.Millisecond)
	}
	return evts, quit
}
