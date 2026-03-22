package main

import (
	"log/slog"
	"os"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/rpc"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/handlers"
	_ "modernc.org/sqlite"
)


func main() {
	server := rpc.NewServer()

	server.Register("initialize", handlers.InitializeHandler{}) 
	
	if err := server.Serve(os.Stdin, os.Stdout); err != nil {
		slog.Error("server error", "error", err)
	}
}
