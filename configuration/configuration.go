package configuration

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strconv"
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
		return err
	}
	return json.Unmarshal(content, &config)
}

func FromFlags() (Configuration, error) {
	config := &Configuration{}
	fromEnv(config)
	setValueFromFlag(&config.Marathon.Endpoint, marathonEndpoint)
	setValueFromFlag(&config.Bamboo.Endpoint, bambooEndpoint)
	setValueFromFlag(&config.Bamboo.Zookeeper.Host, bambooZkHost)
	setValueFromFlag(&config.Bamboo.Zookeeper.Path, bambooZkPath)
	setInt64ValueFromFlag(&config.Bamboo.Zookeeper.ReportingDelay, bambooZkReportingDelay)
	setValueFromFlag(&config.HAProxy.TemplatePath, haproxyTemplatePath)
	setValueFromFlag(&config.HAProxy.OutputPath, haproxyOutputPath)
	setValueFromFlag(&config.HAProxy.ReloadCommand, haproxyReloadCommand)
	setBoolValueFromFlag(&config.StatsD.Enabled, statsdEnabled)
	setValueFromFlag(&config.StatsD.Host, statsdHost)
	setValueFromFlag(&config.StatsD.Prefix, statsdPrefix)
	if config.Bamboo.Endpoint == "" {
		return *config, errors.New("The bamboo endpoint cannot be EMPTY!")
	}
	jsonBytes, err := json.Marshal(*config)
	if err != nil {
		return *config, err
	}
	var buffer bytes.Buffer
	err = json.Indent(&buffer, jsonBytes, "", "    ")
	log.Printf("The configuration from flags (& env):\n%s\n\n", buffer.String())
	return *config, err
}

func FromFile(filePath string) (Configuration, error) {
	conf := &Configuration{}
	err := conf.FromFile(filePath)
	fromEnv(conf)
	return *conf, err
}

func fromEnv(conf *Configuration) {
	setValueFromEnv(&conf.Marathon.Endpoint, "MARATHON_ENDPOINT")
	setValueFromEnv(&conf.Bamboo.Endpoint, "BAMBOO_ENDPOINT")
	setValueFromEnv(&conf.Bamboo.Zookeeper.Host, "BAMBOO_ZK_HOST")
	setValueFromEnv(&conf.Bamboo.Zookeeper.Path, "BAMBOO_ZK_PATH")
	setValueFromEnv(&conf.HAProxy.TemplatePath, "HAPROXY_TEMPLATE_PATH")
	setValueFromEnv(&conf.HAProxy.OutputPath, "HAPROXY_OUTPUT_PATH")
	setValueFromEnv(&conf.HAProxy.ReloadCommand, "HAPROXY_RELOAD_CMD")
	setBoolValueFromEnv(&conf.StatsD.Enabled, "STATSD_ENABLED")
	setValueFromEnv(&conf.StatsD.Host, "STATSD_HOST")
	setValueFromEnv(&conf.StatsD.Prefix, "STATSD_PREFIX")
}

func setValueFromEnv(field *string, envVar string) {
	env := os.Getenv(envVar)
	if len(env) > 0 {
		log.Printf("Using environment override %s=%s", envVar, env)
		*field = env
	}
}

func setBoolValueFromEnv(field *bool, envVar string) {
	env := os.Getenv(envVar)
	if len(env) > 0 {
		log.Printf("Using environment override %s=%t", envVar, env)
		x, err := strconv.ParseBool(env)
		if err != nil {
			log.Printf("Error converting boolean value: %s\n", err)
		}
		*field = x
	} else {
		log.Printf("Environment variable not set: %s", envVar)
	}
}
