package api

import (
	"io"
	"net/http"
)

// Status Handler
func HandleStatus(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "OK")
}
