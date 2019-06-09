package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/jaqmol/approx/errormsg"
	"github.com/jaqmol/approx/processorconf"
)

// NewApproxHTTPServer ...
func NewApproxHTTPServer(conf *processorconf.ProcessorConf) *ApproxHTTPServer {
	return &ApproxHTTPServer{
		errMsg:      &errormsg.ErrorMsg{Processor: "approx_http_server"},
		conf:        conf,
		output:      conf.Outputs[0],
		input:       conf.Inputs[0],
		idCounter:   0,
		respWrForID: make(map[int]http.ResponseWriter),
		mutex:       &sync.Mutex{},
	}
}

// ApproxHTTPServer ...
type ApproxHTTPServer struct {
	errMsg      *errormsg.ErrorMsg
	conf        *processorconf.ProcessorConf
	output      *bufio.Writer
	input       *bufio.Reader
	idCounter   int
	respWrForID map[int]http.ResponseWriter
	mutex       *sync.Mutex
}

// InitRequestProxy ...
func (a *ApproxHTTPServer) InitRequestProxy() {
	http.HandleFunc(a.conf.Envs["ENDPOINT"], func(w http.ResponseWriter, r *http.Request) {
		id := a.registerRespWriterWithID(w)

		reqMsg, err := RequestMessageFromHTTPRequest(id, r)
		if err != nil {
			a.errMsg.Log(&id, err.Error())
			return
		}

		reqMsgBytes, err := json.Marshal(reqMsg)
		if err != nil {
			a.errMsg.Log(&id, "Error marshalling request message: %v", err.Error())
			return
		}

		reqMsgBytes = append(reqMsgBytes, '\n')
		_, err = a.output.Write(reqMsgBytes)
		if err != nil {
			a.errMsg.Log(&id, "Error writing request message to output: %v", err.Error())
			return
		}

		err = a.output.Flush()
		if err != nil {
			a.errMsg.Log(&id, "Error flushing written message to output: %v", err.Error())
			return
		}
	})
}

// StartServer ...
func (a *ApproxHTTPServer) StartServer() {
	addr := fmt.Sprintf(":%v", a.conf.Envs["PORT"])
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		a.errMsg.LogFatal(nil, "Error starting server: %v", err.Error())
	}
}

// InitResponseListener ...
func (a *ApproxHTTPServer) InitResponseListener() {
	var hardErr error
	for hardErr == nil {
		var resMsgBytes []byte
		resMsgBytes, hardErr = a.input.ReadBytes('\n')
		if hardErr != nil {
			break
		}

		var resMsg ResponseMessage
		err := json.Unmarshal(resMsgBytes, &resMsg)
		if err != nil {
			a.errMsg.Log(nil, "Error unmarshalling response message: %v", err.Error())
			continue
		}

		w := a.writeHeadersWithRespWriter(&resMsg)

		body, err := base64.StdEncoding.DecodeString(resMsg.Result.BodyB64)
		if err != nil {
			a.errMsg.Log(nil, "Error decoding base64 body: %v", err.Error())
			continue
		}

		nwr, err := w.Write(body)
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

func (a *ApproxHTTPServer) writeHeadersWithRespWriter(resMsg *ResponseMessage) http.ResponseWriter {
	w := a.unregisterRespWriter(resMsg.ID)
	for key, values := range resMsg.Result.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resMsg.Result.Status)
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
