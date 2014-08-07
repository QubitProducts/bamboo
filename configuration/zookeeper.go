package configuration

import "strings"

/*
	Zookeeper configuration set
 */
type Zookeeper struct {
	// comma separated host:port connection strings set
	Host string
	// zookeeper path
	Path string
	// TODO: authentication parameters for zookeeper
}

func (zk Zookeeper) ConnectionString() []string {
	return strings.Split(zk.Host, ",")
}
