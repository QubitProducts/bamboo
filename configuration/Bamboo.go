package configuration

type Bamboo struct {
	// Service host
	Host string

	// Routing configuration storage
	Zookeeper Zookeeper
}
