package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"

	"github.com/NFT-com/indexer/contracts"
	"github.com/NFT-com/indexer/events"
	"github.com/NFT-com/indexer/parse"
	"github.com/NFT-com/indexer/parse/cryptokitties"
	ethParse "github.com/NFT-com/indexer/parse/ethereum"
	"github.com/NFT-com/indexer/parse/opensea"
	"github.com/NFT-com/indexer/source"
	"github.com/NFT-com/indexer/source/ethereum"
	"github.com/NFT-com/indexer/store"
	"github.com/NFT-com/indexer/subscriber"
)

// FIXME: ethParse no longer needs an alias once we remove the test code which
//        currently imports the parse package.

func main() {
	if err := run(); err != nil {
		fmt.Printf("failure: %v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func run() error {
	ctx := context.Background()

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

	if len(os.Args) < 2 {
		return fmt.Errorf("required argument: <node_url>")
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

	contractStore, err := store.NewMockStore()
	if err != nil {
		return err
	}

	manager := contracts.New(log, client, contractStore)
	parser := ethParse.NewParser(log, client, manager)

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

	subs, err := subscriber.NewSubscriber(log, parser, sources)
	if err != nil {
		return err
	}

	failed := make(chan error)
	done := make(chan struct{})
	eventChannel := make(chan events.Event)
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
			event := <-eventChannel

			log.Info().Interface("event", event).Msg("received")

			// FIXME this code is just for testing purposes
			switch ev := event.(type) {
			case *events.Transfer:
				if ev.Address == common.HexToAddress("0x06012c8cf97bead5deae237070f9587f8e7a266d") {
					nft := parse.NFT{
						ID:      int64(ev.NftID),
						Address: ev.Address.Hex(),
						Chain:   ev.Chain(),
						Network: ev.Network(),
					}

					err = cryptokitties.Handler(client)(ctx, nft)
					if err != nil {
						log.Error().Err(err).Msg("could not handle CryptoKitties event")
					}
				}
			case *events.OrdersMatched:
				nft := parse.NFT{
					ID:      int64(ev.Price),
					Address: ev.Address.Hex(),
					Chain:   ev.Chain(),
					Network: ev.Network(),
				}

				err = opensea.Handler(client)(ctx, nft)
				if err != nil {
					log.Error().Err(err).Msg("could not handle OpenSea event")
				}
			}
		}
	}()

	select {
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
