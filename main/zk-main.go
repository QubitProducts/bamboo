package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/samuel/go-zookeeper/zk"

	"bamboo/qzk"
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

func main() {

	servers := []string{"localhost"}

	qzk.
	evts, quit := qzk.ListenToZooKeeper(servers, "/foo")
	go showEvents(evts)
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
	quit <- true
	fmt.Printf("exiting!\n")
}
