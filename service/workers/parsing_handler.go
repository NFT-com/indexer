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
	logs, err := fetch.Logs(ctx, parsing.ContractAddresses, parsing.EventHashes, parsing.StartHeight, parsing.EndHeight)
	if err != nil {
		return nil, fmt.Errorf("could not fetch logs: %w", err)
	}

	p.log.Debug().
		Int("logs", len(logs)).
		Msg("event logs fetched")

	// For each log, try to parse it into the respective events.
	var transfers []*events.Transfer
	var sales []*events.Sale
	timestamps := make(map[uint64]time.Time)
	for _, log := range logs {

		// skip logs for reverted transactions
		if log.Removed {
			p.log.Trace().
				Str("transaction", log.TxHash.Hex()).
				Uint("index", log.Index).
				Msg("skipping log for reverted transaction")
			continue
		}

		// keep track of all heightSet we need to process to get timestamps
		timestamps[log.BlockNumber] = time.Time{}

		eventType := log.Topics[0]
		switch eventType.String() {

		case params.HashERC721Transfer:

			transfer, err := parsers.ERC721Transfer(log)
			if err != nil {
				return nil, fmt.Errorf("could not parse ERC721 transfer: %w", err)
			}
			transfers = append(transfers, transfer)

			p.log.Trace().
				Str("transaction", log.TxHash.Hex()).
				Uint("index", log.Index).
				Str("collection_address", transfer.CollectionAddress).
				Str("token_id", transfer.TokenID).
				Str("sender_address", transfer.SenderAddress).
				Str("receiver_address", transfer.ReceiverAddress).
				Str("token_count", transfer.TokenCount).
				Msg("ERC721 transfer parsed")

		case params.HashERC1155Transfer:

			transfer, err := parsers.ERC1155Transfer(log)
			if err != nil {
				return nil, fmt.Errorf("could not parse ERC1155 transfer: %w", err)
			}
			transfers = append(transfers, transfer)

			p.log.Trace().
				Str("transaction", log.TxHash.Hex()).
				Uint("index", log.Index).
				Str("collection_address", transfer.CollectionAddress).
				Str("token_id", transfer.TokenID).
				Str("sender_address", transfer.SenderAddress).
				Str("receiver_address", transfer.ReceiverAddress).
				Str("token_count", transfer.TokenCount).
				Msg("ERC1155 transfer parsed")

		case params.HashERC1155Batch:

			batch, err := parsers.ERC1155Batch(log)
			if err != nil {
				return nil, fmt.Errorf("could not parse ERC1155 batch: %w", err)
			}
			transfers = append(transfers, batch...)
			for _, transfer := range batch {

				p.log.Trace().
					Str("transaction", log.TxHash.Hex()).
					Uint("index", log.Index).
					Str("collection_address", transfer.CollectionAddress).
					Str("token_id", transfer.TokenID).
					Str("sender_address", transfer.SenderAddress).
					Str("receiver_address", transfer.ReceiverAddress).
					Str("token_count", transfer.TokenCount).
					Msg("ERC115 batch parsed")
			}

		case params.HashOpenSeaWyvernTrade:

			sale, err := parsers.OpenSeaWyvernSale(log)
			if err != nil {
				return nil, fmt.Errorf("could not parse sale: %w", err)
			}
			sales = append(sales, sale)

			p.log.Trace().
				Str("transaction", log.TxHash.Hex()).
				Uint("index", log.Index).
				Msg("OpenSea Wyvern sale parsed")

		case params.HashOpenSeaSeaportTrade:

			sale, err := parsers.OpenSeaSeaportSale(log)
			if err != nil {
				return nil, fmt.Errorf("could not parse sale: %w", err)
			}
			sales = append(sales, sale)

			p.log.Trace().
				Str("transaction", log.TxHash.Hex()).
				Uint("index", log.Index).
				Msg("OpenSea Seaport sale parsed")
		}
	}

	p.log.Info().
		Int("logs", len(logs)).
		Int("transfers", len(transfers)).
		Int("sales", len(sales)).
		Msg("all logs parsed")

	// Get all the headers to assign timestamps to the events.
	for height := range timestamps {
		requests++
		header, err := api.HeaderByNumber(ctx, big.NewInt(0).SetUint64(height))
		if err != nil {
			return nil, fmt.Errorf("could not get header: %w", err)
		}
		timestamps[height] = time.Unix(int64(header.Time), 0)
	}

	p.log.Info().
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
		Requests:  requests,
	}

	return &result, nil
}
