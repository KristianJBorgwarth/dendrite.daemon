package rpc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
)

type Handler interface {
	Handle(json.RawMessage) (any, *Error)
}

type Server struct {
	handlers map[string]Handler
}

func NewServer() *Server {
	return &Server{
		handlers: make(map[string]Handler),
	}
}

func (s *Server) Register(method string, handler Handler) {
	s.handlers[method] = handler
}

func (s *Server) Serve(r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Bytes()
		var req Request
		if err := json.Unmarshal(line, &req); err != nil {
			s.write(w, Response{
				Jsonrpc: "2.0",
				Error:   &Error{Code: -32700, Message: "parse error"},
			})
			continue
		}

		s.handle(w, req)

	}

	return scanner.Err()
}

func (s *Server) handle(w io.Writer, req Request) {
	handler, ok := s.handlers[req.Method]

	if !ok {
		s.respond(w, req.ID, nil, &Error{
			Code:    -32601,
			Message: "method not found",
		})
		return
	}

	result, err := handler.Handle(req.Params)

	if req.ID == nil {
		return
	}

	s.respond(w, req.ID, result, err)
}

func (s *Server) respond(w io.Writer, id *int, result any, err *Error) {
	resp := Response{
		Jsonrpc: "2.0",
		ID:      id,
		Result:  result,
		Error:   err,
	}

	s.write(w, resp)
}

func (s *Server) write(w io.Writer, resp Response) {
	data, _ := json.Marshal(resp)
	fmt.Fprintln(w, string(data)) 
}
