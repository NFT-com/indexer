package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/rs/zerolog"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/NFT-com/indexer/function/handlers/action"
	processor "github.com/NFT-com/indexer/function/processors/action"
	"github.com/NFT-com/indexer/function/processors/action/erc721/owner-change"
	"github.com/NFT-com/indexer/function/processors/action/erc721metadata/addition"
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

	handler := action.NewHandler(log, func(client networks.Network) ([]processor.Processor, error) {
		processors := make([]processor.Processor, 0, 0)

		erc721Processor, err := addition.NewProcessor(log, client)
		if err != nil {
			return nil, fmt.Errorf("could not create addition erc721 metadata processor: %w", err)
		}
		processors = append(processors, erc721Processor)

		ownerChangeProcessor, err := owner_change.NewProcessor(log, client)
		if err != nil {
			return nil, fmt.Errorf("could not create owner change erc721 processor: %w", err)
		}
		processors = append(processors, ownerChangeProcessor)

		return processors, nil
	})
	lambda.Start(handler.Handle)

	return nil
}
