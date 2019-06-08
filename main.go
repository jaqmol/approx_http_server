package main

import (
	"github.com/jaqmol/approx/errormsg"
	"github.com/jaqmol/approx/processorconf"
)

func main() {
	conf := processorconf.NewProcessorConf("approx_http_server", []string{"ENDPOINT", "PORT"})

	if len(conf.Outputs) != 1 {
		errormsg.LogFatal("approx_http_server", nil, -5001, "HTTP server expects exactly 1 output, but got %v", len(conf.Outputs))
	}
	if len(conf.Inputs) != 1 {
		errormsg.LogFatal("approx_http_server", nil, -5002, "HTTP server expects exactly 1 input, but got %v", len(conf.Inputs))
	}

	if len(conf.Envs["ENDPOINT"]) == 0 {
		errormsg.LogFatal("approx_http_server", nil, -5003, "HTTP server expects value for env ENDPOINT")
	}
	if len(conf.Envs["PORT"]) == 0 {
		errormsg.LogFatal("approx_http_server", nil, -5004, "HTTP server expects value for env PORT")
	}

	ahs := NewApproxHTTPServer(conf)
	go func() {
		ahs.InitRequestProxy()
		ahs.StartServer()
	}()
	ahs.InitResponseListener()
}
