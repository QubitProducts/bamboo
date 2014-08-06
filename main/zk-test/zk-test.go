package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/samuel/go-zookeeper/zk"

	conf "bamboo/configuration"
	"bamboo/zk"
)

func showEvents(evts chan zk.Event) {
	ok := true
	for ok {
		ev, isOk := <-evts
		ok = isOk
		if ev.State != zk.StateDisconnected {
			fmt.Printf("a change occurred: %+v\n", ev)
		}
	}
}

var config *conf.Configuration

func init() {
}

func main() {

	evts, quit := qzk.ListenToZooKeeper(config.ServicesMapping.Zookeeper)
	go showEvents(evts)
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
	quit <- true
	fmt.Printf("exiting!\n")
}
