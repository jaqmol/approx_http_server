package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

// NewApproxHTTPServer ...
func NewApproxHTTPServer(conf *ProcessorConfig) *ApproxHTTPServer {
	return &ApproxHTTPServer{
		conf:        conf,
		idCounter:   0,
		respWrForID: make(map[int]http.ResponseWriter),
		mutex:       &sync.Mutex{},
	}
}

// ApproxHTTPServer ...
type ApproxHTTPServer struct {
	conf        *ProcessorConfig
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
			LogErrorMessage(&id, -5005, err)
			return
		}

		reqMsgBytes, err := json.Marshal(reqMsg)
		if err != nil {
			specError := fmt.Errorf("Error marshalling request message: %v", err.Error())
			LogErrorMessage(&id, -5006, specError)
			return
		}

		reqMsgBytes = append(reqMsgBytes, []byte("\n")...)

		nwr, err := a.conf.Outputs[0].Write(reqMsgBytes)
		if err != nil {
			specError := fmt.Errorf("Error writing request message to output: %v", err.Error())
			LogErrorMessage(&id, -5007, specError)
			return
		}
		if nwr < len(reqMsgBytes) {
			specError := fmt.Errorf("Only %v of %v bytes written to output", nwr, len(reqMsgBytes))
			LogErrorMessage(&id, -5008, specError)
			return
		}
	})
}

// StartServer ...
func (a *ApproxHTTPServer) StartServer() {
	addr := fmt.Sprintf(":%v", a.conf.Envs["PORT"])
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		specError := fmt.Errorf("Error starting server: %v", err.Error())
		LogFatalErrorMessage(nil, -5009, specError)
	}
}

// InitResponseListener ...
func (a *ApproxHTTPServer) InitResponseListener() {
	reader := bufio.NewReader(a.conf.Inputs[0])
	var hardErr error
	for hardErr == nil {
		var resMsgBytes []byte
		resMsgBytes, hardErr = reader.ReadBytes('\n')
		// = a.conf.Inputs[0].Read(resMsgBuffer)
		if hardErr != nil {
			return
		}

		var resMsg ResponseMessage
		err := json.Unmarshal(resMsgBytes, &resMsg)
		if err != nil {
			specError := fmt.Errorf("Error unmarshalling response message: %v", err.Error())
			LogErrorMessage(nil, -5010, specError)
			continue
		}

		w := a.unregisterRespWriter(resMsg.ID)

		for key, values := range resMsg.Result.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		w.WriteHeader(resMsg.Result.Status)

		body, err := base64.StdEncoding.DecodeString(resMsg.Result.BodyB64)
		if err != nil {
			specError := fmt.Errorf("Error decoding base64 body: %v", err.Error())
			LogErrorMessage(nil, -5011, specError)
			continue
		}

		nwr, err := w.Write(body)
		if err != nil {
			specError := fmt.Errorf("Error writing response body: %v", err.Error())
			LogErrorMessage(nil, -5012, specError)
			continue
		}
		if nwr < len(body) {
			specError := fmt.Errorf("Only %v of %v bytes written to response", nwr, len(body))
			LogErrorMessage(nil, -5013, specError)
			continue
		}
	}

	if hardErr == io.EOF {
		specError := fmt.Errorf("Unexpected EOL listening for response input")
		LogFatalErrorMessage(nil, -5014, specError)
	} else {
		specError := fmt.Errorf("Unexpected error listening for response input: %v", hardErr.Error())
		LogFatalErrorMessage(nil, -5015, specError)
	}
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
