package service

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"go.uber.org/zap"
)

//~~~~~~~~~~~~~~~~~~~~~~| ILogger

// ILogger - microservice logger interface
type ILogger interface {
	Infof(pattern string, args ...interface{})  //!< formatted message
	Errorf(pattern string, args ...interface{}) //!< formatted message
	Info(message string, pairs ...interface{})  //!< message + list of (key, value, key, value, ...)
	Error(message string, pairs ...interface{}) //!< message + list of (key, value, key, value, ...)
}

//~~~~~~~~~~~~~~~~~~~~~~| Logger

type logger struct {
	logger zap.SugaredLogger
}

// NewLogger - create a defaul logger
func NewLogger(file string) ILogger {
	logger := &logger{}
	logger.Init(file)
	return logger
}

func (lg *logger) Init(file string) {
	// try to cretae the @file
	if stat, err := os.Stat(file); err == nil {
		if stat.IsDir() {
			Panic(fmt.Errorf(`Cannot use directory as a log file: %s`, file), err)
		}
	} else {
		root := path.Dir(file)
		os.Mkdir(root, 0777)

		log, err := os.Create(file)
		if err != nil {
			Panic(fmt.Errorf(`Cannot create a log file: "%s"`, file), err)
		}
		log.Chmod(777)
		log.Close()
	}

	// init the logger's entery
	JSON := []byte(fmt.Sprintf(`{
		"level"				: "debug",
		"encoding"			: "json",
		"outputPaths"		: ["stdout", "%s"],
		"errorOutputPaths"	: ["stderr"],
		"encoderConfig"		:
		{
		  "messageKey"		: "message",
		  "levelKey"		: "level",
		  "levelEncoder"	: "lowercase"
		}
	  }`, file))

	config := zap.Config{}

	if err := json.Unmarshal(JSON, &config); err != nil {
		Panic(fmt.Errorf(`Error during a logger configuratin`), err)
	}

	logger, err := config.Build()
	if err != nil {
		Panic(fmt.Errorf(`Error during a logger configuratin`), err)
	}
	lg.logger = *logger.Sugar()

	lg.Info("Logger started")
}

func (lg *logger) Infof(pattern string, args ...interface{}) {
	lg.logger.Infof(pattern, args...)
}

func (lg *logger) Errorf(pattern string, args ...interface{}) {
	lg.logger.Errorf(pattern, args...)
}

func (lg *logger) Info(message string, pairs ...interface{}) {
	lg.logger.Info(message, pairs)
}

func (lg *logger) Error(message string, pairs ...interface{}) {
	lg.logger.Error(message, pairs)
}
