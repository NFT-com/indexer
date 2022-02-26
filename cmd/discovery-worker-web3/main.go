package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/NFT-com/indexer/handler/discovery/web3"
	"github.com/NFT-com/indexer/queue/producer"
	"github.com/adjust/rmq/v4"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog"
)

const (
	RMQTagEnvVar              = "RMQ_TAG"
	RedisNetworkEnvVar        = "REDIS_NETWORK"
	RedisURLEnvVar            = "REDIS_URL"
	RedisDatabaseEnvVar       = "REDIS_DATABASE"
	RedisParseQueueNameEnvVar = "REDIS_PARSE_QUEUE_NAME"
	LogLevelEnvVar            = "LOG_LEVEL"

	DefaultRMQTag              = "discovery-worker-web3"
	DefaultRedisNetwork        = "tcp"
	DefaultRedisDatabase       = "1"
	DefaultRedisParseQueueName = "parse"
	DefaultLogLevel            = "info"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("failure: %v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func run() error {
	logLevel, ok := os.LookupEnv(LogLevelEnvVar)
	if !ok {
		logLevel = DefaultLogLevel
	}

	// Logger initialization.
	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	log := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel)
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		return err
	}
	log = log.Level(level)
	log.Info().Msg("asd")

	rmqTag, ok := os.LookupEnv(RMQTagEnvVar)
	if !ok {
		rmqTag = DefaultRMQTag
	}

	redisNetwork, ok := os.LookupEnv(RedisNetworkEnvVar)
	if !ok {
		redisNetwork = DefaultRedisNetwork
	}

	redisURL, ok := os.LookupEnv(RedisURLEnvVar)
	if !ok {
		return errors.New("missing redis url")
	}

	redisDatabaseString, ok := os.LookupEnv(RedisDatabaseEnvVar)
	if !ok {
		redisDatabaseString = DefaultRedisDatabase
	}

	redisParseQueueName, ok := os.LookupEnv(RedisParseQueueNameEnvVar)
	if !ok {
		redisParseQueueName = DefaultRedisParseQueueName
	}

	redisDatabase, err := strconv.Atoi(redisDatabaseString)
	if err != nil {
		return err
	}

	failed := make(chan error)
	connection, err := rmq.OpenConnection(rmqTag, redisNetwork, redisURL, redisDatabase, failed)
	if err != nil {
		return err
	}

	prod, err := producer.NewProducer(connection)
	if err != nil {
		return err
	}

	handler, err := web3.NewWeb3(log, redisParseQueueName, prod)
	if err != nil {
		return err
	}

	lambda.Start(handler.Handle)

	return nil
}
