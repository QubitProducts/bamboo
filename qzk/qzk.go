package qzk

import (
	"fmt"
	"time"

	"github.com/samuel/go-zookeeper/zk"

	c "bamboo/configuration"
)

const (
	ErrNoSuchPath = 1
)

type Err int32

func pollZooKeeper(c zk.Conn, path string, evts chan zk.Event, quit chan bool) {
	_, _, selfCh, err := c.GetW(path)
	_, _, childrenCh, err := c.ChildrenW(path)
	if err != nil {
		fmt.Printf("%+v", err)
		panic(err)
	}

	keepGoing := true
	for keepGoing {
		select {
		case selfEv := <-selfCh:
			evts <- selfEv
			_, _, selfCh, _ = c.GetW(path)
		case childEv := <-childrenCh:
			evts <- childEv
			_, _, childrenCh, _ = c.ChildrenW(path)
		case _ = <-quit:
			close(evts)
			keepGoing = false
		}
	}

}

func debounce(ch chan zk.Event, delay time.Duration) chan zk.Event {
	debounced := make(chan zk.Event)
	var t time.Timer
	var latest zk.Event
	go func() {
		latest = <-ch
		t = *time.NewTimer(delay)
		for {
			select {
			case ev := <-ch:
				latest = ev
				t.Reset(delay)
			case _ = <-t.C:
				debounced <- latest
			}
		}
	}()

	return debounced
}

func ListenToZooKeeper(config c.Zookeeper, deb bool) (chan zk.Event, chan bool) {
	c, _, err := zk.Connect(config.ConnectionString(), time.Second)

	if err != nil {
		panic(err)
	}

	return ListenToConn(*c, config.Path, deb)
}

func ListenToConn(c zk.Conn, path string, deb bool) (chan zk.Event, chan bool) {

	quit := make(chan bool)
	evts := make(chan zk.Event)

	go pollZooKeeper(c, path, evts, quit)

	if deb {
		evts = debounce(evts, 100*time.Millisecond)
	}
	return evts, quit
}
