package main

import (
	"fmt"
	"time"

	"github.com/samuel/go-zookeeper/zk"
	"github.com/zenazn/goji"

	"bamboo/api"
	"bamboo/configuration"
	"bamboo/qzk"
	"bamboo/services/cmd"
	"bamboo/services/haproxy"
)

/* HTTP Service */
func main() {
	conf := cmd.GetConfiguration()
	conns := listenToZookeeper(conf)

	initServer(conf, conns)
}

func initServer(conf configuration.Configuration, conns Conns) {

	stateAPI := api.State{Config: conf, Zookeeper: conns.DomainMapping}
	domainAPI := api.Domain{Config: conf, Zookeeper: conns.DomainMapping}

	// Status live information
	goji.Get("/status", api.HandleStatus)

	// All state API
	goji.Get("/api/state", stateAPI.Get)

	// Domains API
	goji.Get("/api/state/domains", domainAPI.All)
	goji.Post("/api/state/domains", domainAPI.Create)
	goji.Delete("/api/state/domains/:id", domainAPI.Delete)
	goji.Put("/api/state/domains/:id", domainAPI.Put)

	goji.Serve()
}

type Conns struct {
	Marathon      *zk.Conn
	DomainMapping *zk.Conn
}

func listenToZookeeper(conf configuration.Configuration) Conns {

	marathonCh, marathonConn := createAndListen(conf.Marathon.Zookeeper)
	domainCh, domainConn := createAndListen(conf.DomainMapping.Zookeeper)

	go func() {
		for {
			select {
			case _ = <-marathonCh:
				fmt.Println("Marathon state changed")
				handleHAPUpdate(conf, marathonConn)
			case _ = <-domainCh:
				fmt.Println("Domain mapping stated changed")
				handleHAPUpdate(conf, marathonConn)
			}
		}
	}()

	return Conns{marathonConn, domainConn}
}

func handleHAPUpdate(conf configuration.Configuration, conn * zk.Conn) {
	err := haproxy.WriteHAProxyConfig(conf.HAProxy, haproxy.GetTemplateData(conf, conn))
	if err != nil {
		panic(err)
	}
	fmt.Println("HAProxy cfg Updated")
}

func createAndListen(conf configuration.Zookeeper) (chan zk.Event, *zk.Conn) {
	conn, _, err := zk.Connect(conf.ConnectionString(), time.Second)

	if err != nil {
		panic(err)
	}

	ch, _ := qzk.ListenToConn(conn, conf.Path, true)

	return ch, conn
}
