// Package server contains the server running the JSON-RPC 2.0 Protocol.
package server

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/handlers"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/rpc"
)

type Server struct {
	handlers map[string]handlers.Handler
}

func NewServer() *Server {
	return &Server{
		handlers: make(map[string]handlers.Handler),
	}
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
		s.respond(w, req.ID, nil, &rpc.Error{
			Code:    -32601,
			Message: "method not found",
		})
		return
	}

	result, err := handler.Handle(ctx, req.Params)
	if err != nil {
		s.respond(w, req.ID, nil, &rpc.Error{
			Code:    -32000,
			Message: err.Error(),
		})
		return
	}

	if req.ID == nil {
		return
	}

	s.respond(w, req.ID, result, nil)
}

func (s *Server) respond(w io.Writer, id *int, result any, err *rpc.Error) {
	resultJSON, _ := json.Marshal(result)
	resp := rpc.Response{
		Jsonrpc: "2.0",
		ID:      id,
		Result:  resultJSON,
		Error:   err,
	}

	s.write(w, resp)
}

func (s *Server) write(w io.Writer, resp rpc.Response) {
	data, _ := json.Marshal(resp)
	fmt.Fprintln(w, string(data))
}
