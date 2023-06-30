package main

import (
	"database/sql"
	"os"
	"time"

	"github.com/OCD-Labs/store-hub/api"
	"github.com/OCD-Labs/store-hub/cache"
	db "github.com/OCD-Labs/store-hub/db/sqlc"
	"github.com/OCD-Labs/store-hub/mailer"
	"github.com/OCD-Labs/store-hub/token"
	"github.com/OCD-Labs/store-hub/util"
	"github.com/OCD-Labs/store-hub/worker"
	"github.com/hibiken/asynq"
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

	dbStore := db.NewSQLTx(dbConn)

	cache, err := cache.NewRedisCache(configs.RedisAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to redis cache")
	}
	log.Info().Msg("redis cache connection established")

	redisOpt := asynq.RedisClientOpt{
		Addr: configs.RedisAddress,
	}
	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	tokenMaker, err := token.NewPasetoMaker(configs.TokenSymmetricKey)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create token maker")
	}

	app, err := api.NewStoreHub(configs, log.Logger, cache, dbStore, taskDistributor, tokenMaker)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialise application")
	}

	log.Info().Msg("starting redis server")
	go runTaskProcessor(configs, redisOpt, dbStore, tokenMaker)

	if err = app.Start(); err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}

func runTaskProcessor(config util.Configs, redisOpt asynq.RedisClientOpt, store db.StoreTx, tokenMaker token.Maker) {
	mailer := mailer.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mailer, config, tokenMaker)
	log.Info().Msg("starting task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}
}