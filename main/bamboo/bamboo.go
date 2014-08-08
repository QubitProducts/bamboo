package main

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"io"
	"net/http"
	"time"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"

	"bamboo/api"
	"bamboo/services/cmd"
)

func hello(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", c.URLParams["name"])
}

func haproxyConfigUpdateHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "haproxy updated!")
}

/* HTTP Service */
func main() {

	iniServer()
}

func iniServer() {

	conf := cmd.GetConfiguration()
	conn, _, err := zk.Connect(conf.DomainMapping.Zookeeper.ConnectionString(), time.Second)

	if err != nil {
		panic(err)
	}

	apiState := api.State{Config: conf, Zookeeper: *conn}
	apiDomain := api.Domain{Config: conf, Zookeeper: *conn}

	goji.Get("/status", api.HandleStatus)

	// All state API
	goji.Get("/api/state", apiState.Get)

	// Domains API
	goji.Get("/api/state/domains", apiDomain.All)
	goji.Post("/api/state/domains", apiDomain.Create)
	goji.Delete("/api/state/domains/:id", apiDomain.Delete)
	goji.Put("/api/state/domains/:id", apiDomain.Put)

	goji.Serve()
}
