package configuration

import (
	"encoding/json"
	"io/ioutil"
)

/*
	Service configuration struct
*/
type Configuration struct {
	// Marathon configuration
	Marathon Marathon

	// Services mapping configuration
	ServicesMapping ServicesMapping

	// HAProxy output configuration
	HAProxy HAProxy
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
	e := json.Unmarshal(content, &config)
	return e
}
