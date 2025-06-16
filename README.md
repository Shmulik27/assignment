# Extraction Manager

## Overview
The `Extraction Manager` is a Go-based application designed to process JSON input files, extract specific fields, and write the results to multiple CSV output files. It ensures that each output file contains a specified maximum number of lines.

## Features
- Reads JSON input files line by line.
- Extracts specific fields from JSON records.
- Writes extracted data to CSV files.
- Automatically creates new output files when the line limit is reached.

## Prerequisites
- Go 1.22 or later installed on your system.
- A valid JSON input file with the required structure.

## Installation

```bash
  go mod tidy
```

## Run application:
```bash
  go run ./cmd/data_extraction
```
The output CSV files will be generated in the project directory.

## Configuration
- `linesPerFile`: Specifies the maximum number of lines per output file. This can be adjusted in the code.
- Input and output file names can be customized in the `writeResults` function.

## Build the application:
```bash
  go build -o bin/assignment ./cmd/data_extraction
```   

## Testing
Run the tests to ensure the application works as expected:
```bash
  go test ./...
```