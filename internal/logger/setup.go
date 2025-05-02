package logger

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

func Setup() {
	slogHandler := tint.NewHandler(os.Stdout, &tint.Options{
		AddSource:  true,
		Level:      slog.LevelDebug,
		TimeFormat: time.StampMilli,
	})
	logger := slog.New(slogHandler)

	slog.SetDefault(logger)
}
