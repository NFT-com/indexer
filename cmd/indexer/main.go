package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/NFT-com/indexer/parse"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"

	"github.com/NFT-com/indexer/contracts"
	"github.com/NFT-com/indexer/events"
	"github.com/NFT-com/indexer/parse/cryptokitties"
	"github.com/NFT-com/indexer/parse/opensea"
	"github.com/NFT-com/indexer/store"
	"github.com/NFT-com/indexer/subscriber"
	"github.com/NFT-com/indexer/subscriber/ethereum"
)

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
	pflag.StringVarP(&flagLogLevel, "log-level", "-l", "info", "log level")

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

	// FIXME: Handle hybrid subscribing (historical + live) instead of one at a time.
	var subs subscriber.Subscriber
	switch {
	case flagStartHeight == 0 && flagEndHeight == 0:
		liveSub := ethereum.NewLive(log, client, manager)

		subs = liveSub

	default:
		historicalSub := ethereum.NewHistorical(log, client, manager, flagStartHeight, flagEndHeight)

		subs = historicalSub
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
