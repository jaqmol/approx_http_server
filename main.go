package main

import (
	"github.com/jaqmol/approx/axenvs"
	"github.com/jaqmol/approx/axmsg"
)

func main() {
	envs := axenvs.NewEnvs("approx_http_server", []string{"ENDPOINT", "PORT"}, nil)
	errMsg := axmsg.Errors{Source: "approx_http_server"}

	if len(envs.Outs) != 1 {
		errMsg.LogFatal(nil, "HTTP server expects exactly 1 output, but got %v", len(envs.Outs))
	}
	if len(envs.Ins) != 1 {
		errMsg.LogFatal(nil, "HTTP server expects exactly 1 input, but got %v", len(envs.Ins))
	}

	if len(envs.Required["ENDPOINT"]) == 0 {
		errMsg.LogFatal(nil, "HTTP server expects value for env ENDPOINT")
	}
	if len(envs.Required["PORT"]) == 0 {
		errMsg.LogFatal(nil, "HTTP server expects value for env PORT")
	}

	ahs := NewApproxHTTPServer(envs)
	go func() {
		ahs.InitRequestProxy()
		ahs.StartServer()
	}()
	ahs.InitResponseListener()
}
