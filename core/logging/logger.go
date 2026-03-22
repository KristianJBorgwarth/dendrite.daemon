// Package logging
package logging

import (
	"io"
	"log/slog"
	"os"
)

func Init() {
	logFile, err := os.OpenFile("/tmp/dendrite.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		panic(err)
	}

	writer := io.MultiWriter(os.Stdout, logFile)

	handler := slog.NewTextHandler(writer, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	logger := slog.New(handler)

	slog.SetDefault(logger)
}
