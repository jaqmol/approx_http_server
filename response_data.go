package main

// ResponseData ...
type ResponseData struct {
	Status  int                 `json:"status"`
	Header  map[string][]string `json:"header"`
	BodyB64 string              `json:"bodyB64,omitempty"`
}
