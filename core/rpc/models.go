// Package rpc provides a simple implementation of JSON-RPC 2.0 protocol for handling RPC requests and notifications.
package rpc

import (
	"encoding/json"
)

type Request struct {
	Jsonrpc string          `json:"jsonrpc"`
	ID      *int             `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

type Response struct {
	Jsonrpc string          `json:"jsonrpc"`
	ID      *int             `json:"id"`
	Result  any `json:"result,omitempty"`
	Error   *Error           `json:"error,omitempty"`
}

type Notification struct {
	Jsonrpc string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
