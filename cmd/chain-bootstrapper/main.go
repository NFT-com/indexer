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

	"github.com/NFT-com/indexer/bootstrapper"
	"github.com/NFT-com/indexer/networks/web3"
	"github.com/NFT-com/indexer/service/client"
)

const (
	defaultHTTPTimeout = time.Second * 30
)

func main() {
	err := run()
	if err != nil {
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
		flagChainID      string
		flagChainURL     string
		flagChainType    string
		flagContract     string
		flagEventType    string
		flagEndIndex     int64
		flagLogLevel     string
		flagStandardType string
		flagStartIndex   int64
	)

	pflag.StringVarP(&flagAPIEndpoint, "api", "a", "", "jobs api base endpoint")
	pflag.StringVarP(&flagChainID, "chain-id", "i", "", "id of the chain")
	pflag.StringVarP(&flagChainURL, "chain-url", "u", "", "url of the chain to connect")
	pflag.StringVarP(&flagChainType, "chain-type", "t", "", "type of chain")
	pflag.StringVarP(&flagContract, "contract", "c", "", "contract to watch")
	pflag.StringVar(&flagEventType, "event", "", "event type to watch")
	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "log level")
	pflag.StringVar(&flagStandardType, "standard-type", "", "standard type")
	pflag.Int64VarP(&flagStartIndex, "start-index", "s", 0, "start index")
	pflag.Int64VarP(&flagEndIndex, "end-index", "e", 0, "end index")
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

	cli := http.DefaultClient
	cli.Timeout = defaultHTTPTimeout

	api := client.New(log,
		client.WithHTTPClient(cli),
		client.WithHost(flagAPIEndpoint),
	)

	network, err := web3.New(ctx, flagChainURL)
	if err != nil {
		return fmt.Errorf("could not create web3 network: %w", err)
	}

	chainID, err := network.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("could not get chain id from network: %w", err)
	}

	if chainID != flagChainID {
		return fmt.Errorf("could not start bootstrapper: chain-id doesn't match id retrieved from chain-url")
	}

	bootstrapper := bootstrapper.New(
		log,
		api,
		flagChainURL,
		flagChainType,
		flagStandardType,
		flagContract,
		flagEventType,
		flagStartIndex,
		flagEndIndex,
	)
	if err != nil {
		return fmt.Errorf("could not create bootstrapper: %w", err)
	}

	done := make(chan struct{})
	go func() {
		log.Info().Msg("bootstrapper starting")

		err = bootstrapper.Bootstrap(ctx)
		if err != nil {
			failed <- fmt.Errorf("could not bootstrap: %w", err)
		}

		close(done)
		log.Info().Msg("bootstrapper done")
	}()

	select {
	case <-sig:
		log.Info().Msg("bootstrapper stopping")
	case err = <-failed:
		log.Error().Err(err).Msg("bootstrapper aborted")
		return err
	}

	select {
	case <-done:
		log.Info().Msg("chain bootstrapper done")
	case <-sig:
		log.Info().Msg("chain bootstrapper stopping")
	case err = <-failed:
		log.Error().Err(err).Msg("chain bootstrapper aborted")
		return err
	}

	bootstrapper.Close()
	network.Close()
	api.Close()

	go func() {
		<-sig
		log.Warn().Msg("forcing exit")
		os.Exit(1)
	}()

	return nil
}
