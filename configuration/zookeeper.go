package configuration

import (
	"strings"
	"time"
)

/*
	Zookeeper configuration set
*/
type Zookeeper struct {
	// comma separated host:port connection strings set
	Host string
	// zookeeper path
	Path string
	// Delay n seconds to report change event
	ReportingDelay int64

	// TODO: authentication parameters for zookeeper
}

func (zk Zookeeper) Delay() time.Duration {
	return time.Duration(zk.ReportingDelay) * time.Second
}

func (zk Zookeeper) ConnectionString() []string {
	return strings.Split(zk.Host, ",")
}
