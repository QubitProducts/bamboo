package api

import(
	"net/http"

	"github.com/samuel/go-zookeeper/zk"
	conf "bamboo/configuration"
)

type Domain struct {
	Config conf.Configuration
	Zookeeper zk.Conn
}

func (d Domain) All(w http.ResponseWriter, r *http.Request) {

}

func (d Domain) Create(w http.ResponseWriter, r *http.Request) {

}

func (d Domain) Delete(w http.ResponseWriter, r *http.Request) {

}

func (d Domain) Put(w http.ResponseWriter, r *http.Request) {

}


