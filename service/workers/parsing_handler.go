package workers

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/rs/zerolog"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/NFT-com/indexer/config/params"
	"github.com/NFT-com/indexer/models/events"
	"github.com/NFT-com/indexer/models/jobs"
	"github.com/NFT-com/indexer/models/results"
	"github.com/NFT-com/indexer/network/ethereum"
	"github.com/NFT-com/indexer/network/web3"
	"github.com/NFT-com/indexer/service/parsers"
)

type ParsingHandler struct {
	log zerolog.Logger
	url string
}

func NewParsingHandler(log zerolog.Logger, url string) *ParsingHandler {

	e := ParsingHandler{
		log: log,
		url: url,
	}

	return &e
}

func (p *ParsingHandler) Handle(ctx context.Context, parsing *jobs.Parsing) (*results.Parsing, error) {

	log := p.log.With().
		Str("job_id", parsing.ID).
		Uint64("chain_id", parsing.ChainID).
		Uint64("start_height", parsing.StartHeight).
		Uint64("end_height", parsing.EndHeight).
		Strs("contract_addresses", parsing.ContractAddresses).
		Strs("event_hashes", parsing.EventHashes).
		Logger()

	log.Info().
		Str("node_url", p.url).
		Msg("initiating connection to Ethereum node")

	var err error
	var api *ethclient.Client
	close := func() {}
	if strings.Contains(p.url, "ethereum.managedblockchain") {
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not load AWS configuration: %w", err)
		}
		api, close, err = ethereum.NewSigningClient(ctx, p.url, cfg)
		if err != nil {
			return nil, fmt.Errorf("could not create signing client (url: %s): %w", p.url, err)
		}
	} else {
		api, err = ethclient.DialContext(ctx, p.url)
		if err != nil {
			return nil, fmt.Errorf("could not create default client (url: %s): %w", p.url, err)
		}
	}
	defer api.Close()
	defer close()

	log.Info().
		Str("node_url", p.url).
		Msg("connection to Ethereum node established")

	fetch := web3.NewLogsFetcher(api)

	// Retrieve the logs for all of the addresses and event types for the given block range.
	requests := uint(1)
	entries, err := fetch.Logs(ctx, parsing.ContractAddresses, parsing.EventHashes, parsing.StartHeight, parsing.EndHeight)
	if err != nil {
		return nil, fmt.Errorf("could not fetch logs: %w", err)
	}

	log.Debug().
		Int("entries", len(entries)).
		Msg("event log entries fetched")

	// For each log, try to parse it into the respective events.
	var transfers []*events.Transfer
	var sales []*events.Sale
	timestamps := make(map[uint64]time.Time)
	for _, entry := range entries {

		// skip logs for reverted transactions
		if entry.Removed {
			log.Trace().
				Str("transaction", entry.TxHash.Hex()).
				Uint("index", entry.Index).
				Msg("skipping log entry for reverted transaction")
			continue
		}

		eventType := entry.Topics[0]
		switch eventType.String() {

		case params.HashERC721Transfer:

			transfer, err := parsers.ERC721Transfer(entry)
			if err != nil {
				log.Warn().
					Err(err).
					Hex("transaction", entry.TxHash[:]).
					Uint("index", entry.Index).
					Msg("could not parse ERC721 transfer, skipping log entry")
				continue
			}

			timestamps[entry.BlockNumber] = time.Time{}
			transfers = append(transfers, transfer)

			log.Trace().
				Str("transaction", entry.TxHash.Hex()).
				Uint("index", entry.Index).
				Str("collection_address", transfer.CollectionAddress).
				Str("token_id", transfer.TokenID).
				Str("sender_address", transfer.SenderAddress).
				Str("receiver_address", transfer.ReceiverAddress).
				Str("token_count", transfer.TokenCount).
				Msg("ERC721 transfer parsed")

		case params.HashERC1155Transfer:

			transfer, err := parsers.ERC1155Transfer(entry)
			if err != nil {
				log.Warn().
					Err(err).
					Hex("transaction", entry.TxHash[:]).
					Uint("index", entry.Index).
					Msg("could not parse ERC1155 transfer, skipping log entry")
				continue
			}

			timestamps[entry.BlockNumber] = time.Time{}
			transfers = append(transfers, transfer)

			log.Trace().
				Str("transaction", entry.TxHash.Hex()).
				Uint("index", entry.Index).
				Str("collection_address", transfer.CollectionAddress).
				Str("token_id", transfer.TokenID).
				Str("sender_address", transfer.SenderAddress).
				Str("receiver_address", transfer.ReceiverAddress).
				Str("token_count", transfer.TokenCount).
				Msg("ERC1155 transfer parsed")

		case params.HashERC1155Batch:

			batch, err := parsers.ERC1155Batch(entry)
			if err != nil {
				log.Warn().
					Err(err).
					Hex("transaction", entry.TxHash[:]).
					Uint("index", entry.Index).
					Msg("could not parse ERC1155 batch, skipping log entry")
				continue
			}

			timestamps[entry.BlockNumber] = time.Time{}
			transfers = append(transfers, batch...)

			for _, transfer := range batch {

				log.Trace().
					Str("transaction", entry.TxHash.Hex()).
					Uint("index", entry.Index).
					Str("collection_address", transfer.CollectionAddress).
					Str("token_id", transfer.TokenID).
					Str("sender_address", transfer.SenderAddress).
					Str("receiver_address", transfer.ReceiverAddress).
					Str("token_count", transfer.TokenCount).
					Msg("ERC115 batch parsed")
			}

		case params.HashWyvernSale:

			sale, err := parsers.WyvernSale(entry)
			if err != nil {
				log.Warn().
					Err(err).
					Hex("transaction", entry.TxHash[:]).
					Uint("index", entry.Index).
					Msg("could not parse Wyvern sale, skipping log entry")
				continue
			}

			timestamps[entry.BlockNumber] = time.Time{}
			sales = append(sales, sale)

			log.Trace().
				Str("transaction", entry.TxHash.Hex()).
				Uint("index", entry.Index).
				Msg("OpenSea Wyvern sale parsed")

		case params.HashSeaportSale:

			sale, err := parsers.SeaportSale(entry)
			if err != nil {
				log.Warn().
					Err(err).
					Hex("transaction", entry.TxHash[:]).
					Uint("index", entry.Index).
					Msg("could not parse Seaport sale, skipping log entry")
				continue
			}

			timestamps[entry.BlockNumber] = time.Time{}
			sales = append(sales, sale)

			log.Trace().
				Str("transaction", entry.TxHash.Hex()).
				Uint("index", entry.Index).
				Msg("OpenSea Seaport sale parsed")
		}
	}

	log.Info().
		Int("entries", len(entries)).
		Int("transfers", len(transfers)).
		Int("sales", len(sales)).
		Msg("all log entries parsed")

	// Get all the headers to assign timestamps to the events.
	for height := range timestamps {
		requests++
		header, err := api.HeaderByNumber(ctx, big.NewInt(0).SetUint64(height))
		if err != nil {
			return nil, fmt.Errorf("could not get header: %w", err)
		}
		timestamps[height] = time.Unix(int64(header.Time), 0)
	}

	log.Info().
		Int("heights", len(timestamps)).
		Msg("block heights retrieved")

	// Go through all logs and assign timestamp of emission
	for _, transfer := range transfers {
		transfer.ChainID = parsing.ChainID
		transfer.EmittedAt = timestamps[transfer.BlockNumber]
	}
	for _, sale := range sales {
		sale.ChainID = parsing.ChainID
		sale.EmittedAt = timestamps[sale.BlockNumber]
	}

	// Put everything together for the result.
	result := results.Parsing{
		Job:       parsing,
		Sales:     sales,
		Transfers: transfers,
	}

	return &result, nil
}
