package main

import (
	"database/sql"
	"os"
	"time"

	"github.com/OCD-Labs/store-hub/api"
	"github.com/OCD-Labs/store-hub/cache"
	db "github.com/OCD-Labs/store-hub/db/sqlc"
	"github.com/OCD-Labs/store-hub/util"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	configs, err := util.ParseConfigs(".")
	if err != nil {
		log.Fatal().Err(err).Msg("error occurred parsing configs")
	}

	if configs.Env == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Caller().Logger()
	}

	log.Info().Msg("connecting to DB")
	dbConn, err := sql.Open(configs.DBDriver, configs.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to DB")
	}

	store := db.NewSQLTx(dbConn)

	cache, err := cache.NewRedisCache(configs.RedisAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to redis cache")
	}
	log.Info().Msg("redis cache connection established")

	app, err := api.NewStoreHub(configs, log.Logger, cache, store)
	if err = app.Start(); err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}