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

	server.NewServer()
}
