package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	conf "github.com/QubitProducts/bamboo/configuration"
	"github.com/QubitProducts/bamboo/qzk"
	"github.com/samuel/go-zookeeper/zk"
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

var configFilePath string
var config conf.Configuration

func init() {
	flag.StringVar(&configFilePath, "config", "config/development.json", "Full path of the configuration JSON file")
}

func main() {

	config, err := conf.FromFile(configFilePath)

	if err != nil {
		panic(err)
	}

	evts, quit := qzk.ListenToZooKeeper(config.Bamboo.Zookeeper, true)
	go showEvents(evts)
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
	quit <- true
	fmt.Printf("exiting!\n")
}
