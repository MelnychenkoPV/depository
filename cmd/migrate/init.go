package main

import (
	"context"
	"errors"
	"log/slog"

	"github.com/MelnychenkoPV/dispensation/cmd"
	"github.com/MelnychenkoPV/dispensation/migration"
	"github.com/golang-migrate/migrate/v4"
	"github.com/spf13/cobra"
)

func run(ctx context.Context, cnf cmd.Config, logger *slog.Logger) chan error {

	outCh := make(chan error)

	go func() {
		defer close(outCh)

		db, err := cmd.CreateDB(ctx, cnf.DB)
		if err != nil {
			outCh <- err
			return
		}

		m, err := migration.NewMigration(db)
		if err != nil {
			outCh <- err
			return
		}
		defer m.Close()

		cmdCh := make(chan error)
		go func() {
			defer close(cmdCh)

			if err := createRootCMD(m, logger).ExecuteContext(ctx); err != nil {
				outCh <- err
				return
			}

			version, dirty, err := m.Version()
			if err != nil {
				outCh <- err
				return
			}
			logger.InfoContext(ctx, "VERSION", slog.Any("version", version), slog.Any("dirty", dirty))
		}()

		select {
		case cmdErr := <-cmdCh:
			outCh <- cmdErr
		case <-ctx.Done():
			srcErr, dbErr := m.Close()
			if srcErr != nil {
				outCh <- srcErr
			}
			if dbErr != nil {
				outCh <- dbErr
			}
			if err := db.Close(); err != nil {
				outCh <- err
			}
		}
	}()

	return outCh
}

func createRootCMD(m *migrate.Migrate, logger *slog.Logger) *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "migration",
		Short: "Migration commands",
	}

	rootCmd.AddCommand(
		&cobra.Command{
			Use: "up",
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := m.Up(); err != nil {
					if errors.Is(err, migrate.ErrNoChange) {
						logger.InfoContext(cmd.Context(), "UP no change")
						return nil
					}
					return err
				}
				logger.InfoContext(cmd.Context(), "UP complete")
				return nil
			},
		},
		&cobra.Command{
			Use: "down",
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := m.Down(); err != nil {
					if errors.Is(err, migrate.ErrNoChange) {
						logger.InfoContext(cmd.Context(), "DOWN no change")
						return nil
					}
					return err
				}
				logger.InfoContext(cmd.Context(), "Down complete")
				return nil
			},
		},
		&cobra.Command{
			Use:   "version",
			Short: "Version of database migration",
			RunE: func(cmd *cobra.Command, args []string) error {
				version, dirty, err := m.Version()
				if err != nil {
					return err
				}
				logger.InfoContext(cmd.Context(), "VERSION", slog.Any("version", version), slog.Any("dirty", dirty))
				return nil
			},
		},
	)

	return rootCmd
}
