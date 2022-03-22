package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/networks"
	"github.com/NFT-com/indexer/parsers"
	"github.com/NFT-com/indexer/workers/parsing"
)

const (
	envVarLogLevel = "LOG_LEVEL"

	defaultLogLevel = "info"
)

func main() {
	err := run()
	if err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	logLevel, ok := os.LookupEnv(envVarLogLevel)
	if !ok {
		logLevel = defaultLogLevel
	}

	// Logger initialization.
	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	log := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel)
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		return fmt.Errorf("could not parse log level: %w", err)
	}
	log = log.Level(level)

	handler := parsing.NewHandler(log, func(client networks.Network) (parsers.Parser, error) {
		parser, err := NewParser(client)
		if err != nil {
			return nil, fmt.Errorf("could not create parser client: %w", err)
		}

		return parser, nil
	})
	lambda.Start(handler.Handle)

	return nil
}
