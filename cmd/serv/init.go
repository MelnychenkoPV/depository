package main

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/MelnychenkoPV/dispensation/cmd"
	//"github.com/MelnychenkoPV/dispensation/internal/handler"
	//"github.com/go-playground/validator/v10"
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

		//validate := validator.New(validator.WithRequiredStructEnabled())

		//goodsRepo := repository.NewGoodsRepository(db)
		//
		//goodsSrv := service.NewGoodsService(goodsRepo)

		srv := &http.Server{
			Addr: cnf.Server.Addr,
			//Handler: handler.CreateHandler(goodsSrv, validate, logger),
		}

		srvCh := make(chan error)
		go func() {
			defer close(srvCh)

			if err := srv.ListenAndServe(); err != nil {
				srvCh <- err
				return
			}
		}()

		select {
		case srvErr := <-srvCh:
			outCh <- srvErr
		case <-ctx.Done():
			stopCtx := context.TODO()
			if err := srv.Shutdown(stopCtx); err != nil {
				outCh <- err
			}
			if err := db.Close(); err != nil {
				outCh <- err
			}
		}
	}()

	return outCh
}
