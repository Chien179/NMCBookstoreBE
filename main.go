package main

import (
	"crypto/tls"
	"database/sql"
	"os"

	"github.com/Chien179/NMCBookstoreBE/src/api"
	db "github.com/Chien179/NMCBookstoreBE/src/db/sqlc"
	"github.com/Chien179/NMCBookstoreBE/src/mail"
	"github.com/Chien179/NMCBookstoreBE/src/util"
	"github.com/Chien179/NMCBookstoreBE/src/worker"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/hibiken/asynq"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}

	store := db.NewStore(conn)

	redisOpt := asynq.RedisClientOpt{
		Addr:      config.RedisAddress,
		Username:  "red-clmdiicjtl8s73aiv8tg",
		Password:  "iKGnFsunMjhr8iQQl4In5A98RjxjRlEK",
		TLSConfig: &tls.Config{MinVersion: tls.VersionTLS12},
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)
	go runTaskProcessor(config, redisOpt, store)
	runGinServer(config, store, taskDistributor)
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create new migrate instance")
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("failed to run migrate up")
	}

	log.Info().Msgf("db migrated successfully")
}

func runTaskProcessor(config util.Config, redisOpt asynq.RedisClientOpt, store db.Store) {
	mailer := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mailer)
	log.Info().Msg("start task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}
}

func connectElasticSearch(config util.Config) (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{config.ELASTIC_ADDRESS},
	}
	es, err := elasticsearch.NewClient(cfg)
	log.Info().Msg("start elastic search")

	return es, err
}

func runGinServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	elastic, err := connectElasticSearch(config)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to elastic")
	}

	server, err := api.NewServer(config, store, elastic, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start server")
	}
}
