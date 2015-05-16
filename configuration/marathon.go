package configuration

import (
	"strings"
)

/*
	Mesos Marathon configuration
*/
type Marathon struct {
	// comma separated marathon http endpoints including port number
	Endpoint string

	// username for HTTP Basic Authentication
	AuthUsername string

	// password for HTTP Basic Authentication
	AuthPassword string
}

func (m Marathon) Endpoints() []string {
	return strings.Split(m.Endpoint, ",")
}
