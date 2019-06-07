package main

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
)

// RequestMessageFromHTTPRequest ...
func RequestMessageFromHTTPRequest(id int, r *http.Request) (*RequestMessage, error) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return &RequestMessage{
		JSONRPC: "2.0",
		ID:      id,
		Method:  r.Method,
		Params: Params{
			URL:     r.URL.String(),
			Header:  r.Header,
			BodyB64: base64.StdEncoding.EncodeToString(bodyBytes),
		},
	}, nil
}

// RequestMessage ...
type RequestMessage struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Method  string `json:"method"`
	Params  Params `json:"params"`
}

// Params ...
type Params struct {
	URL     string              `json:"url"`
	Header  map[string][]string `json:"header"`
	BodyB64 string              `json:"bodyB64,omitempty"`
}
