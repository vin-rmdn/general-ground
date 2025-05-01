package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/vin-rmdn/general-ground/cmd/server"
)

func main() {
	logger := slog.New(
		tint.NewHandler(os.Stdout, &tint.Options{
			AddSource:  true,
			Level:      slog.LevelDebug,
			TimeFormat: time.StampMilli,
			NoColor:    false,
		}),
	)

	slog.SetDefault(logger)

	instance, err := server.New()
	if err != nil {
		slog.Error("Failed to create server instance", "error", err)
		os.Exit(1)
	}

	if err := instance.Start(); err != nil {
		slog.Error("Failed to start server", "error", err)
		os.Exit(1)
}
}
