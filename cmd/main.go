package main

import (
	"github.com/KristianJBorgwarth/dendrite.daemon/core/http"
	"github.com/KristianJBorgwarth/dendrite.daemon/core/logging"
)

func main() {
	logging.Init()
	srv := server.New(":6969")
	srv.Start()
}
