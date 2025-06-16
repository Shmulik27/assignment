package config

import (
	"assignment/pkg/logger"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"os"
)

type AppConfig struct {
	InputFileName      string `json:"inputFileName"`
	OutputFileName     string `json:"outputFileName"`
	NumWorkers         int    `json:"numWorkers"`
	LinesPerFile       int    `json:"linesPerFile"`
	LinesChannelSize   int    `json:"linesChannelSize"`
	ResultsChannelSize int    `json:"resultsChannelSize"`
}

func LoadConfig(configFile string) (*AppConfig, error) {
	file, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config AppConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		logger.Error("Failed to decode configuration", logrus.Fields{"error": err})
		return nil, err
	}
	return &config, nil
}
