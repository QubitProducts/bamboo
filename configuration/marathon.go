package configuration

/*
	Mesos Marathon configuration
 */
type Marathon struct {
	// marathon http endpoint including port number
	Endpoint string
	// Zookeeper setting for marathon
	Zookeeper Zookeeper
}
