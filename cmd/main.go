package main

import (
	"log/slog"
	"os"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/handlers"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/server"
	_ "modernc.org/sqlite"
)

func main() {
	server := server.NewServer()

	server.Register("initialize", handlers.InitializeHandler{})

	if err := server.Run(os.Stdin, os.Stdout); err != nil {
		slog.Error("server error", "error", err)
	}
}
