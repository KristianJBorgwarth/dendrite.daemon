// Package server provides the implementation of the server component of the application.
package server

import (
	"log/slog"
	"net/http"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/rpc"
)

type Server struct {
	addr string
}

func New(log *slog.Logger, addr string) *Server {
	return &Server{
		addr: addr,
	}
}

func (s *Server) Start() error {
	slog.Info("server started")
	mux := http.NewServeMux()

	mux.HandleFunc("/rpc", rpc.RPCHandler)

	return http.ListenAndServe(s.addr, mux)
}
