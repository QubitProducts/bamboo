package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/samuel/go-zookeeper/zk"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"

	"bamboo/api"
	"bamboo/configuration"
	"bamboo/qzk"
	"bamboo/services/cmd"
	"bamboo/services/haproxy"
)

func hello(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", c.URLParams["name"])
}

func haproxyConfigUpdateHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "haproxy updated!")
}

/* HTTP Service */
func main() {

	conf := cmd.GetConfiguration()
	conn, _, err := zk.Connect(conf.DomainMapping.Zookeeper.ConnectionString(), time.Second)

	if err != nil {
		panic(err)
	}

	initServer(conf, conn)
	go listenToZookeeper(conf, conn)
}

func initServer(conf configuration.Configuration, conn *zk.Conn) {

	stateAPI := api.State{Config: conf, Zookeeper: conn}
	domainAPI := api.Domain{Config: conf, Zookeeper: conn}

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

func listenToZookeeper(conf configuration.Configuration, conn *zk.Conn) {
	ch, err := qzk.ListenToConn(conn, conf.DomainMapping.Zookeeper.Path, true)
	if err != nil {
		panic(err)
	}

	for {
		_ = <-ch
		writeErr := haproxy.WriteHAProxyConfig(conf.HAProxy, data)
		if writeErr != nil {
			panic(err)
		}
	}
}
