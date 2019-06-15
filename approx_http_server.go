package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/jaqmol/approx/axenvs"
	"github.com/jaqmol/approx/axmsg"
)

// NewApproxHTTPServer ...
func NewApproxHTTPServer(envs *axenvs.Envs) *ApproxHTTPServer {
	ins, outs := envs.InsOuts()
	return &ApproxHTTPServer{
		errMsg:      &axmsg.Errors{Source: "approx_http_server"},
		envs:        envs,
		output:      axmsg.NewWriter(&outs[0]),
		input:       axmsg.NewReader(&ins[0]),
		idCounter:   0,
		respWrForID: make(map[int]http.ResponseWriter),
		mutex:       &sync.Mutex{},
	}
}

// ApproxHTTPServer ...
type ApproxHTTPServer struct {
	errMsg      *axmsg.Errors
	envs        *axenvs.Envs
	output      *axmsg.Writer
	input       *axmsg.Reader
	idCounter   int
	respWrForID map[int]http.ResponseWriter
	mutex       *sync.Mutex
}

// InitRequestProxy ...
func (a *ApproxHTTPServer) InitRequestProxy() {
	http.HandleFunc(
		a.envs.Required["ENDPOINT"],
		func(w http.ResponseWriter, r *http.Request) {
			id := a.registerRespWriterWithID(w)

			reqAction, err := RequestActionFromHTTPRequest(id, r)
			err = a.output.Write(reqAction)
			if err != nil {
				a.errMsg.Log(&id, "Error writing request message to output: %v", err.Error())
				return
			}
		},
	)
}

// StartServer ...
func (a *ApproxHTTPServer) StartServer() {
	addr := fmt.Sprintf(":%v", a.envs.Required["PORT"])
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		a.errMsg.LogFatal(nil, "Error starting server: %v", err.Error())
	}
}

// InitResponseListener ...
func (a *ApproxHTTPServer) InitResponseListener() {
	var hardErr error
	for hardErr == nil {
		action, data, hardErr := a.readResponseActionAndData()
		if hardErr != nil {
			break
		}

		writer := a.writeHeadersWithRespWriter(*action.ID, data)

		body, err := base64.StdEncoding.DecodeString(data.BodyB64)
		if err != nil {
			a.errMsg.Log(nil, "Error decoding base64 body: %v", err.Error())
			continue
		}

		nwr, err := writer.Write(body)
		if err != nil {
			a.errMsg.Log(nil, "Error writing response body: %v", err.Error())
			continue
		}
		if nwr < len(body) {
			a.errMsg.Log(nil, "Only %v of %v bytes written to response", nwr, len(body))
			continue
		}
	}

	if hardErr == io.EOF {
		a.errMsg.LogFatal(nil, "Unexpected EOL listening for response input")
	} else {
		a.errMsg.LogFatal(nil, "Unexpected error listening for response input: %v", hardErr.Error())
	}
}

func (a *ApproxHTTPServer) readResponseActionAndData() (*axmsg.Action, *ResponseData, error) {
	action, rawData, err := a.input.Read()
	if err != nil {
		return nil, nil, err
	}

	var data ResponseData
	err = json.Unmarshal(rawData, &data)
	if err != nil {
		return nil, nil, err
	}
	return action, &data, nil
}

func (a *ApproxHTTPServer) writeHeadersWithRespWriter(id int, data *ResponseData) http.ResponseWriter {
	w := a.unregisterRespWriter(id)
	for key, values := range data.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(data.Status)
	return w
}

func (a *ApproxHTTPServer) registerRespWriterWithID(w http.ResponseWriter) int {
	a.mutex.Lock()
	a.idCounter++
	id := a.idCounter
	a.respWrForID[id] = w
	a.mutex.Unlock()
	return id
}

func (a *ApproxHTTPServer) unregisterRespWriter(id int) http.ResponseWriter {
	a.mutex.Lock()
	w := a.respWrForID[id]
	delete(a.respWrForID, id)
	a.mutex.Unlock()
	return w
}
