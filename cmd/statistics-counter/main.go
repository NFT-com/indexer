package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/NFT-com/indexer/network/ethereum"
	"github.com/NFT-com/indexer/network/web3"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"
	"go.uber.org/ratelimit"
)

const (
	success = 0
	failure = 1
)

func main() {
	os.Exit(run())
}

func run() int {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	var (
		flagLogLevel string

		flagNodeURL string

		flagStartingHeight uint64
		flagAddresses      []string
		flagEventHashes    []string

		flagRateLimit int
	)

	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "severity level for log output")

	pflag.StringVarP(&flagNodeURL, "node", "n", "", "ethereum node url")

	pflag.Uint64VarP(&flagStartingHeight, "starting-height", "s", 0, "counter starting block")
	pflag.StringSliceVarP(&flagAddresses, "addresses", "a", []string{}, "addresses to count")
	pflag.StringSliceVarP(&flagEventHashes, "hashes", "h", []string{}, "hashes to count")

	pflag.IntVarP(&flagRateLimit, "rate-limit", "r", 100, "node requests per second")

	pflag.Parse()

	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	log := zerolog.New(os.Stderr).With().Timestamp().Logger()
	level, err := zerolog.ParseLevel(flagLogLevel)
	if err != nil {
		log.Error().Err(err).Str("log_level", flagLogLevel).Msg("could not parse log level")
		return failure
	}
	log = log.Level(level)

	if flagStartingHeight == 0 {
		log.Error().Msg("starting height must be defined")
		return failure
	}

	var api *ethclient.Client
	close := func() {}
	if strings.Contains(flagNodeURL, "ethereum.managedblockchain") {
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			log.Error().Err(err).Msg("could not load AWS configuration")
			return failure
		}
		api, close, err = ethereum.NewSigningClient(ctx, flagNodeURL, cfg)
		if err != nil {
			log.Error().Err(err).Str("url", flagNodeURL).Msg("could not create signing client")
			return failure
		}
	} else {
		api, err = ethclient.DialContext(ctx, flagNodeURL)
		if err != nil {
			log.Error().Err(err).Str("url", flagNodeURL).Msg("could not create default client")
			return failure
		}
	}
	defer api.Close()
	defer close()

	header, err := api.HeaderByNumber(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg("could not get latest header")
		return failure
	}

	fetcher := web3.NewLogsFetcher(api)
	counter := NewCounter()

	limiter := ratelimit.New(flagRateLimit)

	go func() {
		for height := flagStartingHeight; height <= header.Number.Uint64(); height += 10 {
			limiter.Take()
			if height%1000 == 0 {
				log.Info().Uint64("height", height).Msg("height processed")
			}

			logs, err := fetcher.Logs(ctx, flagAddresses, flagEventHashes, height, height+10)
			if err != nil {
				log.Error().Err(err).Msg("could not fetch logs")
				continue
			}

			for _, log := range logs {
				counter.Count(log)
			}
		}
		sig <- os.Interrupt
	}()

	log.Info().Uint64("start", flagStartingHeight).Uint64("end", header.Number.Uint64()).Msg("counter started")

	<-sig

	fmt.Println(counter.String())

	log.Info().Msg("initialized shutdown")

	go func() {
		<-sig
		log.Fatal().Msg("forced shutdown")
	}()

	log.Info().Msg("shutdown complete")

	return success
}

type Counter struct {
	counts map[string]uint64
}

func NewCounter() *Counter {
	c := Counter{
		counts: make(map[string]uint64),
	}

	return &c
}

func (c *Counter) Count(log types.Log) {
	address := log.Address.String()
	hash := log.Topics[0].String()

	key := strings.Join([]string{address, hash}, ",")

	currentValue, _ := c.counts[key]
	c.counts[key] = currentValue + 1
}

func (c *Counter) String() string {
	if len(c.counts) == 0 {
		return "No logs Found!"
	}

	output := fmt.Sprintf("%s | %s | %s\n", "Address", "Event Hash", "Count")
	for key, count := range c.counts {
		keys := strings.Split(key, ",")
		output += fmt.Sprintf("%s | %s | %d\n", keys[0], keys[1], count)
	}
	return output
}
