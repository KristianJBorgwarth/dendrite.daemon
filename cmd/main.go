package main

import (
	"log/slog"
	"os"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/handlers"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/logging"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/server"
	_ "modernc.org/sqlite"
)

func main() {
	logging.Init()
	server := server.NewServer()

	server.RegisterHandler("initialize", handlers.NewInitializeHandler())
	server.RegisterHandler("create_note", handlers.NewCreateNoteHandler())

	if err := server.Run(os.Stdin, os.Stdout); err != nil {
		slog.Error("server error", "error", err)
	}
}
