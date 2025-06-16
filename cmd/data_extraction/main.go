package main

import (
	"assignment/config"
	"assignment/internal/service"
	"assignment/pkg/logger"
	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize logger with default configuration
	logConfig := logger.LogConfig{
		Level:      "info",
		Format:     "text",
		OutputPath: "",
	}
	if err := logger.InitLogger(logConfig); err != nil {
		panic(err)
	}

	// Load configuration
	config, err := config.LoadConfig("config/app_configuration.json")
	if err != nil {
		logger.Fatal("Failed to load configuration", logrus.Fields{"error": err})
	}

	// Extract input file
	extractionManager := service.NewExtractionManager(
		config.InputFileName,
		config.OutputFileName,
		config.NumWorkers,
		config.LinesPerFile,
		config.LinesChannelSize,
		config.ResultsChannelSize,
	)
	extractionManager.Extract()
}
