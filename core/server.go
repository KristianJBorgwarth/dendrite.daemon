// Package server contains core functionality for dendrite server
package server

import (
	"bufio"
	"io"
	"log/slog"
	"os"
)

type Server struct {
	in  *bufio.Reader
	out io.Writer
	log *slog.Logger
}

func NewServer() *Server {
	return &Server{
		in:  bufio.NewReader(os.Stdin),
		out: os.Stdout,
		log: slog.Default(),
	}
}
