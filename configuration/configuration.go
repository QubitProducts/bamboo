package configuration

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

/*
	Service configuration struct
*/
type Configuration struct {
	// Marathon configuration
	Marathon Marathon

	// Bamboo specific configuration
	Bamboo Bamboo

	// HAProxy output configuration
	HAProxy HAProxy

	// StatsD configuration
	StatsD StatsD
}

/*
	Returns Configuration struct from a given file path

	Parameters:
		filePath: full file path to the JSON configuration
*/
func (config *Configuration) FromFile(filePath string) error {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return json.Unmarshal(content, &config)
}

func FromFile(filePath string) (Configuration, error) {
	conf := &Configuration{}
	err := conf.FromFile(filePath)
	setValueFromEnv(&conf.Marathon.Endpoint, "MARATHON_ENDPOINT")

	setValueFromEnv(&conf.Bamboo.Endpoint, "BAMBOO_ENDPOINT")
	setValueFromEnv(&conf.Bamboo.Zookeeper.Host, "BAMBOO_ZK_HOST")
	setValueFromEnv(&conf.Bamboo.Zookeeper.Path, "BAMBOO_ZK_PATH")

	setValueFromEnv(&conf.HAProxy.TemplatePath, "HAPROXY_TEMPLATE_PATH")
	setValueFromEnv(&conf.HAProxy.OutputPath, "HAPROXY_OUTPUT_PATH")
	setValueFromEnv(&conf.HAProxy.ReloadCommand, "HAPROXY_RELOAD_CMD")
	return *conf, err
}

func setValueFromEnv(field *string, envVar string) {
	env := os.Getenv(envVar)
	if len(env) > 0 {
		log.Printf("Using environment override %s=%s", envVar, env)
		*field = env
	}
}
