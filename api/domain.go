package api

import (
	"errors"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/zenazn/goji/web"
	"github.com/samuel/go-zookeeper/zk"


	conf "bamboo/configuration"
	service "bamboo/services/domain"

)

type Domain struct {
	Config    conf.Configuration
	Zookeeper *zk.Conn
}


type DomainModel struct {
	Id string `param:"id"`
	Value string `param:"value"`
}

func (d Domain) All(w http.ResponseWriter, r *http.Request) {
	domains, err := service.All(d.Zookeeper, d.Config.DomainMapping.Zookeeper)

	if err != nil {
		fmt.Println(err)
		responseError(w, err.Error())
		return
	}

	responseJSON(w, domains)
}

func (d Domain) Create(w http.ResponseWriter, r *http.Request) {
	domainModel, err := extractDomainModel(r)

	if err != nil {
		responseError(w, err.Error())
		return
	}

	_, err2 := service.Create(d.Zookeeper, d.Config.DomainMapping.Zookeeper, domainModel.Id, domainModel.Value)
	if err2 != nil {
		responseError(w, "Service id might already exist")
		return
	}

	responseJSON(w, domainModel)
}

func (d Domain) Put(c web.C, w http.ResponseWriter, r *http.Request) {
	identifier := c.URLParams["id"]
	domainModel, err := extractDomainModel(r)
	if err != nil {
		responseError(w, err.Error())
		return
	}

	_, err1 := service.Put(d.Zookeeper, d.Config.DomainMapping.Zookeeper, identifier, domainModel.Value)
	if err1 != nil {
		responseError(w, err1.Error())
		return
	}

	responseJSON(w, domainModel)
}


func (d Domain) Delete(c web.C, w http.ResponseWriter, r *http.Request) {
	identifier := c.URLParams["id"]
	err := service.Delete(d.Zookeeper, d.Config.DomainMapping.Zookeeper, identifier)
	if err != nil {
		responseError(w, err.Error())
		return
	}

	responseJSON(w, new(map[string]string))
}


func extractDomainModel(r *http.Request) (DomainModel, error) {
	var domainModel DomainModel
	payload := make([]byte, r.ContentLength)
	r.Body.Read(payload)

	err1 := json.Unmarshal(payload, &domainModel)
	if err1 != nil {
		return domainModel, errors.New("Unable to decode JSON request")
	}

	return domainModel, nil
}


func responseError(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusInternalServerError)
}

func responseJSON(w http.ResponseWriter, data interface {}) {
	w.Header().Set("Content-Type", "application/json")
	bites, _ := json.Marshal(data)
	w.Write(bites)
}
