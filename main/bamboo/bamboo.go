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

	state := api.State{Config: conf, Zookeeper: *conn}

	goji.Get("/status", api.HandleStatus)

	goji.Get("/api/state", state.Get)

	goji.Post("/api/haproxy/update", haproxyConfigUpdateHandler)
	goji.Serve()
}
