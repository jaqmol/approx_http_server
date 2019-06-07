package main

import (
	"encoding/json"
	"log"
	"os"
)

// LogFatalErrorMessage ...
func LogFatalErrorMessage(id *int, code int, err error) {
	LogErrorMessage(id, code, err)
	os.Exit(1)
}

// LogErrorMessage ...
func LogErrorMessage(id *int, code int, err error) {
	LogErrorMessageWithData(id, code, err, nil)
}

// LogErrorMessageWithData ...
func LogErrorMessageWithData(id *int, code int, err error, data interface{}) {
	msg := NewErrorMessage(id, code, err.Error(), data)
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Fatalf("Error marshalling error message: %v\n", err.Error())
	}
	msgBytes = append(msgBytes, []byte("\n")...)
	os.Stderr.Write(msgBytes)
}

// NewErrorMessage ...
func NewErrorMessage(id *int, code int, message string, data interface{}) *ErrorMessage {
	return &ErrorMessage{
		JSONRPC: "2.0",
		ID:      id,
		Error: Error{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
}

// ErrorMessage ...
type ErrorMessage struct {
	JSONRPC string `json:"jsonrpc"`
	ID      *int   `json:"id"`
	Error   Error  `json:"error"`
}

// Error ...
type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
