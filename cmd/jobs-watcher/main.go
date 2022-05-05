package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/adjust/rmq/v4"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"

	"github.com/NFT-com/indexer/config/params"
	"github.com/NFT-com/indexer/service/pipeline"
	"github.com/NFT-com/indexer/service/watcher"
	"github.com/NFT-com/indexer/storage/jobs"
)

func main() {
	err := run()
	if err != nil {
		// TODO: Improve this mixing logging
		// https://github.com/NFT-com/indexer/issues/32
		log.Fatal(err)
	}
}

func run() error {

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	var (
		flagLogLevel string

		flagJobsDB       string
		flagRedisNetwork string
		flagRedisURL     string
		flagRedisDB      int

		flagOpenConnections uint
		flagIdleConnections uint
		flagReadInterval    time.Duration
	)

	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "severity level for log output")

	pflag.StringVarP(&flagJobsDB, "jobs-database", "j", "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=jobs sslmode=disable", "Postgres connection details for jobs database")
	pflag.StringVarP(&flagRedisURL, "redis-url", "u", "127.0.0.1:6379", "URL for Redis server connection")
	pflag.IntVarP(&flagRedisDB, "redis-database", "d", 1, "Redis database number")

	pflag.UintVar(&flagOpenConnections, "db-connection-limit", 16, "maximum number of open database connections")
	pflag.UintVar(&flagIdleConnections, "db-idle-connection-limit", 4, "maximum number of idle database connections")
	pflag.DurationVar(&flagReadInterval, "read-interval", 200*time.Millisecond, "interval between checks for job reading")

	pflag.Parse()

	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	log := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel)
	level, err := zerolog.ParseLevel(flagLogLevel)
	if err != nil {
		return fmt.Errorf("could not parse log level: %w", err)
	}
	log = log.Level(level)

	failed := make(chan error)
	redisConnection, err := rmq.OpenConnection(params.PipelineIndexer, flagRedisNetwork, flagRedisURL, flagRedisDB, failed)
	if err != nil {
		return fmt.Errorf("could not open connection with redis: %w", err)
	}

	produce, err := pipeline.NewProducer(redisConnection, params.QueueParsing, params.QueueAction)
	if err != nil {
		return fmt.Errorf("could not create message producer: %w", err)
	}

	jobsDB, err := sql.Open(params.DialectPostgres, flagJobsDB)
	if err != nil {
		log.Error().Err(err).Msg("could not open SQL connection")
		return err
	}
	jobsDB.SetMaxOpenConns(int(flagOpenConnections))
	jobsDB.SetMaxIdleConns(int(flagIdleConnections))

	parsingRepo := jobs.NewParsingRepository(jobsDB)
	actionRepo := jobs.NewActionRepository(jobsDB)

	watch := watcher.New(log, parsingRepo, actionRepo, produce, flagReadInterval)

	log.Info().Msg("jobs watcher starting")
	watch.Watch()

	select {
	case <-sig:
		log.Info().Msg("jobs watcher stopping")
		watch.Close()
	case err = <-failed:
		log.Error().Err(err).Msg("jobs watcher aborted")
		return err
	}

	go func() {
		<-sig
		log.Warn().Msg("forcing exit")
		os.Exit(1)
	}()

	log.Info().Msg("jobs watcher done")

	return nil
}
