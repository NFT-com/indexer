package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go/aws/credentials"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"

	"github.com/NFT-com/indexer/config/params"
	"github.com/NFT-com/indexer/network/amb"
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

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("could not load AWS config")
	}

	client := &http.Client{}
	creds := credentials.NewEnvCredentials()
	_, err = creds.Get()
	if err == nil {
		signer := v4.NewSigner(creds)
		transport := amb.NewRoundTripper(
			signer,
			cfg.Region,
			params.AWSManagedBlockchain,
			http.DefaultTransport,
		)
		client.Transport = transport
	}

	handler := lambdas.NewAdditionHandler(log, client)
	lambda.Start(handler.Handle)

	os.Exit(0)
}
