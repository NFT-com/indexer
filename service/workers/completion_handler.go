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
	saleLookup := make(map[string][]*events.Sale)
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
		transactionTransfers, ok := transferLookup[transactionHash]
		if !ok {
			log.Warn().
				Str("transaction", transactionHash).
				Msg("no transfers found for transaction")
			continue
		}

		// Sort everything by log index
		sort.Slice(sales, func(i, j int) bool {
			return sales[i].EventIndex < sales[i].EventIndex
		})
		sort.Slice(transactionTransfers, func(i, j int) bool {
			return transactionTransfers[i].EventIndex < transactionTransfers[i].EventIndex
		})

		lastSaleIndex := transactionTransfers[0].EventIndex - 1
		// Link transfers to a specific sale
		for _, sale := range sales {

			// Get all transfers for the current sale
			var transfers []*events.Transfer

			// if there is only one transfer it means native token was used to pay
			if sale.EventIndex-lastSaleIndex == 2 {

				transfers = append(transfers, &events.Transfer{
					CollectionAddress: params.AddressZero,
					TokenCount:        sale.CurrencyValue,
					SenderAddress:     sale.SellerAddress,
					ReceiverAddress:   sale.BuyerAddress,
				})
			}

			for _, transfer := range transactionTransfers {

				if transfer.EventIndex >= sale.EventIndex {
					break
				}

				if sale.Hash() != transfer.Hash() && sale.PaymentHash() != transfer.Hash() {
					continue
				}

				if transfer.EventIndex <= lastSaleIndex {
					continue
				}

				transfers = append(transfers, transfer)
			}

			lastSaleIndex = sale.EventIndex

			// if there is more than 1 nft transfer and 1 token transfer
			// (could be 2 nft transfers next cases will cover this)
			if len(transfers) > 2 {
				continue
			}

			switch len(transfers) {

			case 2:

				var currencyTransfer *events.Transfer
				var nftTransfer *events.Transfer

				switch {

				case transfers[0].TokenCount == sale.CurrencyValue:

					// the transfer with the index 0 is the currency transfer
					currencyTransfer = transfers[0]
					nftTransfer = transfers[1]

				case transfers[1].TokenCount == sale.CurrencyValue:

					// the transfer with the index 1 is the currency transfer
					currencyTransfer = transfers[1]
					nftTransfer = transfers[0]

				default:

					p.log.Warn().
						Str("sale_id", sale.ID).
						Msg("invalid transfers found, skipping")

					continue
				}

				sale.TokenID = nftTransfer.TokenID
				sale.CollectionAddress = nftTransfer.CollectionAddress

				sale.CurrencyAddress = currencyTransfer.CollectionAddress

			case 1:

				p.log.Warn().
					Str("sale_id", sale.ID).
					Msg("not all transfers were found, skipping")

				continue

			default:

				p.log.Warn().
					Str("sale_id", sale.ID).
					Msg("no transfers found, skipping")

				continue
			}
		}
	}

	// Put everything together for the result.
	result := results.Completion{
		Job:      completion,
		Requests: requests,
	}

	return &result, nil
}
