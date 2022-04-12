package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/pflag"

	"github.com/NFT-com/indexer/networks/web3"
	"github.com/NFT-com/indexer/service/client"
	"github.com/NFT-com/indexer/watcher/chain"
)

const (
	defaultHTTPTimeout = time.Second * 30

	// This defaults the batch of historic data to 200 messages per request each second
	defaultBatchDelay = time.Second
	defaultBatch      = 200
)

func main() {
	err := run()
	if err != nil {
		// TODO: Improve this mixing logging
		// https://github.com/NFT-com/indexer/issues/32
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()

	// Signal catching for clean shutdown.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	// Command line parameter initialization.
	var (
		flagAPIEndpoint  string
		flagBatch        int64
		flagBatchDelay   time.Duration
		flagChainID      string
		flagChainURL     string
		flagChainType    string
		flagContract     string
		flagEventType    string
		flagLogLevel     string
		flagStandardType string
	)

	pflag.StringVarP(&flagAPIEndpoint, "api", "a", "", "jobs api base endpoint")
	pflag.Int64VarP(&flagBatch, "batch", "b", defaultBatch, "number of jobs to publish each batch")
	pflag.DurationVar(&flagBatchDelay, "batch-delay", defaultBatchDelay, "delay between each batch request")
	pflag.StringVarP(&flagChainID, "chain-id", "i", "", "id of the chain")
	pflag.StringVarP(&flagChainURL, "chain-url", "u", "", "url of the chain to connect")
	pflag.StringVarP(&flagChainType, "chain-type", "t", "", "type of chain")
	pflag.StringVarP(&flagContract, "contract", "c", "", "contract to watch")
	pflag.StringVarP(&flagEventType, "event", "e", "", "event type to watch")
	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "log level")
	pflag.StringVar(&flagStandardType, "standard-type", "", "standard type")
	pflag.Parse()

	// Logger initialization.
	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	log := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel)
	level, err := zerolog.ParseLevel(flagLogLevel)
	if err != nil {
		return fmt.Errorf("could not parse log level: %w", err)
	}
	log = log.Level(level)

	failed := make(chan error)

	network, err := web3.New(ctx, flagChainURL)
	if err != nil {
		return fmt.Errorf("could not create web3 network: %w", err)
	}

	chainID, err := network.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("could not get chain id from network: %w", err)
	}

	if chainID != flagChainID {
		return fmt.Errorf("could not start watcher: mismatch between chain ID and chain URL")
	}

	cli := http.DefaultClient
	cli.Timeout = defaultHTTPTimeout

	api := client.New(log,
		client.WithHTTPClient(cli),
		client.WithHost(flagAPIEndpoint),
	)

	highestJobIndex := "1"
	highestJob, err := api.GetHighestBlockNumberParsingJob(flagChainURL, flagChainType, flagContract, flagStandardType, flagEventType)
	if err == nil {
		highestJobIndex = highestJob.BlockNumber
	}

	cfg := chain.Config{
		ChainURL:     flagChainURL,
		ChainID:      chainID,
		ChainType:    flagChainType,
		StandardType: flagStandardType,
		Contract:     flagContract,
		EventType:    flagEventType,
		StartIndex:   highestJobIndex,
		Batch:        flagBatch,
		BatchDelay:   flagBatchDelay,
	}

	watcher, err := chain.NewWatcher(log, ctx, api, network, cfg)
	if err != nil {
		return fmt.Errorf("could not create watcher: %w", err)
	}

	go func() {
		log.Info().Msg("chain watcher starting")

		err = watcher.Watch(ctx)
		if err != nil {
			failed <- fmt.Errorf("could not watch chain: %w", err)
		}

		log.Info().Msg("chain watcher done")
	}()

	select {
	case <-sig:
		log.Info().Msg("chain watcher stopping")
		network.Close()
		watcher.Close()
		api.Close()
	case err = <-failed:
		log.Error().Err(err).Msg("chain watcher aborted")
		return err
	}

	go func() {
		<-sig
		log.Warn().Msg("forcing exit")
		os.Exit(1)
	}()

	return nil
}
