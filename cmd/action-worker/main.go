package main

import (
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/credentials"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"

	"github.com/NFT-com/indexer/aws"
	"github.com/NFT-com/indexer/service/lambdas"
)

const (
	envLevel  = "LOG_LEVEL"
	envRegion = "AWS_REGION"

	defaultLevel  = "info"
	defaultRegion = "eu-west-1"

	service = "managedblockchain"
)

func main() {

	level, ok := os.LookupEnv(envLevel)
	if !ok {
		level = defaultLevel
	}

	awsRegion, ok := os.LookupEnv(envRegion)
	if !ok {
		awsRegion = defaultRegion
	}

	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	log := zerolog.New(os.Stderr).With().Timestamp().Logger()
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		log.Fatal().Err(err).Msg("could not parse log level")
	}
	log = log.Level(lvl)

	client := &http.Client{}

	creds := credentials.NewEnvCredentials()
	_, err = creds.Get()
	if err == nil {
		signer := v4.NewSigner(creds)
		transport := aws.NewInjectorRoundTripper(
			signer,
			awsRegion,
			service,
			http.DefaultTransport,
		)
		client.Transport = transport
	}

	handler := lambdas.NewAdditionHandler(log, client)
	lambda.Start(handler.Handle)

	os.Exit(0)
}
