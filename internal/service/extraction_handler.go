package service

import (
	"assignment/pkg/logger"
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
)

type Record struct {
	Spins      int    `json:"spins"`
	ServerTime string `json:"server_time"`
}

var (
	successfulLines int
	failedLines     int
)

type ExtractionManager struct {
	inputFileName  string
	outputFileName string
	numWorkers     int
	linesPerFile   int           // Max Number of lines per output file
	linesChannel   chan string   // buffered channel for lines
	resultsChannel chan []string // Buffered channel for results
}

func NewExtractionManager(inputFileName, outputFileName string, numWorkers, linesPerFile, linesChannelSize, resultsChannelSize int) *ExtractionManager {
	if inputFileName == "" || outputFileName == "" {
		log.Fatalf("Input or output file name cannot be empty")
	}
	if numWorkers <= 0 || linesPerFile <= 0 || linesChannelSize <= 0 || resultsChannelSize <= 0 {
		log.Fatalf("Configuration values must be greater than zero")
	}

	return &ExtractionManager{
		inputFileName:  inputFileName,
		outputFileName: outputFileName,
		numWorkers:     numWorkers,
		linesPerFile:   linesPerFile,
		linesChannel:   make(chan string, linesChannelSize),
		resultsChannel: make(chan []string, resultsChannelSize),
	}
}

// Extract reads the input file, processes it with multiple workers, and writes the results to output files.
func (p *ExtractionManager) Extract() {
	// Open input file
	inputFile, err := os.Open(p.inputFileName)
	if err != nil {
		logger.Fatal("Error opening input file", logrus.Fields{"error": err})
	}
	defer inputFile.Close()

	p.TriggerWorkers(p.linesChannel, p.resultsChannel)
	p.readInputFile(inputFile, p.linesChannel)
	p.writeResults(p.resultsChannel)

	logger.Info("Processing completed", logrus.Fields{
		"inputFile":       p.inputFileName,
		"outputFile":      p.outputFileName,
		"successfulLines": successfulLines,
		"failedLines":     failedLines,
	})
}

// Read input file line by line and send to workers
func (p *ExtractionManager) readInputFile(inputFile *os.File, lines chan<- string) {
	scanner := bufio.NewScanner(inputFile)
	go func() {
		for scanner.Scan() {
			lines <- scanner.Text()
		}
		close(lines)
	}()
}

// TriggerWorkers manages the worker goroutines,
// ensuring they are started and that the results channel is closed when all workers are done.
func (p *ExtractionManager) TriggerWorkers(lines chan string, results chan []string) {
	var wg sync.WaitGroup
	// Start worker goroutines
	for i := 0; i < p.numWorkers; i++ {
		wg.Add(1)
		go worker(lines, results, &wg)
	}

	// Start a goroutine to close the results channel after workers are done
	go func() {
		wg.Wait()
		close(results)
	}()
}

// responsible to process lines and send extracted data to the results channel
func worker(lines chan string, results chan []string, wg *sync.WaitGroup) {
	defer wg.Done()
	for line := range lines {
		var record Record
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			failedLines++
			logger.Warning("Malformed JSON skipped", logrus.Fields{
				"error": err,
			})
			continue
		}
		successfulLines++
		results <- []string{strconv.Itoa(record.Spins), record.ServerTime}
	}
}

// writeResults listen to result channel and writes the processed results to CSV files.
func (p *ExtractionManager) writeResults(results chan []string) {
	fileIndex := 0
	currentLineCount := 0
	var outputFile *os.File
	var writer *csv.Writer

	for result := range results {
		// Initialize the writer if it is nil
		if writer == nil {
			var err error
			outputFileName := fmt.Sprintf("output-%d.csv", fileIndex)
			outputFile, err = os.Create(outputFileName)
			if err != nil {
				logger.Fatal("Error creating output file", logrus.Fields{"error": err})
			}
			writer = csv.NewWriter(outputFile)
			fileIndex++
		}

		// Write the result to the current file
		if err := writer.Write(result); err != nil {
			logger.Error("Error writing to output file", logrus.Fields{"error": err})
			return
		}
		currentLineCount++

		// Open a new file if needed
		if currentLineCount == p.linesPerFile {
			writer.Flush()
			outputFile.Close()
			writer = nil // Reset writer to trigger reinitialization
			currentLineCount = 0
		}
	}

	// Flush and close the last file
	if writer != nil {
		writer.Flush()
		outputFile.Close()
	}
}

// input {"apple": 1, "banana":2}
func weightedRandomChoice(input map[string]int) string {
	totalWeight := 0
	for _, weight := range input {
		totalWeight += weight
	}

	randomValue := rand.Intn(totalWeight)
	for item, weight := range input {
		if randomValue < weight {
			return item
		}
		randomValue -= weight
	}
	return ""
}

// weightedRandomChoice more efficient version
func weightedRandomChoiceEfficient(input map[string]int) string {
	totalWeight := 0
	for _, weight := range input {
		totalWeight += weight
	}

	randomValue := rand.Intn(totalWeight)
	for item, weight := range input {
		if randomValue < weight {
			return item
		}
		randomValue -= weight
	}
	return ""
}
