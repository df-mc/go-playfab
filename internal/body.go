package internal

import (
	"go/token"
	"strconv"
	"strings"
)

type Result[T any] struct {
	StatusCode int    `json:"code,omitempty"`
	Data       T      `json:"data,omitempty"`
	Status     string `json:"status,omitempty"`
}

type Error struct {
	StatusCode int                 `json:"code,omitempty"`
	Type       string              `json:"error,omitempty"`
	Code       int                 `json:"errorCode,omitempty"`
	Details    map[string][]string `json:"errorDetails,omitempty"`
	Message    string              `json:"errorMessage,omitempty"`
	Status     string              `json:"status,omitempty"`
}

func (err Error) Error() string {
	b := &strings.Builder{}
	b.WriteString("playfab:")

	b.WriteByte(' ')
	b.WriteString(strconv.Itoa(err.Code))

	if err.Type != "" {
		b.WriteByte(' ')
		b.WriteByte(byte(token.LPAREN))
		b.WriteString(err.Type)
		b.WriteByte(byte(token.RPAREN))
	}
	if err.Message != "" && err.Message != err.Type {
		b.WriteByte(byte(token.COLON))
		b.WriteByte(' ')
		b.WriteString(strconv.Quote(err.Message))
	}
	return b.String()
}
