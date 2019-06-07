package main

// ResponseMessage ...
type ResponseMessage struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  Result `json:"result"`
}

// Result ...
type Result struct {
	Status  int                 `json:"status"`
	Header  map[string][]string `json:"header"`
	BodyB64 string              `json:"bodyB64,omitempty"`
}
