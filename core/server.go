// Package server contains core functionality for dendrite server
package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/models"
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

// Run starts the server loop, continuously reading and processing incoming messages
func (s *Server) Run() error {
	for {
		length, err := s.readContentLength()
		if err != nil {
			return err
		}

		body := make([]byte, length)
		if _, err := io.ReadFull(s.in, body); err != nil {
			return err
		}

		s.log.Debug("received", "body", string(body))

		// Try request first
		var req types.Request
		if err := json.Unmarshal(body, &req); err == nil && req.Method != "" {
			s.handleRequest(req)
			continue
		}

		// fallback: notification
		var notif types.Notification
		if err := json.Unmarshal(body, &notif); err == nil {
			s.handleNotification(notif)
			continue
		}
	}
}

func (s *Server) readContentLength() (int, error) {
	for {
		line, err := s.in.ReadString('\n')
		if err != nil {
			return 0, err
		}

		line = strings.TrimSpace(line)
		if line == "" {
			break
		}

		if after, ok :=strings.CutPrefix(line, "Content-Length:"); ok  {
			val := strings.TrimSpace(after)
			return strconv.Atoi(val)
		}
	}
	return 0, fmt.Errorf("missing Content Length")
}

func (s *Server) handleRequest(req types.Request) {
	s.log.Info("handling request", "method", req.Method, "id", req.ID)
}

func (s *Server) handleNotification(notif types.Notification) {
	s.log.Info("handling notification", "method", notif.Method)
}	

