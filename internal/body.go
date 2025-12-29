package internal

import (
	"fmt"
)

// Result represents a successful response in PlayFab API.
// Make sure to specify the T generic type to whatever you want in the Data.
type Result[T any] struct {
	// StatusCode is the HTTP status code of the response, e.g. 200.
	StatusCode int `json:"code,omitempty"`
	// Data is the payload of the Result.
	Data T `json:"data,omitempty"`
	// Status is the HTTP status of the response, e.g. OK.
	Status string `json:"status,omitempty"`
}

// Error represents an error included in the response body.
type Error struct {
	// StatusCode is the HTTP status code of the response, e.g. 404.
	StatusCode int `json:"code,omitempty"`
	// Type indicates the type of the Error.
	Type string `json:"error,omitempty"`
	// Code is the numerical error code for the Error.
	Code int `json:"errorCode,omitempty"`
	// Details encapsulates custom data specific to the Error.
	Details map[string][]string `json:"errorDetails,omitempty"`
	// Message optionally contains a human-readable message describing the Error.
	Message string `json:"errorMessage,omitempty"`
	// Status is the HTTP status of the response, e.g. Not Found.
	Status string `json:"status,omitempty"`
}

// Error returns a string representation of the Error.
func (err Error) Error() string {
	s := fmt.Sprintf("playfab: %d", err.Code)
	if err.Type != "" {
		s += "(" + err.Type + ")"
	}
	if err.Message != "" {
		s += ": " + err.Message
	}
	return s
}
