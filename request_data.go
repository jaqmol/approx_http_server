package main

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"

	"github.com/jaqmol/approx/axmsg"
)

// RequestActionFromHTTPRequest ...
func RequestActionFromHTTPRequest(id int, r *http.Request) (*axmsg.Action, error) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return axmsg.NewAction(
		&id,
		nil,
		"http_request",
		nil,
		RequestData{
			Method:  r.Method,
			URL:     r.URL.String(),
			Header:  r.Header,
			BodyB64: base64.StdEncoding.EncodeToString(bodyBytes),
		},
	), nil
}

// RequestData ...
type RequestData struct {
	Method  string              `json:"method"`
	URL     string              `json:"url"`
	Header  map[string][]string `json:"header"`
	BodyB64 string              `json:"bodyB64,omitempty"`
}
