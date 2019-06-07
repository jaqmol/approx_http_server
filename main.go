package main

import (
	"fmt"
)

func main() {
	conf := NewProcessorConfig([]string{"ENDPOINT", "PORT"})

	if len(conf.Outputs) != 1 {
		LogFatalErrorMessage(nil, -5001, fmt.Errorf("HTTP server expects exactly 1 output, but got %v", len(conf.Outputs)))
	}
	if len(conf.Inputs) != 1 {
		LogFatalErrorMessage(nil, -5002, fmt.Errorf("HTTP server expects exactly 1 input, but got %v", len(conf.Inputs)))
	}

	if len(conf.Envs["ENDPOINT"]) == 0 {
		LogFatalErrorMessage(nil, -5003, fmt.Errorf("HTTP server expects value for env ENDPOINT"))
	}
	if len(conf.Envs["PORT"]) == 0 {
		LogFatalErrorMessage(nil, -5004, fmt.Errorf("HTTP server expects value for env PORT"))
	}

	ahs := NewApproxHTTPServer(conf)
	go func() {
		ahs.InitRequestProxy()
		ahs.StartServer()
	}()
	ahs.InitResponseListener()
}
