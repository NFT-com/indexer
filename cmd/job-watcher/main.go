package main

import (
	"database/sql"
	"os"
	"os/signal"
	"time"

	_ "github.com/lib/pq"

	"github.com/adjust/rmq/v4"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"

	"github.com/NFT-com/indexer/config/params"
	"github.com/NFT-com/indexer/service/pipeline"
	"github.com/NFT-com/indexer/service/watcher"
	"github.com/NFT-com/indexer/storage/jobs"
)

const (
	success = 0
	failure = 1
)

func main() {
	os.Exit(run())
}

func run() int {

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	var (
		flagLogLevel string

		flagRedisDB  int
		flagRedisURL string
		flagJobDB    string

		flagOpenConnections uint
		flagIdleConnections uint
		flagReadInterval    time.Duration
	)

	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "severity level for log output")

	pflag.StringVarP(&flagJobDB, "job-database", "j", "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=postgres sslmode=disable", "postgresql connection details for job database")
	pflag.IntVarP(&flagRedisDB, "redis-database", "d", 1, "redis database number")
	pflag.StringVarP(&flagRedisURL, "redis-url", "u", "127.0.0.1:6379", "redis server url")

	pflag.UintVar(&flagOpenConnections, "db-connection-limit", 16, "maximum number of open database connections")
	pflag.UintVar(&flagIdleConnections, "db-idle-connection-limit", 4, "maximum number of idle database connections")
	pflag.DurationVar(&flagReadInterval, "read-interval", 100*time.Millisecond, "interval between checks for job reading")

	pflag.Parse()

	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	log := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel)
	level, err := zerolog.ParseLevel(flagLogLevel)
	if err != nil {
		log.Error().Err(err).Str("log_level", flagLogLevel).Msg("could not parse log level")
		return failure
	}
	log = log.Level(level)

	failed := make(chan error)
	redisConnection, err := rmq.OpenConnection(params.PipelineIndexer, "tcp", flagRedisURL, flagRedisDB, failed)
	if err != nil {
		log.Error().Err(err).Str("redis_url", flagRedisURL).Msg("could not connect to Redis")
		return failure
	}

	produce, err := pipeline.NewProducer(redisConnection, params.QueueParsing, params.QueueAction)
	if err != nil {
		log.Error().Err(err).Msg("could not create pipeline producer")
		return failure
	}

	jobDB, err := sql.Open(params.DialectPostgres, flagJobDB)
	if err != nil {
		log.Error().Err(err).Str("job_db", flagJobDB).Msg("could not open job database")
		return failure
	}
	jobDB.SetMaxOpenConns(int(flagOpenConnections))
	jobDB.SetMaxIdleConns(int(flagIdleConnections))

	parsingRepo := jobs.NewParsingRepository(jobDB)
	actionRepo := jobs.NewActionRepository(jobDB)

	watch := watcher.New(log, parsingRepo, actionRepo, produce, flagReadInterval)

	watch.Watch()

	log.Info().Msg("job watcher started")
	select {
	case <-sig:
		log.Info().Msg("job watcher stopping")
		watch.Close()
	case err = <-failed:
		log.Error().Err(err).Msg("job watcher aborted")
		return failure
	}

	go func() {
		<-sig
		log.Warn().Msg("forcing exit")
		os.Exit(1)
	}()

	log.Info().Msg("job watcher done")

	return success
}
