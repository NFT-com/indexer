package main

import (
	"database/sql"
	"os"
	"os/signal"
	"time"

	_ "github.com/lib/pq"

	"github.com/nsqio/go-nsq"
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

		flagJobsDB    string
		flagNSQServer string

		flagOpenConnections uint
		flagIdleConnections uint
		flagReadInterval    time.Duration
	)

	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "severity level for log output")

	pflag.StringVarP(&flagJobsDB, "jobs-database", "j", "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=jobs sslmode=disable", "Postgres connection details for jobs database")
	pflag.StringVarP(&flagNSQServer, "nsq-server", "q", "127.0.0.1:4150", "address for NSQ server to produce messages")

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

	nsqCfg := nsq.NewConfig()
	producer, err := nsq.NewProducer(flagNSQServer, nsqCfg)
	if err != nil {
		log.Error().Err(err).Str("nsq_address", flagNSQServer).Msg("could not create NSQ producer")
		return failure
	}
	defer producer.Stop()

	handler, err := pipeline.NewJobCreator(producer, params.TopicParsing, params.TopicAction)
	if err != nil {
		log.Error().Err(err).Msg("could not create pipeline producer")
		return failure
	}

	jobsDB, err := sql.Open(params.DialectPostgres, flagJobsDB)
	if err != nil {
		log.Error().Err(err).Str("jobs_database", flagJobsDB).Msg("could not open jobs database")
		return failure
	}
	jobsDB.SetMaxOpenConns(int(flagOpenConnections))
	jobsDB.SetMaxIdleConns(int(flagIdleConnections))

	parsingRepo := jobs.NewParsingRepository(jobsDB)
	actionRepo := jobs.NewActionRepository(jobsDB)

	watch := watcher.New(log, parsingRepo, actionRepo, handler, flagReadInterval)
	defer watch.Close()

	watch.Watch()

	log.Info().Msg("jobs watcher started")

	<-sig

	log.Info().Msg("jobs watcher stopping")

	go func() {
		<-sig
		log.Warn().Msg("forcing exit")
		os.Exit(1)
	}()

	log.Info().Msg("jobs watcher done")

	return success
}
