package main

import (
	"fmt"
	"os"
	"time"

	parseDefault "github.com/NFT-com/indexer/handler/parse/default"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog"
)

const (
	LogLevelEnvVar = "LOG_LEVEL"

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

	handler, err := parseDefault.NewDefault(log)
	if err != nil {
		return err
	}

	lambda.Start(handler.Handle)

	return nil
}
