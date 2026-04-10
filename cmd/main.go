package main

import (
	"log/slog"
	"os"

	"github.com/KristianJBorgwarth/dendrite.daemon/core/handlers/completion"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/handlers/note"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/handlers/vault"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/logging"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/server"
	_ "modernc.org/sqlite"
)

func main() {
	logging.Init()
	server := server.NewServer()

	server.RegisterHandler("initialize", vault.NewInitializeHandler())
	server.RegisterHandler("create_note", note.NewCreateNoteHandler())
	server.RegisterHandler("completion/link", completion.NewCompleteLinkHandler())

	if err := server.Run(os.Stdin, os.Stdout); err != nil {
		slog.Error("server error", "error", err)
	}
}
