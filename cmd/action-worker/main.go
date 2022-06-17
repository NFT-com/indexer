package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/NFT-com/indexer/service/lambdas"
)

func main() {

	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	log := zerolog.New(os.Stderr).With().Timestamp().Logger()

	envLevel, ok := os.LookupEnv("LOG_LEVEL")
	if ok {
		level, err := zerolog.ParseLevel(envLevel)
		if err != nil {
			log.Warn().Str("LOG_LEVEL", envLevel).Msg("invalid log level, using default")
		} else {
			log = log.Level(level)
		}
	}

	envNodeURL, ok := os.LookupEnv("NODE_URL")
	if !ok {
		log.Fatal().Msg("missing node URL, aborting execution")
	}

	handler := lambdas.NewAdditionHandler(log, envNodeURL)
	lambda.Start(handler.Handle)

	os.Exit(0)
}
