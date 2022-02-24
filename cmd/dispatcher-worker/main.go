package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"

	"github.com/NFT-com/indexer/block"
	"github.com/NFT-com/indexer/networks/ethereum"
	"github.com/NFT-com/indexer/source"
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
		flagStartHeight int64
		flagEndHeight   int64
		flagLogLevel    string
		flagTestMode    bool
		flagLambdaURL   string
		flagRegion      string
	)

	pflag.Int64VarP(&flagStartHeight, "start", "s", 0, "height at which to start indexing")
	pflag.Int64VarP(&flagEndHeight, "end", "e", 0, "height at which to stop indexing")
	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "log level")
	pflag.BoolVarP(&flagTestMode, "test", "t", false, "test mode")
	pflag.StringVarP(&flagLambdaURL, "lambda-url", "u", "", "lambda url")
	pflag.StringVarP(&flagRegion, "aws-region", "r", "eu-west-1", "aws region")

	pflag.Parse()

	if len(os.Args) < 2 {
		return fmt.Errorf("required arguments: <node_url>")
	}
	nodeURL := os.Args[1]

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

	subs, err := subscriber.NewSubscriber(log, sources)
	if err != nil {
		return err
	}

	failed := make(chan error)
	done := make(chan struct{})
	blockChannel := make(chan *block.Block, 1000)

	go func() {
		log.Info().Msg("Launching subscriber")
		if err := subs.Subscribe(ctx, blockChannel); err != nil {
			failed <- err
		}
		log.Info().Msg("Stopped subscriber")
		close(done)
	}()

	go func() {
		for {
			select {
			case b := <-blockChannel:
				log.Info().Interface("block", b).Msg("received")
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