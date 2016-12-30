package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/cloverstd/bamboo/configuration"
	"github.com/cloverstd/bamboo/services/haproxy"
	"github.com/cloverstd/bamboo/services/service"
)

type StateAPI struct {
	Config  *configuration.Configuration
	Storage service.Storage
}

func (state *StateAPI) Get(w http.ResponseWriter, r *http.Request) {
	templateData, _ := haproxy.GetTemplateData(state.Config, state.Storage)
	payload, _ := json.Marshal(templateData)
	io.WriteString(w, string(payload))
}
