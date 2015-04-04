package configuration

type Bamboo struct {
	// Service host
	Endpoint string
	
	// Service socket binding
	Bind	 string

	// Routing configuration storage
	Zookeeper Zookeeper
}
