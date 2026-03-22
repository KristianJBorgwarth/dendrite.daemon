// Package server provides the implementation of the server component of the application.
package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/KristianJBorgwarth/dendrite.daemon/config"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/handlers"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/rpc"
)

type Server struct {
	addr   string
	config *config.Config
}

func New(addr string) *Server {
	return &Server{
		addr: addr,
	}
}

func (s *Server) Start() error {
	slog.Info("server started")
	mux := http.NewServeMux()

	mux.HandleFunc("/rpc", s.RPCHandler)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	return http.ListenAndServe(s.addr, mux)
}

func (s *Server) RPCHandler(w http.ResponseWriter, r *http.Request) {
	var req rpc.Request

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	resp := rpc.Response{
		Jsonrpc: "2.0",
		ID:      req.ID,
	}

	switch req.Method {
	case "initialize":
		cfg, err := handlers.Initialize(req.Params)
		if err != nil {
			resp.Error = map[string]any{
				"code":    -1,
				"message": err.Error(),
			}
		} else {
			s.config = cfg
			resp.Result = map[string]any{"status": "ok"}
		}
	default:
		slog.Warn("Unknown method", "method", req.Method)
		resp.Error = map[string]any{
			"code":    -32601,
			"message": "method not found",
		}
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to write response", "error", err)
	}
}
