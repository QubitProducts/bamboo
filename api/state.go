package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/samuel/go-zookeeper/zk"

	"bamboo/configuration"
	"bamboo/services/haproxy"
)

type State struct {
	Config    configuration.Configuration
	Zookeeper zk.Conn
}

func (state State) Get(w http.ResponseWriter, r *http.Request) {
	payload, _ := json.Marshal(haproxy.GetTemplateData(state.Config))
	io.WriteString(w, string(payload))
}
