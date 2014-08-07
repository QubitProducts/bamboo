package api

import (
	"io"
	"fmt"
	"encoding/json"

	"bamboo/services/cmd"
//	"bamboo/services/domain"
	"bamboo/services/marathon"
	"net/http"

)

func HandleState(w http.ResponseWriter, r *http.Request) {
	config := cmd.GetConfiguration()
	fmt.Printf("configpath: %v\n", cmd.GetConfigFilePath())
	fmt.Printf("marathon: %v\n", config.HAProxy)

	apps, _ := marathon.FetchApps(config.Marathon.Endpoint)
//	services, _ := domain.FetchAll(config.ServicesMapping.Zookeeper)
	services := map[string]string{}

	data := struct {
		Apps     []marathon.App
		Services map[string]string
	}{apps, services}

	payload, _ := json.Marshal(data)

	io.WriteString(w, string(payload))
}
