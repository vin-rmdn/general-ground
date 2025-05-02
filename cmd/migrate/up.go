package migrate

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // import required for postgres driver initialization
	_ "github.com/golang-migrate/migrate/v4/source/file"       // import required for file source initialization
	"github.com/spf13/viper"
	"github.com/urfave/cli/v3"
	"github.com/vin-rmdn/general-ground/internal/version"
)

var upCommand = &cli.Command{
	Name:                  "up",
	Aliases:               []string{"u", "upgrade", "^"},
	Usage:                 "Upgrade the database to the latest version",
	UsageText:             "general_ground migrate up [options]",
	ArgsUsage:             "argsusage",
	Version:               version.Version,
	Description:           "Migrates database to the latest version to be used with the server",
	DefaultCommand:        "defaultcommand",
	Category:              "database",
	EnableShellCompletion: true,
	Before:                setupMigrate,
	After:                 nil,
	Action:                migrateUp,
	Authors:               []any{"vin-rmdn"},
	Suggest:               true,
}

func migrateUp(ctx context.Context, c *cli.Command) error {
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

	slog.Debug("Starting migration process")

	if err = migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		slog.Error("Failed to apply migrations", "error", err)

		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		slog.Info("No new migrations to apply")
	}

	slog.Info("Migration process completed successfully")

	return nil
}
