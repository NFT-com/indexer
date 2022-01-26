package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"

	dispatcher "github.com/NFT-com/indexer/dispatch/aws"
	"github.com/NFT-com/indexer/event"
	"github.com/NFT-com/indexer/networks/ethereum"
	"github.com/NFT-com/indexer/source"
	"github.com/NFT-com/indexer/store/mock"
	"github.com/NFT-com/indexer/subscriber"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("failure: %v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Signal catching for clean shutdown.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	// Command line parameter initialization.
	var (
		flagEndHeight   int64
		flagStartHeight int64
		flagLogLevel    string
	)

	pflag.Int64VarP(&flagStartHeight, "start", "s", 0, "height at which to start indexing")
	pflag.Int64VarP(&flagEndHeight, "end", "e", 0, "height at which to stop indexing")
	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "log level")

	pflag.Parse()

	if len(os.Args) < 4 {
		return fmt.Errorf("required argument: <node_url> <network> <chain>")
	}
	nodeURL := os.Args[1]
	network := os.Args[2]
	chain := os.Args[3]

	// Logger initialization.
	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	log := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel)
	level, err := zerolog.ParseLevel(flagLogLevel)
	if err != nil {
		return err
	}
	log = log.Level(level)

	client, err := ethclient.Dial(nodeURL)
	if err != nil {
		return err
	}

	// TODO: Currently, we omit the case where start height is 0 and end height is non-zero,
	//       since this use-case (indexing part of the historical data from the beginning)
	//       is not yet relevant. It can be handled later if it becomes so.
	//       See https://github.com/NFT-com/indexer/issues/3.
	var sources []source.Source
	if flagStartHeight != 0 {
		historical, err := ethereum.NewHistorical(ctx, log, client, flagStartHeight, flagEndHeight)
		if err != nil {
			return err
		}

		sources = append(sources, historical)
	}
	if flagEndHeight == 0 {
		live, err := ethereum.NewLive(ctx, log, client)
		if err != nil {
			return err
		}

		sources = append(sources, live)
	}

	parser := ethereum.NewParser(log, client, network, chain)

	subs, err := subscriber.NewSubscriber(log, parser, sources)
	if err != nil {
		return err
	}

	// FIXME: Update to handle not only local
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.AnonymousCredentials,
		Region:      aws.String("us-west-2"),
	}))

	lambdaClient := lambda.New(sess, &aws.Config{
		Endpoint: aws.String("http://127.0.0.1:3001"),
	})

	dispatcher := dispatcher.New(lambdaClient, mock.New())

	failed := make(chan error)
	done := make(chan struct{})
	eventChannel := make(chan *event.Event)
	go func() {
		log.Info().Msg("Launching subscriber")
		if err := subs.Subscribe(ctx, eventChannel); err != nil {
			failed <- err
		}
		log.Info().Msg("Stopped subscriber")
		close(done)
	}()

	go func() {
		for {
			select {
			case e := <-eventChannel:
				log.Info().Interface("event", e).Msg("received")
				if err := dispatcher.Dispatch(ctx, e); err != nil {
					//log.Error().Err(err).Str("event", e.ID).Msg("failed to dispatch event")
				}
			case <-done:
				return
			}
		}
	}()

	select {
	case <-done:
		return nil
	case <-sig:
		if err := subs.Close(); err != nil {
			failed <- err
		}
	case err := <-failed:
		return err
	}

	go func() {
		<-sig
		log.Fatal().Msg("forced interruption")
	}()

	<-done

	return nil
}
