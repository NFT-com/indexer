package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/adjust/rmq/v4"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"

	"github.com/NFT-com/indexer/queue/consumer"
	"github.com/NFT-com/indexer/service/client"
)

const (
	defaultHTTPTimeout      = time.Second * 30
	defaultPollDuration     = time.Second * 20
	defaultParsingQueueName = "parsing"
)

func main() {
	err := run()
	if err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	// Signal catching for clean shutdown.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	// Command line parameter initialization.
	var (
		flagAPIEndpoint          string
		flagConsumerPrefetch     int64
		flagConsumerPollDuration time.Duration
		flagParsingQueueName     string
		flagLogLevel             string
		flagRMQTag               string
		flagRedisDatabase        int
		flagRedisNetwork         string
		flagRedisURL             string
	)

	pflag.StringVarP(&flagAPIEndpoint, "api", "a", "", "jobs api base hostname and port")
	pflag.Int64VarP(&flagConsumerPrefetch, "prefetch", "p", 5, "amount of message to prefetch in the consumer")
	pflag.DurationVarP(&flagConsumerPollDuration, "poll-duration", "i", defaultPollDuration, "time for each consumer poll")
	pflag.StringVarP(&flagParsingQueueName, "parsing-queue", "q", defaultParsingQueueName, "name of the queue for parsing")
	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "log level")
	pflag.StringVarP(&flagRMQTag, "tag", "c", "parsing-agent", "rmq consumer tag")
	pflag.IntVar(&flagRedisDatabase, "database", 1, "redis database number")
	pflag.StringVarP(&flagRedisNetwork, "network", "n", "tcp", "redis network type")
	pflag.StringVarP(&flagRedisURL, "url", "u", "", "redis server connection url")
	pflag.Parse()

	// Logger initialization.
	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	log := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel)
	level, err := zerolog.ParseLevel(flagLogLevel)
	if err != nil {
		return fmt.Errorf("could not parse log level: %w", err)
	}
	log = log.Level(level)

	cli := http.DefaultClient
	cli.Timeout = defaultHTTPTimeout

	api := client.New(log,
		client.WithHTTPClient(cli),
		client.WithHost(flagAPIEndpoint),
	)

	consumer := consumer.NewParsingConsumer(log, api)

	failed := make(chan error)

	redisClient := redis.NewClient(&redis.Options{
		Network: flagRedisNetwork,
		Addr:    flagRedisURL,
		DB:      flagRedisDatabase,
	})

	rmqConnection, err := rmq.OpenConnectionWithRedisClient(flagRMQTag, redisClient, failed)
	if err != nil {
		return fmt.Errorf("could not open redis connection: %w", err)
	}

	queue, err := rmqConnection.OpenQueue(flagParsingQueueName)
	if err != nil {
		return fmt.Errorf("could not open redis queue: %w", err)
	}

	err = queue.StartConsuming(flagConsumerPrefetch, flagConsumerPollDuration)
	if err != nil {
		return fmt.Errorf("could not start consuming process: %w", err)
	}

	consumerName, err := queue.AddConsumer(flagRMQTag, consumer)
	if err != nil {
		return fmt.Errorf("could not add parsing consumer: %w", err)
	}

	log.Info().Str("name", consumerName).Msg("started parsing dispatcher")

	select {
	case <-sig:
		rmqConnection.StopAllConsuming()
	case err := <-failed:
		return err
	}

	go func() {
		<-sig
		log.Fatal().Msg("forced interruption")
	}()

	log.Info().Str("name", consumerName).Msg("stopped parsing dispatcher gracefully")

	return nil
}
