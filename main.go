package main

import (
	"log/slog"
	"os"

	"github.com/vin-rmdn/general-ground/cmd/server"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))

	slog.SetDefault(logger)

	server.NewServer()
}