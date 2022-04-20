package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/pflag"

	models "github.com/NFT-com/indexer/models/chain"
	"github.com/NFT-com/indexer/networks/web3"
	"github.com/NFT-com/indexer/service/client"
	"github.com/NFT-com/indexer/service/postgres"
	"github.com/NFT-com/indexer/watcher/chain"
)

const (
	databaseDriver = "postgres"

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
		flagAPIEndpoint      string
		flagBatch            int64
		flagBatchDelay       time.Duration
		flagChainID          string
		flagChainURL         string
		flagChainType        string
		flagDBConnectionInfo string
		flagLogLevel         string
	)

	pflag.StringVarP(&flagAPIEndpoint, "api", "a", "", "jobs api base endpoint")
	pflag.Int64VarP(&flagBatch, "batch", "b", defaultBatch, "number of jobs to publish each batch")
	pflag.DurationVar(&flagBatchDelay, "batch-delay", defaultBatchDelay, "delay between each batch request")
	pflag.StringVarP(&flagChainID, "chain-id", "i", "", "id of the chain")
	pflag.StringVarP(&flagChainURL, "chain-url", "u", "", "url of the chain to connect")
	pflag.StringVarP(&flagChainType, "chain-type", "t", "", "type of chain")
	pflag.StringVarP(&flagDBConnectionInfo, "db", "d", "", "database connection string")
	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "log level")
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

	db, err := sql.Open(databaseDriver, flagDBConnectionInfo)
	if err != nil {
		return fmt.Errorf("could not open SQL connection: %w", err)
	}

	store, err := postgres.NewStore(db)
	if err != nil {
		return fmt.Errorf("could not create store: %w", err)
	}

	storedChain, err := store.Chain(chainID)
	if err != nil {
		return fmt.Errorf("could not get chain from database: %w", err)
	}

	collections, contracts, err := getCollections(store, storedChain.ID)
	if err != nil {
		return fmt.Errorf("could not get collections: %w", err)
	}

	standards, contractStandards, err := getStandards(store, collections)
	if err != nil {
		return fmt.Errorf("could not get standars: %w", err)
	}

	standardsEventTypes, err := getEventTypes(store, standards)
	if err != nil {
		return fmt.Errorf("could not get event types: %w", err)
	}

	highestJobIndexes, startingBlock := getHighestJobBlockNumberForCollections(api, flagChainURL, flagChainType, contracts, contractStandards, standardsEventTypes)
	fmt.Println(highestJobIndexes, startingBlock)
	cfg := chain.Config{
		ChainURL:      flagChainURL,
		ChainID:       chainID,
		ChainType:     flagChainType,
		Contracts:     contracts,
		Standards:     contractStandards,
		EventTypes:    standardsEventTypes,
		StartingBlock: startingBlock,
		StartIndexes:  highestJobIndexes,
		Batch:         flagBatch,
		BatchDelay:    flagBatchDelay,
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

func getCollections(store *postgres.Store, chainID string) ([]models.Collection, []string, error) {
	collections, err := store.Collections(chainID)
	if err != nil {
		return nil, nil, fmt.Errorf("could not get collections from database: %w", err)
	}

	contracts := make([]string, 0, len(collections))
	for _, collection := range collections {
		contracts = append(contracts, collection.Address)
	}

	return collections, contracts, nil
}

func getStandards(store *postgres.Store, collections []models.Collection) ([]models.Standard, map[string][]string, error) {
	standards := make([]models.Standard, 0, len(collections))
	contractStandards := make(map[string][]string, len(collections))
	for _, collection := range collections {
		var err error
		storedStandards, err := store.Standards(collection.ID)
		if err != nil {
			return nil, nil, fmt.Errorf("could not get standards from database: %w", err)
		}

		names := make([]string, 0, len(storedStandards))
		for _, standard := range storedStandards {
			names = append(names, standard.Name)
		}

		standards = append(standards, storedStandards...)
		contractStandards[collection.Address] = names
	}

	return standards, contractStandards, nil
}

func getEventTypes(store *postgres.Store, standards []models.Standard) (map[string][]string, error) {
	eventTypes := make(map[string][]string, len(standards))

	for _, standard := range standards {
		types, err := store.EventTypes(standard.ID)
		if err != nil {
			return nil, fmt.Errorf("could not get event types from database: %w", err)
		}

		ids := make([]string, 0, len(types))
		for _, eventType := range types {
			ids = append(ids, eventType.ID)
		}

		eventTypes[standard.Name] = ids
	}

	return eventTypes, nil
}

func getHighestJobBlockNumberForCollections(api *client.Client, chainURL, chainType string, contracts []string, standards map[string][]string, eventTypes map[string][]string) (chain.Indexes, *big.Int) {
	startingBlock := big.NewInt(-1)

	highestJobIndexes := make(chain.Indexes, len(contracts))
	for _, contract := range contracts {
		stands := standards[contract]

		collectionIndexes, lowestBlock := getHighestJobNumberForStandards(api, chainURL, chainType, contract, stands, eventTypes)
		if len(collectionIndexes) == 0 {
			continue
		}

		if startingBlock.CmpAbs(big.NewInt(-1)) == 0 ||
			startingBlock.CmpAbs(lowestBlock) > 0 {

			startingBlock.SetBytes(lowestBlock.Bytes())
		}

		highestJobIndexes[contract] = collectionIndexes
	}

	return highestJobIndexes, startingBlock
}

func getHighestJobNumberForStandards(api *client.Client, chainURL, chainType, contract string, standards []string, eventTypes map[string][]string) (chain.CollectionIndexes, *big.Int) {
	startingBlock := big.NewInt(-1)

	collectionIndexes := make(chain.CollectionIndexes, len(standards))
	for _, standard := range standards {
		eTypes := eventTypes[standard]

		eventTypesIndexes, lowestBlock := getHighestJobNumberForEventTypes(api, chainURL, chainType, contract, standard, eTypes)
		if len(eventTypesIndexes) == 0 {
			continue
		}

		if startingBlock.CmpAbs(big.NewInt(-1)) == 0 ||
			startingBlock.CmpAbs(lowestBlock) > 0 {

			startingBlock.SetBytes(lowestBlock.Bytes())
		}

		collectionIndexes[standard] = eventTypesIndexes
	}

	return collectionIndexes, startingBlock
}

func getHighestJobNumberForEventTypes(api *client.Client, chainURL, chainType, contract, standard string, eTypes []string) (chain.EventTypesIndexes, *big.Int) {
	startingBlock := big.NewInt(-1)

	eventTypesIndexes := make(chain.EventTypesIndexes, len(eTypes))
	for _, eType := range eTypes {
		highestJobIndex := big.NewInt(1)
		highestJob, err := api.GetHighestBlockNumberParsingJob(chainURL, chainType, contract, standard, eType)
		if err == nil {
			highestJobIndex.SetString(highestJob.BlockNumber, 0)
		}

		if startingBlock.CmpAbs(big.NewInt(-1)) == 0 ||
			startingBlock.CmpAbs(highestJobIndex) > 0 {

			startingBlock.SetBytes(highestJobIndex.Bytes())
		}

		eventTypesIndexes[eType] = highestJobIndex
	}

	return eventTypesIndexes, startingBlock
}
