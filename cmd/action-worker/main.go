package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/NFT-com/indexer/service/lambdas"
)

const (
	defaultLevel = "info"
	envLevel     = "LOG_LEVEL"
)

func main() {

	level, ok := os.LookupEnv(envLevel)
	if !ok {
		level = defaultLevel
	}

	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	log := zerolog.New(os.Stderr).With().Timestamp().Logger()
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		log.Fatal().Err(err).Msg("could not parse log level")
	}
	log = log.Level(lvl)

	handler := lambdas.NewActionHandler(log)
	lambda.Start(handler.Handle)

	os.Exit(0)
}
