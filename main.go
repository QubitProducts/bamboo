package main

import (
	"fmt"
    "io"
    "net/http"
    "github.com/zenazn/goji"
    "github.com/zenazn/goji/web"
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

//
/* HTTP Service */
func main() {
	goji.Get("/status", Status)
	goji.Post("/api/haproxy/update", haproxyConfigUpdateHandler)
	goji.Serve()
}
