package api

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/OCD-Labs/store-hub/util"
	"github.com/rs/zerolog"
)

type StoreHub struct {
	configs util.Configs
	logger  zerolog.Logger
}

func NewStoreHub(configs util.Configs, logger zerolog.Logger) *StoreHub {
	return &StoreHub{ configs, logger }
}

func (s *StoreHub) Start() error {
	srv := http.Server{
		Addr:         s.configs.ServerAddr,
		Handler:      s.setupRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		ErrorLog:     log.New(s.logger, "", 0),
	}

	recvErr := make(chan error)
	go func() {
		recvSig := make(chan os.Signal, 1)

		signal.Notify(recvSig, syscall.SIGINT, syscall.SIGTERM)
		sig := <-recvSig

		s.logger.Info().Str("signal received", sig.String()).Msg("shutdown started")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			recvErr <- err
		}

		s.logger.Info().Str("server address", srv.Addr).Msg("completing background tasks: ")

		recvErr <- nil
	}()

	s.logger.Info().
		Str("server address", srv.Addr).
		Str("environment", s.configs.Env).
		Msg("starting server...")

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <- recvErr
	if err != nil {
		return err
	}

	s.logger.Info().Msg("server shutdown completed")

	return nil
}