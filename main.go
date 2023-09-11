package main

import (
	"database/sql"
	"embed"
	"flag"
	"io/fs"
	"os"
	"time"

	"github.com/OCD-Labs/store-hub/api"
	"github.com/OCD-Labs/store-hub/cache"
	db "github.com/OCD-Labs/store-hub/db/sqlc"
	"github.com/OCD-Labs/store-hub/mailer"
	"github.com/OCD-Labs/store-hub/near"
	"github.com/OCD-Labs/store-hub/token"
	"github.com/OCD-Labs/store-hub/util"
	"github.com/OCD-Labs/store-hub/worker"
	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//go:embed doc
var swaggerDocs embed.FS

func main() {
	configs, err := util.ParseConfigs(".")
	if err != nil {
		log.Fatal().Err(err).Msg("error occurred parsing configs")
	}

	flag.Float64Var(&configs.Limiter.RPS, "limiter-rps", 3, "Rate limiter maximum requests per second")
	flag.IntVar(&configs.Limiter.Burst, "limiter-burst", 5, "Rate limiter maximum burst")
	flag.BoolVar(&configs.Limiter.Enabled, "limiter-enabled", true, "Enable rate limiter")
	flag.Parse()

	if configs.Env == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Caller().Logger()
	}

	log.Info().Msg("connecting to DB")
	dbConn, err := sql.Open(configs.DBDriver, configs.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to DB")
	}

	runDBMigrations(configs.MigrationURL, configs.DBSource, configs)
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

	// Retrieve the swagger-ui files.
	swaggerFiles, err := fs.Sub(swaggerDocs, "doc/swagger")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get subcontent from swaggerDocs")
	}

	app, err := api.NewStoreHub(configs, log.Logger, cache, dbStore, taskDistributor, tokenMaker, swaggerFiles)
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

func runDBMigrations(migrationURL string, dbSource string, configs util.Configs) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create a new migrate instance")
	}

	if configs.Env == "production" {
		if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal().Err(err).Msg("failed to run migrateup")
		}
	}

	log.Info().Msg("db migrated successfully")
}


func setupNEAR(configs util.Configs) (err error) {
	if util.CommandExists("npm") {
		if err = near.InstallNearCLI(); err != nil {
			return err
		}

		if err = near.SetupNearMasterAccount(configs.NEARNetwork, configs.NEARAccountID, configs.NEARPubKey, configs.NEARPrivKey); err != nil {
			return err
		}
	}

	return nil
}