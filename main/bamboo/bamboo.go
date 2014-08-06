package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"

	conf "bamboo/configuration"
	"bamboo/marathon"
)

func hello(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", c.URLParams["name"])
}

// Status Handler
func Status(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "OK")
}

func haproxyConfigUpdateHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "haproxy updated!")
}

func marathonAppsHandler(w http.ResponseWriter, r *http.Request) {
	apps, _ := marathon.Apps(config.Marathon.Endpoint)
	payload, _ := json.Marshal(apps)
	io.WriteString(w, string(payload))
}

// Commandline arguments
var configFilePath string

// shared configuration
var config *conf.Configuration

func init() {
	flag.StringVar(&configFilePath, "config", "config/development.json", "Full path of the configuration JSON file")
}

/* HTTP Service */
func main() {
	// Parsing commandline options
	flag.Parse()

	config = &conf.Configuration{}
	err := config.FromFile(configFilePath)

	if err != nil {
		panic(err)
	}

	fmt.Println("", config.Marathon.Endpoint)

	goji.Get("/status", Status)
	goji.Post("/api/haproxy/update", haproxyConfigUpdateHandler)
	goji.Get("/api/marathon/apps", marathonAppsHandler)
	goji.Serve()
}
