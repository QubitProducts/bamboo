package cmd

import (
	"flag"
	conf "bamboo/configuration"
)

// Commandline arguments
var configFilePath string
// shared configuration
var config conf.Configuration
var configLoaded bool = false

func init() {
	flag.StringVar(&configFilePath, "config", "config/development.json", "Full path of the configuration JSON file")
}

func parseFlagOnce() {
	if !flag.Parsed() {
		flag.Parse()
	}

	if configLoaded != true {
		config = conf.Configuration{}
		config.FromFile(configFilePath)
		configLoaded = true
	}
}

func GetConfigFilePath() string {
	parseFlagOnce()
	return configFilePath
}

func GetConfiguration() conf.Configuration {
	parseFlagOnce()
	return config
}
