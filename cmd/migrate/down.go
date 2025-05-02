package migrate

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v3"
	"github.com/vin-rmdn/general-ground/internal/version"
)

var downCommand = &cli.Command{
	Name:                  "down",
	Aliases:               []string{"d", "downgrade", "_"},
	Usage:                 "Remove all applied migrations",
	UsageText:             "general_ground migrate down [options]",
	ArgsUsage:             "argsusage",
	Version:               version.Version,
	Description:           "Migrates database to the latest version to be used with the server",
	DefaultCommand:        "defaultcommand",
	Category:              "database",
	EnableShellCompletion: true,
	Before:                setupMigrate,
	After:                 nil,
	Action:                migrateDown,
	Authors:               []any{"vin-rmdn"},
	Suggest:               true,
}

func migrateDown(ctx context.Context, c *cli.Command) error {
	const postgresJdbcUrlTemplate = "postgres://%s:%s@%s:%d/%s?sslmode=disable"
	jdbcUrl := fmt.Sprintf(
		postgresJdbcUrlTemplate,
		viper.GetString("POSTGRESQL_USER"),
		viper.GetString("POSTGRESQL_PASSWORD"),
		viper.GetString("POSTGRESQL_HOST"),
		viper.GetUint16("POSTGRESQL_PORT"),
		viper.GetString("POSTGRESQL_DATABASE_NAME"),
	)

	migrator, err := migrate.New("file://database/postgresql/migrations", jdbcUrl)
	if err != nil {
		slog.Error("Failed to create migrator", "error", err)

		return fmt.Errorf("failed to create migrator: %w", err)
	}

	slog.Debug("Starting migration downgrade process")

	err = migrator.Down()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		slog.Error("Failed to downgrade migrations", "error", err)

		return fmt.Errorf("failed to downgrade migrations: %w", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		slog.Info("No migrations to downgrade")
	}

	slog.Debug("Migration downgrade process completed successfully")

	return nil
}
