package api

import (
	"io"
	"net/http"
	"encoding/json"

	"bamboo/services/haproxy"
)

func HandleState(w http.ResponseWriter, r *http.Request) {
	payload, _ := json.Marshal(haproxy.GetTemplateData())
	io.WriteString(w, string(payload))
}
