package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/rs/zerolog"

	"github.com/aws/aws-lambda-go/lambda"

	handler "github.com/NFT-com/indexer/function/handlers/parsing"
	"github.com/NFT-com/indexer/function/processors/parsing"
	"github.com/NFT-com/indexer/function/processors/parsing/erc721"
	"github.com/NFT-com/indexer/function/processors/parsing/opensea"
	"github.com/NFT-com/indexer/networks"
)

const (
	envVarLogLevel = "LOG_LEVEL"

	defaultLogLevel = "info"
)

func main() {
	err := run()
	if err != nil {
		// TODO: Improve this mixing logging
		// https://github.com/NFT-com/indexer/issues/32
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

	handler := handler.NewHandler(log, func(client networks.Network) ([]parsing.Parser, error) {
		parsers := make([]parsing.Parser, 0, 2)
		parsers = append(parsers, erc721.NewParser())

		openseaParser, err := opensea.NewParser(client)
		if err != nil {
			return nil, fmt.Errorf("could not create opensea parser: %w", err)
		}
		parsers = append(parsers, openseaParser)

		return parsers, nil
	})
	lambda.Start(handler.Handle)

	return nil
}
