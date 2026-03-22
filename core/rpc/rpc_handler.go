// Package rpc provides a simple implementation of JSON-RPC 2.0 protocol for handling RPC requests and notifications.
package rpc

import (
	"encoding/json"
	"net/http"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/handlers"
)

type Request struct {
	Jsonrpc string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

type Response struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  any    `json:"result,omitempty"`
	Error   any    `json:"error,omitempty"`
}

type Notification struct {
	Jsonrpc string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

func RPCHandler(w http.ResponseWriter, r *http.Request) {
	var req Request

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	resp := Response{
		Jsonrpc: "2.0",
		ID:      req.ID,
	}

	switch req.Method {
	case "initialize":
		result, err := handlers.Initialize(req.Params)
		if err != nil {
			resp.Error = err.Error()
		} else {
			resp.Result = result
		}
	default:
		resp.Error = "Method not found"
	}
}
