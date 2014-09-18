package configuration

type Bamboo struct {
	// Service host
	Endpoint string

	// Routing configuration storage
	Zookeeper Zookeeper
}
