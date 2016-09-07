package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	conf "github.com/QubitProducts/bamboo/configuration"
	"github.com/QubitProducts/bamboo/services/service"
	"github.com/go-martini/martini"
)

type ServiceAPI struct {
	Config  *conf.Configuration
	Storage service.Storage
}

func (d *ServiceAPI) All(w http.ResponseWriter, r *http.Request) {
	services, err := d.Storage.All()

	if err != nil {
		responseError(w, err.Error())
		return
	}

	byId := make(map[string]service.Service, len(services))
	for _, s := range services {
		byId[s.Id] = s
	}

	responseJSON(w, byId)
}

func (d *ServiceAPI) Create(w http.ResponseWriter, r *http.Request) {
	d.updateService(w, r)
}

func (d *ServiceAPI) Put(params martini.Params, w http.ResponseWriter, r *http.Request) {
	d.updateService(w, r)
}

func (d *ServiceAPI) updateService(w http.ResponseWriter, r *http.Request) {
	service, err := extractService(r)
	if err != nil {
		responseError(w, err.Error())
		return
	}

	err = d.Storage.Upsert(service)
	if err != nil {
		responseError(w, err.Error())
		return
	}

	responseJSON(w, service)
}

func (d *ServiceAPI) Delete(params martini.Params, w http.ResponseWriter, r *http.Request) {
	serviceID := params["_1"]
	if len(serviceID) == 0 {
		responseError(w, "can not use empty ID")
		return
	}

	if !strings.HasPrefix(serviceID, "/") {
		serviceID = "/" + serviceID
	}

	err := d.Storage.Delete(serviceID)
	if err != nil {
		responseError(w, err.Error())
		return
	}

	responseJSON(w, new(map[string]string))
}

func extractService(r *http.Request) (service.Service, error) {
	var serviceModel service.Service
	payload, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(payload, &serviceModel)
	if err != nil {
		return serviceModel, errors.New("Unable to decode JSON request")
	}

	if len(serviceModel.Id) == 0 {
		return serviceModel, errors.New("can not use empty ID")
	}

	if !strings.HasPrefix(serviceModel.Id, "/") {
		serviceModel.Id = "/" + serviceModel.Id
	}

	return serviceModel, nil
}

func responseError(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusBadRequest)
}

func responseJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	bites, _ := json.Marshal(data)
	w.Write(bites)
}
