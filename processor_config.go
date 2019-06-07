package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

// ProcessorConfig ...
type ProcessorConfig struct {
	Envs    map[string]string
	Inputs  []io.Reader
	Outputs []io.Writer
}

// NewProcessorConfig ...
func NewProcessorConfig(requiredEnvs []string) *ProcessorConfig {
	return &ProcessorConfig{
		Envs: readAllEnvs(requiredEnvs),
		Inputs: openInputs(
			readAllPrefixedEnvs(
				"IN_",
				parseIntEnv("IN_COUNT", 0),
				"stdin",
			),
		),
		Outputs: openOutputs(
			readAllPrefixedEnvs(
				"OUT_",
				parseIntEnv("OUT_COUNT", 0),
				"stdout",
			),
		),
	}
}

func readAllEnvs(names []string) map[string]string {
	values := make(map[string]string)
	for _, name := range names {
		values[name] = readEnv(name)
	}
	return values
}

func readEnv(name string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		LogFatalErrorMessage(nil, -1001, fmt.Errorf("Required env %v not found", name))
	}
	return value
}

func openInputs(inValues []string) []io.Reader {
	inputs := make([]io.Reader, 0)
	for _, name := range inValues {
		if name == "stdin" {
			inputs = append(inputs, os.Stdin)
		} else {
			f, err := os.OpenFile(name, os.O_RDONLY, 0600)
			if err != nil {
				LogFatalErrorMessage(nil, -1002, fmt.Errorf("Error opening named pipe %v for reading: %v", name, err.Error()))
			}
			inputs = append(inputs, f)
		}
	}
	return inputs
}

func openOutputs(outValues []string) []io.Writer {
	outputs := make([]io.Writer, 0)
	for _, name := range outValues {
		if name == "stdout" {
			outputs = append(outputs, os.Stdout)
		} else {
			f, err := os.OpenFile(name, os.O_RDWR, 0600)
			if err != nil {
				LogFatalErrorMessage(nil, -1003, fmt.Errorf("Error opening named pipe %v for writing: %v", name, err.Error()))
			}
			outputs = append(outputs, f)
		}
	}
	return outputs
}

func readAllPrefixedEnvs(prefix string, count int, fallback string) []string {
	acc := make([]string, 0)
	if count == 0 {
		acc = append(acc, fallback)
	} else {
		for i := 0; i < count; i++ {
			name := readIndexedEnv(prefix, i)
			acc = append(acc, name)
		}
	}
	return acc
}

func parseIntEnv(name string, fallback int) int {
	valueStr, ok := os.LookupEnv(name)
	if !ok {
		return fallback
	}
	value64, err := strconv.ParseInt(valueStr, 10, 32)
	if err != nil {
		LogFatalErrorMessage(nil, -1004, fmt.Errorf("Error parsing int env %v: %v", name, err.Error()))
	}
	return int(value64)
}

func readIndexedEnv(prefix string, index int) string {
	name := fmt.Sprintf("%v%v", prefix, index)
	value, ok := os.LookupEnv(name)
	if !ok {
		log.Fatalf("Required env %v not found\n", name)
		LogFatalErrorMessage(nil, -1005, fmt.Errorf("Required env %v not found", name))
	}
	return value
}
