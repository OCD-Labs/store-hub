package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/OCD-Labs/store-hub/api"
	"github.com/OCD-Labs/store-hub/util"
)

func main()  {
	configs, err := util.ParseConfigs(".")
	if err != nil {
		log.Fatal().Err(err).Msg("error occurred parsing configs")
	}

	if configs.Env == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Caller().Logger()
	}

	app := api.NewStoreHub(configs, log.Logger)
	if err = app.Start(); err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}