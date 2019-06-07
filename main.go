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

	// idCounter := 0
	// responseWriterForID := make(map[int]http.ResponseWriter)
	// mutex := &sync.Mutex{}

	// go func() {
	// 	http.HandleFunc(conf.Envs["ENDPOINT"], func(w http.ResponseWriter, r *http.Request) {
	// 		mutex.Lock()
	// 		idCounter++
	// 		id := idCounter
	// 		responseWriterForID[id] = w
	// 		mutex.Unlock()

	// 		reqMsg, err := RequestMessageFromHTTPRequest(id, r)
	// 		if err != nil {
	// 			LogErrorMessage(&id, -5005, err)
	// 			return
	// 		}

	// 		reqMsgBytes, err := json.Marshal(reqMsg)
	// 		if err != nil {
	// 			specError := fmt.Errorf("Error marshalling request message: %v", err.Error())
	// 			LogErrorMessage(&id, -5006, specError)
	// 			return
	// 		}

	// 		nwr, err := conf.Outputs[0].Write(reqMsgBytes)
	// 		if err != nil {
	// 			specError := fmt.Errorf("Error writing request message to output: %v", err.Error())
	// 			LogErrorMessage(&id, -5007, specError)
	// 			return
	// 		}
	// 		if nwr < len(reqMsgBytes) {
	// 			specError := fmt.Errorf("Only %v of %v bytes written to output", nwr, len(reqMsgBytes))
	// 			LogErrorMessage(&id, -5008, specError)
	// 			return
	// 		}
	// 	})

	// 	addr := fmt.Sprintf(":%v", conf.Envs["PORT"])
	// 	err := http.ListenAndServe(addr, nil)
	// 	if err != nil {
	// 		specError := fmt.Errorf("Error starting server: %v", err.Error())
	// 		LogFatalErrorMessage(nil, -5009, specError)
	// 	}
	// }()

	// go func() {
	// 	var hardErr error
	// 	for hardErr == nil {
	// 		resMsgBuffer := make([]byte, 0)
	// 		var nrd int
	// 		nrd, hardErr = conf.Inputs[0].Read(resMsgBuffer)
	// 		if hardErr != nil {
	// 			return
	// 		}

	// 		var resMsg ResponseMessage
	// 		err := json.Unmarshal(resMsgBuffer[:nrd], &resMsg)
	// 		if err != nil {
	// 			specError := fmt.Errorf("Error unmarshalling response message: %v", err.Error())
	// 			LogErrorMessage(nil, -5010, specError)
	// 			return
	// 		}

	// 		mutex.Lock()
	// 		w := responseWriterForID[resMsg.ID]
	// 		delete(responseWriterForID, resMsg.ID)
	// 		mutex.Unlock()

	// 		for key, values := range resMsg.Result.Header {
	// 			for _, value := range values {
	// 				w.Header().Add(key, value)
	// 			}
	// 		}

	// 		w.WriteHeader(resMsg.Result.Status)

	// 		body, err := base64.StdEncoding.DecodeString(resMsg.Result.BodyB64)
	// 		if err != nil {
	// 			specError := fmt.Errorf("Error decoding base64 body: %v", err.Error())
	// 			LogErrorMessage(nil, -5011, specError)
	// 			return
	// 		}

	// 		nwr, err := w.Write(body)
	// 		if err != nil {
	// 			specError := fmt.Errorf("Error writing response body: %v", err.Error())
	// 			LogErrorMessage(nil, -5012, specError)
	// 			return
	// 		}
	// 		if nwr < len(body) {
	// 			specError := fmt.Errorf("Only %v of %v bytes written to response", nwr, len(body))
	// 			LogErrorMessage(nil, -5013, specError)
	// 			return
	// 		}
	// 	}

	// 	if hardErr == io.EOF {
	// 		specError := fmt.Errorf("Unexpected EOL listening for response input")
	// 		LogFatalErrorMessage(nil, -5014, specError)
	// 	} else {
	// 		specError := fmt.Errorf("Unexpected error listening for response input: %v", hardErr.Error())
	// 		LogFatalErrorMessage(nil, -5015, specError)
	// 	}
	// }()
}
