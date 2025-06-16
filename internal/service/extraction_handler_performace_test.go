package service

import (
	"assignment/pkg/logger"
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

func TestPerformanceParse(t *testing.T) {
	logger.InitLogger() // Initialize the logger

	inputFileName := "performance_test_input.json"
	outputFileName := "output-%d.csv"
	numLines := 1000000
	generateLargeInputFile(inputFileName, numLines)
	defer os.Remove(inputFileName)

	numWorkers := 8
	linesPerFile := 10000
	linesChannelSize := 100
	resultsChannelSize := 100

	parser := NewExtractionManager(inputFileName, outputFileName, numWorkers, linesPerFile, linesChannelSize, resultsChannelSize)

	startTime := time.Now()
	parser.Extract()
	elapsedTime := time.Since(startTime)

	outputFiles := 0
	for i := 0; ; i++ {
		fileName := fmt.Sprintf(outputFileName, i)
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			break
		}
		outputFiles++
		defer os.Remove(fileName)
	}

	log.Printf("Performance Test Results:")
	log.Printf("Input Lines: %d", numLines)
	log.Printf("Output Files: %d", outputFiles)
	log.Printf("Execution Time: %s", elapsedTime)

	if outputFiles == 0 {
		t.Errorf("No output files were created")
	}
	if elapsedTime > 2*time.Minute {
		t.Errorf("Execution time exceeded threshold: %s", elapsedTime)
	}
}
func generateLargeInputFile(fileName string, numLines int) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Failed to create input file: %v", err)
	}
	defer file.Close()

	for i := 0; i < numLines; i++ {
		line := fmt.Sprintf(`{"spins": %d, "server_time": "2025-05-24 00:00:%02d.99999 UTC"}`, i%100, i%60)
		file.WriteString(line + "\n")
	}
}
