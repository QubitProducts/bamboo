package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"

	"bamboo/api"
)

func hello(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", c.URLParams["name"])
}

func haproxyConfigUpdateHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "haproxy updated!")
}

/* HTTP Service */
func main() {
	// Parsing commandline options
	goji.Get("/status", api.HandleStatus)

	goji.Get("/api/state", api.HandleState)

	goji.Post("/api/haproxy/update", haproxyConfigUpdateHandler)
	goji.Serve()
}



