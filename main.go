package main

import (
	"context"
	"os"

	"github.com/urfave/cli/v3"
	"github.com/vin-rmdn/general-ground/cmd/server"
	"github.com/vin-rmdn/general-ground/internal/version"
)

func main() {
	cmd := &cli.Command{
		Name:                  "general_ground",
		Aliases:               []string{"serve"},
		Usage:                 "Service for a mock chatting platform",
		UsageText:             "general_ground {command} [options]",
		ArgsUsage:             "argsusage",
		Version:               version.Version,
		Description:           "Service for a mock chatting platform, complete with a server and a database migration tool, and many other features to come.",
		DefaultCommand:        "defaultcommand",
		Category:              "category",
		Commands:              []*cli.Command{server.Command},
		Flags:                 []cli.Flag{},
		EnableShellCompletion: true,
		Authors:               []any{"vin-rmdn"},
		Suggest:               true,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		os.Exit(1)
	}
}
