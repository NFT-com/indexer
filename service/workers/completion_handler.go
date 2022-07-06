package workers

import (
	"context"
	"fmt"
	"strings"

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

type CompletionHandler struct {
	log zerolog.Logger
	url string
}

func NewCompletionHandler(log zerolog.Logger, url string) *CompletionHandler {

	e := CompletionHandler{
		log: log,
		url: url,
	}

	return &e
}

func (p *CompletionHandler) Handle(ctx context.Context, completion *jobs.Completion) (*results.Completion, error) {

	log := p.log.With().
		Str("job_id", completion.ID).
		Uint64("chain_id", completion.ChainID).
		Uint64("block_number", completion.BlockNumber).
		Str("transaction_hash", completion.TransactionHash).
		Strs("event_hashes", completion.EventHashes).
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

	// Retrieve the logs for all the addresses and event types for the given block range.
	requests := uint(1)
	logs, err := fetch.Logs(ctx, nil, completion.EventHashes, completion.BlockNumber, completion.BlockNumber)
	if err != nil {
		return nil, fmt.Errorf("could not fetch logs: %w", err)
	}

	p.log.Debug().
		Int("logs", len(logs)).
		Msg("event logs fetched")

	sale := events.Sale{
		ID:                completion.SaleID,
		ChainID:           completion.ChainID,
		CollectionAddress: "",
		TokenID:           "",
		BlockNumber:       completion.BlockNumber,
		TransactionHash:   completion.TransactionHash,
		SellerAddress:     completion.Seller,
		BuyerAddress:      completion.Buyer,
	}

	// For each log, try to parse it into the respective events.
	var transferCount int
	for _, log := range logs {

		// skip logs for reverted transactions
		if log.Removed {
			p.log.Trace().
				Uint("index", log.Index).
				Msg("skipping log for reverted transaction")
			continue
		}

		// skip logs that are not related to the transaction hash we are looking for
		if strings.ToLower(log.TxHash.Hex()) != strings.ToLower(completion.TransactionHash) {
			p.log.Trace().
				Uint("index", log.Index).
				Msg("skipping log unrelated to transaction of interest")
			continue
		}

		eventType := log.Topics[0]
		switch eventType.String() {

		case params.HashERC721Transfer:

			if len(log.Topics) != 4 {
				p.log.Warn().
					Uint("index", log.Index).
					Int("topics", len(log.Topics)).
					Msg("skipping log invalid topic length")
				continue
			}

			transfer, err := parsers.ERC721Transfer(log)
			if err != nil {
				return nil, fmt.Errorf("could not parse sale ERC721 transfer: %w", err)
			}

			if transfer.SenderAddress != completion.Seller ||
				transfer.ReceiverAddress != completion.Buyer {
				p.log.Trace().
					Uint("index", log.Index).
					Int("topics", len(log.Topics)).
					Msg("skipping log unrelated seller and buyer")
				continue
			}

			sale.CollectionAddress = transfer.CollectionAddress
			sale.TokenID = transfer.TokenID

			transferCount++

			p.log.Trace().
				Uint("index", log.Index).
				Str("collection_address", transfer.CollectionAddress).
				Str("token_id", transfer.TokenID).
				Msg("sale ERC721 transfer parsed")

		case params.HashERC1155Transfer:

			transfer, err := parsers.ERC1155Transfer(log)
			if err != nil {
				return nil, fmt.Errorf("could not parse ERC1155 transfer: %w", err)
			}

			if transfer.SenderAddress != completion.Seller ||
				transfer.ReceiverAddress != completion.Buyer {
				p.log.Trace().
					Uint("index", log.Index).
					Int("topics", len(log.Topics)).
					Msg("skipping log unrelated seller and buyer")
				continue
			}

			sale.CollectionAddress = transfer.CollectionAddress
			sale.TokenID = transfer.TokenID

			transferCount++

			p.log.Trace().
				Uint("index", log.Index).
				Str("collection_address", transfer.CollectionAddress).
				Str("token_id", transfer.TokenID).
				Msg("sale ERC1155 transfer parsed")
		}
	}

	if transferCount == 0 {
		return nil, fmt.Errorf("no transfer event found for sale")
	}

	if transferCount > 1 {
		return nil, fmt.Errorf("multiple transfer events found for sale")
	}

	// Put everything together for the result.
	result := results.Completion{
		Job:      completion,
		Sale:     &sale,
		Requests: requests,
	}

	return &result, nil
}
