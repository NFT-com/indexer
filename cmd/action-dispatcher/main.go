package main

import (
	"context"
	"database/sql"
	"os"
	"os/signal"
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"
	"go.uber.org/ratelimit"

	_ "github.com/lib/pq"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"

	"github.com/NFT-com/indexer/config/params"
	"github.com/NFT-com/indexer/service/pipeline"
	"github.com/NFT-com/indexer/storage/graph"
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

		flagGraphDB    string
		flagJobsDB     string
		flagNSQLookups []string
		flagLambdaName string

		flagOpenConnections   uint
		flagIdleConnections   uint
		flagRateLimit         uint
		flagLambdaConcurrency uint

		flagDryRun bool
	)

	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "severity level for log output")

	pflag.StringVarP(&flagGraphDB, "graph-database", "g", "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=graph sslmode=disable", "Postgres connection details for graph database")
	pflag.StringVarP(&flagJobsDB, "jobs-database", "j", "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=jobs sslmode=disable", "Postgres connection details for jobs database")
	pflag.StringSliceVarP(&flagNSQLookups, "nsq-lookups", "", []string{"127.0.0.1:4161"}, "addresses for NSQ lookups to bootstrap consuming")
	pflag.StringVarP(&flagLambdaName, "lambda-name", "n", "action-worker", "name of the Lambda function to invoke")

	pflag.UintVar(&flagOpenConnections, "db-connection-limit", 128, "maximum number of database connections, -1 for unlimited")
	pflag.UintVar(&flagIdleConnections, "db-idle-connection-limit", 32, "maximum number of idle connections")
	pflag.UintVar(&flagRateLimit, "rate-limit", 100, "maximum number of API requests per second")
	pflag.UintVar(&flagLambdaConcurrency, "lambda-concurrency", 900, "maximum number of concurrent Lambda invocations")

	pflag.BoolVar(&flagDryRun, "dry-run", false, "executing as dry run disables invocation of Lambda function")

	pflag.Parse()

	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	log := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel)
	level, err := zerolog.ParseLevel(flagLogLevel)
	if err != nil {
		log.Error().Err(err).Str("log_level", flagLogLevel).Msg("could not parse log level")
		return failure
	}
	log = log.Level(level)

	awsCfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("could not load AWS configuration")
		return failure
	}

	jobsDB, err := sql.Open(params.DialectPostgres, flagJobsDB)
	if err != nil {
		log.Error().Err(err).Str("jobs_database", flagJobsDB).Msg("could not connect to jobs database")
		return failure
	}
	jobsDB.SetMaxOpenConns(int(flagOpenConnections))
	jobsDB.SetMaxIdleConns(int(flagIdleConnections))

	actionRepo := jobs.NewActionRepository(jobsDB)

	graphDB, err := sql.Open(params.DialectPostgres, flagGraphDB)
	if err != nil {
		log.Error().Err(err).Str("graph_database", flagGraphDB).Msg("could not connect to graph database")
		return failure
	}
	graphDB.SetMaxOpenConns(int(flagOpenConnections))
	graphDB.SetMaxIdleConns(int(flagIdleConnections))

	collectionRepo := graph.NewCollectionRepository(graphDB)
	nftRepo := graph.NewNFTRepository(graphDB)
	ownerRepo := graph.NewOwnerRepository(graphDB)
	traitRepo := graph.NewTraitRepository(graphDB)

	nsqCfg := nsq.NewConfig()
	err = nsqCfg.Set("max-in-flight", 2*flagLambdaConcurrency)
	if err != nil {
		log.Error().Err(err).Msg("could not update NSQ configuration")
		return failure
	}

	consumer, err := nsq.NewConsumer(params.TopicAction, params.PipelineIndexer, nsqCfg)
	if err != nil {
		log.Error().Err(err).Str("topic", params.TopicAction).Str("channel", params.PipelineIndexer).Msg("could not create NSQ consumer")
		return failure
	}
	defer consumer.Stop()

	lambda := lambda.NewFromConfig(awsCfg)
	limit := ratelimit.New(int(flagRateLimit))
	handler := pipeline.NewActionHandler(context.Background(), log, lambda, flagLambdaName, actionRepo, collectionRepo, nftRepo, ownerRepo, traitRepo, limit, flagDryRun)
	consumer.AddConcurrentHandlers(handler, int(flagLambdaConcurrency))

	err = consumer.ConnectToNSQLookupds(flagNSQLookups)
	if err != nil {
		log.Error().Err(err).Strs("nsq_lookups", flagNSQLookups).Msg("could not connect to NSQ lookups")
		return failure
	}

	log.Info().Msg("action dispatcher started")

	<-sig

	log.Info().Msg("initialized shutdown")

	go func() {
		<-sig
		log.Fatal().Msg("forced shutdown")
	}()

	log.Info().Msg("shutdown complete")

	return success
}
