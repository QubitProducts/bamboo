package main

import (
	"fmt"
  	"io"
  	"net/http"
  	"encoding/json"
  	"github.com/zenazn/goji"
 	"github.com/zenazn/goji/web"
	"./marathon"
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
	apps, _ := marathon.Apps("http://aws.fe-marathon.qutics.com:8080")
	payload, _ := json.Marshal(apps)
	io.WriteString(w, string(payload))
}

/* HTTP Service */
func main() {
	goji.Get("/status", Status)
	goji.Post("/api/haproxy/update", haproxyConfigUpdateHandler)
	goji.Get("/api/marathon/apps", marathonAppsHandler)
	goji.Serve()
}
