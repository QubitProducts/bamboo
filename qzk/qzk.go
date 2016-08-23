package qzk

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	c "github.com/QubitProducts/bamboo/configuration"
	"github.com/samuel/go-zookeeper/zk"
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

func delay(ch chan zk.Event, delay time.Duration) chan zk.Event {
	delayed := make(chan zk.Event)
	go func() {
		var ev zk.Event
		for {
			ev = <-ch
			time.AfterFunc(delay, func() {
				delayed <- ev
			})
		}
	}()

	return delayed
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

	return ListenToConn(c, config.Path, deb, config.Delay())
}

func zkNodeCreateByPath(path string, c *zk.Conn) error {
	return zkCreateNodes("", strings.Split(path, "/"), c)
}

func zkCreateNodes(path string, nodes []string, c *zk.Conn) error {
	if len(nodes) > 0 {
		// strings.Split will return empty-strings for leading split chars, lets skip over these.
		if len(nodes[0]) == 0 {
			return zkCreateNodes(path, nodes[1:], c)
		}
		fqPath := path + "/" + nodes[0]
		log.Printf("Creating path: %v", fqPath)
		exists, _, err := c.Exists(fqPath)
		if err != nil {
			return err
		}
		if !exists {
			_, err := c.Create(fqPath, []byte{}, 0, zk.WorldACL(zk.PermAll))
			if err != nil {
				return err
			}
		}
		return zkCreateNodes(fqPath, nodes[1:], c)
	}
	return nil
}

func ListenToConn(c *zk.Conn, path string, deb bool, repDelay time.Duration) (chan zk.Event, chan bool) {
	exists, _, err := c.Exists(path)

	if err != nil {
		logger.Fatalf("Couldn't determine whether node %v exists in Zookeeper", path)
	}
	if !exists {
		logger.Printf("Node '%v' does not exist in Zookeeper, creating...", path)
		err := zkNodeCreateByPath(path, c)
		if err != nil {
			logger.Fatalf("Unable to create path '%v': %v", path, err)
		}

	}

	quit := make(chan bool)
	evts := make(chan zk.Event)

	go pollZooKeeper(c, path, evts, quit)

	if deb {
		evts = debounce(evts, 100*time.Millisecond)
	}
	if repDelay > 0 {
		evts = delay(evts, repDelay)
	}
	return evts, quit
}

func nodeExists(c *zk.Conn, path string) bool {
	return false
}
