// Package logging
package logging

import (
	"log/slog"
	"os"
)

func Init() {
	logFile, err := os.OpenFile("/tmp/dendrite.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		panic(err)
	}

	handler := slog.NewTextHandler(logFile, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	logger := slog.New(handler)

	slog.SetDefault(logger)
}
