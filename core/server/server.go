// Package server contains the server running the JSON-RPC 2.0 Protocol.
package server

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/handlers"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/rpc"
)

type Server struct {
	Db *sql.DB
	handlers map[string]handlers.Handler
}

func NewServer() *Server {
	return &Server{
		handlers: make(map[string]handlers.Handler),
	}
}

func (s *Server) InitDatabase(db *sql.DB) {
	s.Db = db
}

func (s *Server) Register(method string, handler handlers.Handler) {
	s.handlers[method] = handler
}

func (s *Server) Run(r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Bytes()
		var req rpc.Request
		if err := json.Unmarshal(line, &req); err != nil {
			s.write(w, rpc.Response{
				Jsonrpc: "2.0",
				Error:   &rpc.Error{Code: -32700, Message: "parse error"},
			})
			continue
		}

		s.handle(w, req)

	}

	return scanner.Err()
}

func (s *Server) handle(w io.Writer, req rpc.Request) {
	ctx := context.Background()
	handler, ok := s.handlers[req.Method]

	if !ok {
		s.write(w, rpc.Response{
			Jsonrpc: "2.0",
			ID:      req.ID,
			Error:   &rpc.Error{Code: -32601, Message: "method not found"},
		})
	return
}

	result := handler.Handle(ctx, req.Params)

	if req.ID == nil {
		return
	}

	s.write(w, *result)
}

func (s *Server) write(w io.Writer, resp rpc.Response) {
	data, _ := json.Marshal(resp)
	fmt.Fprintln(w, string(data)) 
}
