package main

import (
  "github.com/KristianJBorgwarth/dendrite.daemon/core"
	_ "modernc.org/sqlite"
)

func main() {
	server := core.NewServer()
	server.Run()
}

