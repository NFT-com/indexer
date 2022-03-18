package main

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/networks"
	"github.com/NFT-com/indexer/parsers"
	"github.com/NFT-com/indexer/workers/parsing"
)

const (
	EnvVarLogLevel = "LOG_LEVEL"

	DefaultLogLevel = "info"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("failure: %v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func run() error {
	logLevel, ok := os.LookupEnv(EnvVarLogLevel)
	if !ok {
		logLevel = DefaultLogLevel
	}

	// Logger initialization.
	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	log := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel)
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		return fmt.Errorf("failed to parse log level: %v", err)
	}
	log = log.Level(level)

	handler := parsing.NewHandler(log, func(_ networks.Network) (parsers.Parser, error) {
		parser := NewParser()

		return parser, nil
	})
	lambda.Start(handler.Handle)

	return nil
}
