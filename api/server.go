package api

import (
	"context"
	"errors"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/OCD-Labs/store-hub/cache"
	db "github.com/OCD-Labs/store-hub/db/sqlc"
	"github.com/OCD-Labs/store-hub/token"
	"github.com/OCD-Labs/store-hub/util"
	"github.com/OCD-Labs/store-hub/worker"
	"github.com/rs/zerolog"
)

type StoreHub struct {
	configs                util.Configs
	logger                 zerolog.Logger
	swaggerFiles           fs.FS
	cache                  cache.Cache
	tokenMaker             token.Maker
	dbStore                db.StoreTx
	taskDistributor        worker.TaskDistributor
	SupportUnauthenticated bool
}

func NewStoreHub(
	configs util.Configs,
	logger zerolog.Logger,
	cache cache.Cache,
	store db.StoreTx,
	taskDistributor worker.TaskDistributor,
	tokenMaker token.Maker,
	swaggerFiles fs.FS,
) (*StoreHub, error) {
	return &StoreHub{
		configs:         configs,
		logger:          logger,
		cache:           cache,
		tokenMaker:      tokenMaker,
		dbStore:         store,
		taskDistributor: taskDistributor,
		swaggerFiles:    swaggerFiles,
	}, nil
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

	err = <-recvErr
	if err != nil {
		return err
	}

	s.logger.Info().Msg("server shutdown completed")

	return nil
}
