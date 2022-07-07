package workers

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/rs/zerolog"

	"github.com/ethereum/go-ethereum/core/types"
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
		Uint64("start_height", completion.StartHeight).
		Uint64("end_height", completion.EndHeight).
		Strs("event_hashes", completion.EventHashes).
		Strs("sale_ids", completion.SaleIDs()).
		Strs("transaction_hashes", completion.TransactionHashes()).
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
	logs, err := fetch.Logs(ctx, nil, completion.EventHashes, completion.StartHeight, completion.EndHeight)
	if err != nil {
		return nil, fmt.Errorf("could not fetch logs: %w", err)
	}

	p.log.Debug().
		Int("logs", len(logs)).
		Msg("event logs fetched")

	// Create a transaction hash look-up.
	hashLookup := make(map[string]struct{})
	for _, sale := range completion.Sales {
		hashLookup[sale.TransactionHash] = struct{}{}
	}

	// First, we create a lookup for logs by transaction hash.
	logLookup := make(map[string][]types.Log)
	for _, log := range logs {

		// Skip logs for transactions we don't care about.
		txHash := log.TxHash.Hex()
		_, ok := hashLookup[txHash]
		if !ok {
			continue
		}

		logLookup[txHash] = append(logLookup[txHash], log)
	}

	// Then we go through each sale...
	for _, sale := range completion.Sales {

		// ... and get the logs for each sale according to its transaction hash.
		logs, ok := logLookup[sale.TransactionHash]
		if !ok {
			p.log.Warn().Msg("no logs for transaction found")
			continue
		}

		// Next, we go through each log for that transaction hash and identify
		// all transfers that are potential matches.
		var matches []*events.Transfer
		for _, log := range logs {

			// skip logs for reverted transactions
			if log.Removed {
				p.log.Trace().
					Uint("index", log.Index).
					Msg("skipping log for reverted transaction")
				continue
			}

			eventHash := log.Topics[0].String()
			switch eventHash {

			case params.HashERC721Transfer:

				if len(log.Topics) != 4 {
					p.log.Warn().
						Uint("index", log.Index).
						Int("topics", len(log.Topics)).
						Msg("skipping log with invalid topic length")
					continue
				}

				transfer, err := parsers.ERC721Transfer(log)
				if err != nil {
					return nil, fmt.Errorf("could not parse sale ERC721 transfer: %w", err)
				}

				matches = append(matches, transfer)

			case params.HashERC1155Transfer:

				transfer, err := parsers.ERC1155Transfer(log)
				if err != nil {
					return nil, fmt.Errorf("could not parse ERC1155 transfer: %w", err)
				}

				matches = append(matches, transfer)

			default:

				p.log.Debug().Str("event_hash", eventHash).Msg("skipping unknown event hash")
			}
		}

		// Finally, we assign the data to the sale if we have exactly one match.
		if len(matches) > 1 {
			p.log.Warn().Msg("found multiple matching transfers for sale, skipping")
			continue
		}

		if len(matches) == 0 {
			p.log.Warn().Msg("found no matching transfers for sale, skipping")
		}

		sale.CollectionAddress = matches[0].CollectionAddress
		sale.TokenID = matches[0].TokenID
	}

	// Put everything together for the result.
	result := results.Completion{
		Job:      completion,
		Requests: requests,
	}

	return &result, nil
}
