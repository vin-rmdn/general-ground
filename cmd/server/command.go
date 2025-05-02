package server

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v3"
	"github.com/vin-rmdn/general-ground/internal/config"
	"github.com/vin-rmdn/general-ground/internal/version"
)

var Command = &cli.Command{
	Name:                  "server",
	Aliases:               []string{"serve", "s"},
	Usage:                 "Start a chat server",
	UsageText:             "general_ground server [options]",
	ArgsUsage:             "argsusage",
	Version:               version.Version,
	Description:           "Start a chat server",
	DefaultCommand:        "defaultcommand",
	Category:              "service",
	Flags:                 []cli.Flag{},
	EnableShellCompletion: true,
	Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
		if err := config.SetupEnvironment(); err != nil {
			slog.Error("Failed to setup environment", "error", err)

			return nil, fmt.Errorf("failed to setup environment: %w", err)
		}

		setupLogger()

		return ctx, nil
	},
	After:   nil, // TODO: add cleanup function
	Action:  execute,
	Authors: []any{"vin-rmdn"},
	Suggest: true,
}

func execute(ctx context.Context, c *cli.Command) error {
	certificateKeyPath := viper.GetString("CERTIFICATE_KEY_PATH")
	certificatePath := viper.GetString("CERTIFICATE_PATH")

	instance, err := New(certificatePath, certificateKeyPath)
	if err != nil {
		slog.Error("Failed to create server instance", "error", err)

		return fmt.Errorf("failed to create server instance: %w", err)
	}

	if err := instance.Start(); err != nil {
		slog.Error("Failed to start server", "error", err)

		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func setupLogger() {
	slogHandler := tint.NewHandler(os.Stdout, &tint.Options{
		AddSource:  true,
		Level:      slog.LevelDebug,
		TimeFormat: time.StampMilli,
	})
	logger := slog.New(slogHandler)

	slog.SetDefault(logger)
}
