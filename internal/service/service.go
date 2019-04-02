package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
)

//go:generate easyjson -output_filename jsons.go .

//~~~~~~~~~~~~~~~~~~~~~~| IService

// IService - microservice inviroment
type IService interface {
	Logger() ILogger                       //!< get service logger or nullptr
	Config(json.Unmarshaler, string) error //!< get a subconfig | \see: FServiceConfig.Configs
}

//~~~~~~~~~~~~~~~~~~~~~~| Service

type service struct {
	config   FServiceConfig    //!< service configuration file
	logger   ILogger           //!< service logger
	subconfs map[string][]byte //!< cashed configuration files
}

// NewService - cretae a default service
func NewService(config FServiceConfig) IService {
	service := &service{
		subconfs: map[string][]byte{},
	}
	service.Init(config)
	return service
}

func (sv *service) Init(config FServiceConfig) {
	sv.config = config
	sv.logger = NewLogger(config.Log)
	if err := sv.loadSubconfigs(); err != nil {
		panic(err)
	}
}

func (sv *service) Logger() ILogger {
	return sv.logger
}

func (sv *service) Config(conf json.Unmarshaler, name string) error {
	if data, ok := sv.subconfs[name]; ok {
		return conf.UnmarshalJSON(data)
	}
	return fmt.Errorf(`Unknown config name: %s`, name)
}

//~~~~~~~| internal

func (sv *service) loadSubconfigs() error {
	for name, file := range sv.config.Configs {
		// create a config relative path
		if !path.IsAbs(file) {
			file = path.Join(sv.config.Root, file)
		}

		// read the file
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return Error(fmt.Errorf(`Error during openin subconfig "%s" : (%s)`, name, file), err)
		}
		sv.subconfs[name] = data
	}
	return nil
}
