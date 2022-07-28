package workers

import (
	"context"
	"fmt"
	"sort"
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

	// This orders for the case that we have multiple sales per transfer.
	saleLookup := make(map[string][]*events.Sale) // FIXME: var name
	for _, sale := range completion.Sales {
		saleLookup[sale.TransactionHash] = append(saleLookup[sale.TransactionHash], sale)
	}

	transferLookup := make(map[string][]*events.Transfer)
	for _, log := range logs {

		if len(log.Topics) == 0 {
			continue
		}

		if log.Removed {
			continue
		}

		eventHash := log.Topics[0].String()
		switch eventHash {

		case params.HashERCTransfer:

			if len(log.Topics) < 3 && len(log.Topics) > 4 {
				p.log.Warn().
					Uint("index", log.Index).
					Int("topics", len(log.Topics)).
					Msg("skipping log with invalid topic length")
				continue
			}

			switch {

			case len(log.Topics) == 3:

				transfer, err := parsers.ERC20Transfer(log)
				if err != nil {
					return nil, fmt.Errorf("could not parse sale ERC20 transfer: %w", err)
				}

				transferLookup[transfer.TransactionHash] = append(transferLookup[transfer.TransactionHash], transfer)

			case len(log.Topics) == 4:

				transfer, err := parsers.ERC721Transfer(log)
				if err != nil {
					return nil, fmt.Errorf("could not parse sale ERC721 transfer: %w", err)
				}

				transferLookup[transfer.TransactionHash] = append(transferLookup[transfer.TransactionHash], transfer)

			}

		case params.HashERC1155Transfer:

			transfer, err := parsers.ERC1155Transfer(log)
			if err != nil {
				return nil, fmt.Errorf("could not parse ERC1155 transfer: %w", err)
			}

			transferLookup[transfer.TransactionHash] = append(transferLookup[transfer.TransactionHash], transfer)

		default:
			continue
		}
	}

	for transactionHash, sales := range saleLookup {
		// Get the events for this transaction
		transactionTransfers, _ := transferLookup[transactionHash]

		// Sort everything by log index
		sort.Slice(sales, func(i, j int) bool {
			return sales[i].EventIndex < sales[i].EventIndex
		})
		sort.Slice(transactionTransfers, func(i, j int) bool {
			return transactionTransfers[i].EventIndex < transactionTransfers[i].EventIndex
		})

		lastSaleIndex := uint(0)
		// Link transfers to a specific sale
		for _, sale := range sales {

			// Link coins transfers
			transfers := make([]*events.Transfer, 0)
			for _, transfer := range transactionTransfers {
				if transfer.EventIndex >= sale.EventIndex {
					break
				}

				if sale.Hash() != transfer.Hash() || sale.PaymentHash() != transfer.Hash() {
					continue
				}

				if transfer.EventIndex <= lastSaleIndex {
					continue
				}

				transfers = append(transfers, transfer)
			}

			coinTransfer := &events.Transfer{
				CollectionAddress: params.AddressZero,
				TokenCount:        sale.CurrencyValue,
				SenderAddress:     sale.SellerAddress,
				ReceiverAddress:   sale.BuyerAddress,
			}
			if len(coinTransfers) != 0 {
				coinTransfer = coinTransfers[0]
			}

			if coinTransfer.TokenCount != sale.CurrencyValue {
				p.log.Warn().
					Str("sale_id", sale.ID).
					Msg("no erc20 transaction found with the required value, skipping")
				continue
			}

			sale.CurrencyValue = coinTransfer.TokenCount
			sale.CurrencyAddress = coinTransfer.CollectionAddress

			nftTransfer := nftTransfers[0]
			sale.CollectionAddress = nftTransfer.CollectionAddress
			sale.TokenID = nftTransfer.TokenID

			lastSaleIndex = sale.EventIndex
		}
	}

	// Put everything together for the result.
	result := results.Completion{
		Job:      completion,
		Requests: requests,
	}

	return &result, nil
}
