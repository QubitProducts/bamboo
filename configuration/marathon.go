package configuration

import (
	"github.com/QubitProducts/bamboo/Godeps/_workspace/src/github.com/samuel/go-zookeeper/zk"
	"strings"
	"time"
)

/*
	Mesos Marathon configuration
*/
type Marathon struct {
	// comma separated marathon http endpoints including port number
	Endpoint       string
	UseZookeeper   bool
	Zookeeper      Zookeeper
	User           string
	Password       string
	UseEventStream bool
}

func (m Marathon) Endpoints() []string {
	if m.UseZookeeper {
		endpoints, err := zkEndpoints(m.Zookeeper)
		if err != nil {
			return nil
		}
		return endpoints
	}
	return strings.Split(m.Endpoint, ",")
}

func zkEndpoints(zkConf Zookeeper) ([]string, error) {
	// Only tested with marathon 0.11.1, assumes http:// for marathon
	const scheme = "http://"
	const leaderNode = "leader"

	conn, _, err := zk.Connect(zkConf.ConnectionString(), time.Second*10)

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	var leaderPath = zkConf.Path + "/" + leaderNode

	keys, _, err := conn.Children(leaderPath)

	if err != nil {
		return nil, err
	}

	endpoints := make([]string, 0, len(keys))

	for _, childPath := range keys {
		data, _, err := conn.Get(leaderPath + "/" + childPath)
		if err != nil {
			return nil, err
		}
		// TODO configurable http://??
		endpoints = append(endpoints, scheme+string(data))
	}

	return endpoints, nil
}
