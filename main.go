package main

import (
	"github.com/jaqmol/approx/axmsg"
	"github.com/jaqmol/approx/processorconf"
)

func main() {
	conf := processorconf.NewProcessorConf("approx_http_server", []string{"ENDPOINT", "PORT"})
	errMsg := axmsg.Errors{Source: "approx_http_server"}

	if len(conf.Outputs) != 1 {
		errMsg.LogFatal(nil, "HTTP server expects exactly 1 output, but got %v", len(conf.Outputs))
	}
	if len(conf.Inputs) != 1 {
		errMsg.LogFatal(nil, "HTTP server expects exactly 1 input, but got %v", len(conf.Inputs))
	}

	if len(conf.Envs["ENDPOINT"]) == 0 {
		errMsg.LogFatal(nil, "HTTP server expects value for env ENDPOINT")
	}
	if len(conf.Envs["PORT"]) == 0 {
		errMsg.LogFatal(nil, "HTTP server expects value for env PORT")
	}

	ahs := NewApproxHTTPServer(conf)
	go func() {
		ahs.InitRequestProxy()
		ahs.StartServer()
	}()
	ahs.InitResponseListener()
}
