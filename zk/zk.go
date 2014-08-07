package qzk

import (
	"time"

	"github.com/samuel/go-zookeeper/zk"

	c "bamboo/configuration"
)

func pollZooKeeper(host []string, path string, evts chan zk.Event, quit chan bool) {
	c, _, err := zk.Connect(host, time.Second)
	if err != nil {
		panic(err)
	}

	_, _, selfCh, err := c.GetW(path)
	_, _, childrenCh, err := c.ChildrenW(path)
	if err != nil {
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

func ListenToZooKeeper(config c.Zookeeper) (chan zk.Event, chan bool) {
	quit := make(chan bool)
	evts := make(chan zk.Event)
	go pollZooKeeper(config.ConnectionString(), config.Path, evts, quit)
	return evts, quit
}
