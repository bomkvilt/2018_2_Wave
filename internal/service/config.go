package service

import (
	"fmt"
	"io/ioutil"
	"path"
)

//~~~~~~~~~~~~~~~~~~~~~~| FServerConfig

// FServiceConfig - server configuration
// easyjson:json
type FServiceConfig struct {
	Root    string            `json:"root"`    //!< config root
	Log     string            `json:"log"`     //!< application relative path to a log file
	Configs map[string]string `json:"configs"` //!< map @Root relative paths of named configs | Note: name => config
}

// LoadConfig - load a configuration file from a disk
func LoadConfig(file string) FServiceConfig {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		Panic(fmt.Errorf(`Unable to load configuration file: %s`, file), err)
	}

	config := FServiceConfig{
		Configs: map[string]string{},
	}
	if err := config.UnmarshalJSON(data); err != nil {
		Panic(fmt.Errorf(`Unexpected error during unmarshaling file: %s`, file), err)
	}
	config.Root = path.Join(path.Dir(file), config.Root)
	return config
}
