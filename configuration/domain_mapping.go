package configuration

/*
	Services mapping information storage configuration
	Only support zookeeper at the moment
*/
type DomainMapping struct {
	Zookeeper Zookeeper
}
