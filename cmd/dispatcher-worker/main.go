package main

import (
	"errors"
	"fmt"
	"github.com/NFT-com/indexer/dispatch"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"os"
	"os/signal"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"
)

const (
	IndexBase = 10
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("failure: %v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func run() error {
	// Signal catching for clean shutdown.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	// Command line parameter initialization.
	var flagLogLevel string

	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "log level")
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

	// TODO: Check if this will use the dispatch url or will use this from the parameter
	client, err := ethclient.Dial(nodeURL)
	if err != nil {
		return err
	}
	defer client.Close()

	var input dispatch.Dispatch // TODO RECEIVE THIS FROM QUEUE

	var (
		startIndex *big.Int
		endIndex   *big.Int
		contracts  []common.Address
	)

	if _, ok := startIndex.SetString(input.StartIndex, IndexBase); !ok {
		return errors.New("failed to parse start index")
	}

	if _, ok := endIndex.SetString(input.EndIndex, IndexBase); !ok {
		return errors.New("failed to parse end index")
	}

	for _, contract := range input.Contracts {
		contracts = append(contracts, common.HexToAddress(contract))
	}

	return nil
}
