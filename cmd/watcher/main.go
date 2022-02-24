package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/adjust/rmq/v4"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"

	"github.com/NFT-com/indexer/dispatch"
	"github.com/NFT-com/indexer/dispatch/redismq"
	"github.com/NFT-com/indexer/networks/ethereum"
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
		flagRMQTag        string
		flagRedisNetwork  string
		flagRedisURL      string
		flagRedisDatabase int
		flagQueueName     string
		flagLogLevel      string
	)

	pflag.StringVarP(&flagRMQTag, "tag", "t", "watcher", "watcher redismq tag")
	pflag.StringVarP(&flagRedisNetwork, "network", "n", "tcp", "network")
	pflag.StringVarP(&flagRedisURL, "url", "u", "", "redis url")
	pflag.IntVarP(&flagRedisDatabase, "database", "d", 1, "redis database")
	pflag.StringVarP(&flagQueueName, "queue", "q", "discovery", "queue name")
	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "log level")
	pflag.Parse()

	if len(os.Args) < 3 {
		return fmt.Errorf("required arguments: <node_url> <chain_type>")
	}
	nodeURL := os.Args[1]
	chainType := os.Args[2]

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

	live, err := ethereum.NewLive(ctx, log, client)
	if err != nil {
		return err
	}

	failed := make(chan error)
	done := make(chan struct{})

	connection, err := rmq.OpenConnection(flagRMQTag, flagRedisNetwork, flagRedisURL, flagRedisDatabase, failed)
	if err != nil {
		return err
	}

	producer, err := redismq.NewProducer(connection)
	if err != nil {
		return err
	}

	go func() {
		log.Info().Msg("watcher started")
		for {
			block := live.Next(ctx)
			if block == nil {
				break
			}

			job := dispatch.DiscoveryJob{
				ChainURL:   nodeURL,
				ChainType:  chainType,
				StartIndex: block.Number,
				EndIndex:   block.Number,
			}

			err := producer.PublishDiscoveryJob(flagQueueName, job)
			if err != nil {
				failed <- err
				return
			}
		}
		log.Info().Msg("watcher stopped")
		close(done)
	}()

	select {
	case <-done:
		client.Close()
		return nil
	case <-sig:
		if err := live.Close(); err != nil {
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
