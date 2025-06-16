package config

import (
	"errors"
	"runtime"
)

var (
	ErrInvalidWorkerCount     = errors.New("worker count must be greater than 0 and not exceed CPU count")
	ErrInvalidChannelSize     = errors.New("channel size must be greater than 0")
	ErrInvalidLinesPerFile    = errors.New("lines per file must be greater than 0")
	ErrInvalidInputFileName   = errors.New("input file name cannot be empty")
	ErrInvalidOutputFileName  = errors.New("output file name cannot be empty")
)

// ValidateConfig validates the application configuration
func (c *AppConfig) ValidateConfig() error {
	if c.NumWorkers <= 0 || c.NumWorkers > runtime.NumCPU() {
		return ErrInvalidWorkerCount
	}
	
	if c.LinesChannelSize <= 0 {
		return ErrInvalidChannelSize
	}
	
	if c.ResultsChannelSize <= 0 {
		return ErrInvalidChannelSize
	}
	
	if c.LinesPerFile <= 0 {
		return ErrInvalidLinesPerFile
	}
	
	if c.InputFileName == "" {
		return ErrInvalidInputFileName
	}
	
	if c.OutputFileName == "" {
		return ErrInvalidOutputFileName
	}
	
	return nil
} 