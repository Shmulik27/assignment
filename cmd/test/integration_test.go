package test

import (
	"assignment/config"
	"assignment/internal/service"
	"assignment/pkg/logger"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestIntegration(t *testing.T) {

	logger.InitLogger()
	configFile := "test_config.json"
	configContent := `{
  "inputFileName": "test_input.json",
  "outputFileName": "output-%d.csv",
  "numWorkers": 2,
  "linesPerFile": 100,
  "linesChannelSize": 50,
  "resultsChannelSize": 50
 }`
	err := ioutil.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}
	defer os.Remove(configFile)

	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	generator := NewInputFileGenerator(cfg.InputFileName)
	generator.Generate()
	defer os.Remove(cfg.InputFileName)

	parser := service.NewExtractionManager(
		cfg.InputFileName,
		cfg.OutputFileName,
		cfg.NumWorkers,
		cfg.LinesPerFile,
		cfg.LinesChannelSize,
		cfg.ResultsChannelSize,
	)
	go parser.Extract() // Ensure the method is invoked

	// Remove all output files after the test
	for i := 0; i < 100; i++ { // Assuming a maximum of 100 output files
		outputFileName := fmt.Sprintf("output-%d.csv", i)
		defer os.Remove(outputFileName)
	}

	// Wait for the files to be created
	time.Sleep(100 * time.Millisecond)

	outputFileName := "output-0.csv"
	file, err := os.Open(outputFileName)
	if err != nil {
		t.Fatalf("Failed to open output file: %v", err)
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if len(content) == 0 {
		t.Errorf("Output file is empty")
	}
}

type InputRecord struct {
	Spins         int    `json:"spins"`
	Time          string `json:"time"`
	ServerTime    string `json:"server_time"`
	InsertionDate string `json:"insertion_date"`
}

type InputFileGenerator struct {
	inputFileName string
}

func NewInputFileGenerator(inputFileName string) *InputFileGenerator {
	return &InputFileGenerator{
		inputFileName: inputFileName,
	}
}

func (i *InputFileGenerator) Generate() {
	file, err := os.Create(i.inputFileName)
	if err != nil {
		fmt.Printf("Error creating input file: %v\n", err)
		return
	}
	defer file.Close()

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10000; i++ {
		record := InputRecord{
			Spins:         rand.Intn(100), // Random spins between 0 and 99
			Time:          randomTimestamp(),
			ServerTime:    randomTimestamp(),
			InsertionDate: randomTimestamp(),
		}
		data, err := json.Marshal(record)
		if err != nil {
			fmt.Printf("Error marshaling JSON: %v\n", err)
			return
		}
		file.Write(data)
		file.WriteString("\n")
	}
}

func randomTimestamp() string {
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	randomTime := start.Add(time.Duration(rand.Int63n(int64(end.Sub(start)))))
	return randomTime.Format("2006-01-02 15:04:05.99999 UTC")
}
