package migrate

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/urfave/cli/v3"
	"github.com/vin-rmdn/general-ground/internal/config"
	"github.com/vin-rmdn/general-ground/internal/logger"
	"github.com/vin-rmdn/general-ground/internal/version"
)

var RootCommand = &cli.Command{
	Name:                  "migrate",
	Aliases:               []string{"m"},
	Usage:                 "Migrates database",
	UsageText:             "general_ground migrate [options]",
	ArgsUsage:             "argsusage",
	Version:               version.Version,
	Description:           "Migrates database to be used with the server",
	Category:              "database",
	Commands:              []*cli.Command{upCommand, downCommand},
	EnableShellCompletion: true,
	Before:                setupMigrate,
	After:                 nil,
	Action:                nil,
	Authors:               []any{"vin-rmdn"},
	Suggest:               true,
}

func setupMigrate(ctx context.Context, _ *cli.Command) (context.Context, error) {
	if err := config.SetupEnvironment(); err != nil {
		slog.Error("Failed to setup environment", "error", err)

		return nil, fmt.Errorf("failed to setup environment: %w", err)
	}

	logger.Setup()

	return ctx, nil
}
