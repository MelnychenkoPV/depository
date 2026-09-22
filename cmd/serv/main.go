package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"

	"github.com/MelnychenkoPV/dispensation/cmd"
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	cnf, err := cmd.ParseConfig()
	if err != nil {
		log.Fatal(err)
		return
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.InfoContext(ctx, "Starting application")

	if err := <-run(ctx, cnf, logger); err != nil {
		logger.ErrorContext(ctx, "run application error", slog.Any("err", err))
		os.Exit(1)
	}
	logger.InfoContext(ctx, "End application")
}
