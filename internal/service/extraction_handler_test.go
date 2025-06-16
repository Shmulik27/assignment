package service

import (
	"io/ioutil"
	"os"
	"sync"
	"testing"
	"time"
)

func TestWorker(t *testing.T) {
	lines := make(chan string, 1)
	results := make(chan []string, 1)

	lines <- `{"spins": 10, "server_time": "2025-05-24 00:00:01.99999 UTC"}`
	close(lines)

	var wg sync.WaitGroup
	wg.Add(1)
	go worker(lines, results, &wg)
	wg.Wait()
	close(results)

	result := <-results
	if result[0] != "10" || result[1] != "2025-05-24 00:00:01.99999 UTC" {
		t.Errorf("Worker failed to parse JSON correctly, got %v", result)
	}
}

func TestWriteResults(t *testing.T) {
	results := make(chan []string, 1)
	results <- []string{"10", "2025-05-24 00:00:01.99999 UTC"}
	close(results)

	parser := NewExtractionManager("test_input.json", "output-%d.csv", 1, 1, 1, 1)
	go parser.writeResults(results) // Ensure the method is invoked
	outputFileName := "output-0.csv"
	defer os.Remove(outputFileName) // Ensure the file is removed after the test

	// Wait for the file to be created
	time.Sleep(100 * time.Millisecond)

	file, err := os.Open(outputFileName)
	if err != nil {
		t.Fatalf("Failed to open output file: %v", err)
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	expected := "10,2025-05-24 00:00:01.99999 UTC\n"
	if string(content) != expected {
		t.Errorf("Output file content mismatch, expected '%s', got '%s'", expected, string(content))
	}
}
