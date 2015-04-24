package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/QubitProducts/bamboo/Godeps/_workspace/src/github.com/go-martini/martini"
	zk "github.com/QubitProducts/bamboo/Godeps/_workspace/src/github.com/samuel/go-zookeeper/zk"
	conf "github.com/QubitProducts/bamboo/configuration"
	"github.com/QubitProducts/bamboo/services/service"
)

type ServiceAPI struct {
	Config    *conf.Configuration
	Zookeeper *zk.Conn
}

func (d *ServiceAPI) All(w http.ResponseWriter, r *http.Request) {
	services, err := service.All(d.Zookeeper, d.Config.Bamboo.Zookeeper)

	if err != nil {
		responseError(w, err.Error())
		return
	}

	responseJSON(w, services)
}

func (d *ServiceAPI) Create(w http.ResponseWriter, r *http.Request) {
	serviceModel, err := extractServiceModel(r)

	if err != nil {
		responseError(w, err.Error())
		return
	}

	_, err2 := service.Create(d.Zookeeper, d.Config.Bamboo.Zookeeper, serviceModel.Id, serviceModel.Acl)
	if err2 != nil {
		responseError(w, "Marathon ID might already exist")
		return
	}

	responseJSON(w, serviceModel)
}

func (d *ServiceAPI) Put(params martini.Params, w http.ResponseWriter, r *http.Request) {

	identity := params["_1"]
	println(identity)

	serviceModel, err := extractServiceModel(r)
	if err != nil {
		responseError(w, err.Error())
		return
	}

	_, err1 := service.Put(d.Zookeeper, d.Config.Bamboo.Zookeeper, identity, serviceModel.Acl)
	if err1 != nil {
		responseError(w, err1.Error())
		return
	}

	responseJSON(w, serviceModel)
}

func (d *ServiceAPI) Delete(params martini.Params, w http.ResponseWriter, r *http.Request) {
	err := service.Delete(d.Zookeeper, d.Config.Bamboo.Zookeeper, params["_1"])
	if err != nil {
		responseError(w, err.Error())
		return
	}

	responseJSON(w, new(map[string]string))
}

func extractServiceModel(r *http.Request) (service.Service, error) {
	var serviceModel service.Service
	payload, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(payload, &serviceModel)
	if err != nil {
		return serviceModel, errors.New("Unable to decode JSON request")
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
