package config

import (
	"encoding/json"
	"fmt"
	"github.com/google/logger"
	"io/ioutil"
	"os"
)

const logPath = "results.log"

type ConfigType struct {
	Secret       string `json:"secret"`
	Verbose      bool   `json:"verbose"`
	Debug        bool   `json:"debug"`
	OutputFolder string `json:"outputFolder"`
}

func Config() ConfigType {
	var config ConfigType
	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	json.Unmarshal(file, &config)
	return config
}

func (config *ConfigType) InitLogging() *logger.Logger {
	lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		logger.Fatalf("Failed to open log file: %v", err)
	}
	// defer lf.Close()

	return logger.Init("LoggerExample", config.Verbose, true, lf)
}
