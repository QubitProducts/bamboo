package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/samuel/go-zookeeper/zk"

	conf "bamboo/configuration"
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

	evts, quit := qzk.ListenToZooKeeper(config.DomainMapping.Zookeeper, true)
	go showEvents(evts)
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
	quit <- true
	fmt.Printf("exiting!\n")
}
