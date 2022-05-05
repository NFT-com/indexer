package lambdas

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/models/events"
	"github.com/NFT-com/indexer/models/inputs"
	"github.com/NFT-com/indexer/models/jobs"
	"github.com/NFT-com/indexer/models/results"
	"github.com/NFT-com/indexer/network/web3"
	"github.com/NFT-com/indexer/service/parsers"
)

type ParsingHandler struct {
	log zerolog.Logger
}

func NewParsingHandler(log zerolog.Logger) *ParsingHandler {

	e := ParsingHandler{
		log: log,
	}

	return &e
}

func (p *ParsingHandler) Handle(ctx context.Context, job *jobs.Parsing) (*results.Parsing, error) {

	var parsing inputs.Parsing
	err := json.Unmarshal(job.Data, &parsing)
	if err != nil {
		return nil, fmt.Errorf("could not decode parsing inputs: %w", err)
	}

	client, err := ethclient.DialContext(ctx, parsing.NodeURL)
	if err != nil {
		return nil, fmt.Errorf("could not connect to node: %w", err)
	}
	defer client.Close()

	fetch := web3.NewLogsFetcher(client)

	// Retrieve the logs for all of the addresses and event types for the given block range.
	logs, err := fetch.Logs(ctx, job.Addresses, job.EventTypes, job.StartHeight, job.EndHeight)
	if err != nil {
		return nil, fmt.Errorf("could not fetch logs: %w", err)
	}

	// For each log, try to parse it into the respective events.
	var transfers []*events.Transfer
	var sales []*events.Sale
	timestamps := make(map[uint64]time.Time)
	standards := make(map[string]string)
	for _, log := range logs {

		// skip logs for reverted transactions
		if log.Removed {
			continue
		}

		// keep track of all heightSet we need to process to get timestamps
		timestamps[log.BlockNumber] = time.Time{}

		eventType := log.Topics[0]
		switch eventType.String() {

		case ERC721TransferHash:

			transfer, err := parsers.ERC721Transfer(log)
			if err != nil {
				return nil, fmt.Errorf("could not parse ERC721 transfer: %w", err)
			}
			transfers = append(transfers, transfer)
			standards[transfer.ID] = jobs.StandardERC721

		case ERC1155TransferHash:

			transfer, err := parsers.ERC1155Transfer(log)
			if err != nil {
				return nil, fmt.Errorf("could not parse ERC1155 transfer: %w", err)
			}
			transfers = append(transfers, transfer)
			standards[transfer.ID] = jobs.StandardERC1155

		case ERC1155BatchHash:

			batch, err := parsers.ERC1155Batch(log)
			if err != nil {
				return nil, fmt.Errorf("could not parse ERC1155 batch: %w", err)
			}
			transfers = append(transfers, batch...)
			for _, transfer := range batch {
				standards[transfer.ID] = jobs.StandardERC1155
			}

		case OpenSeaTradeHash:

			sale, err := parsers.OpenSeaSale(log)
			if err != nil {
				return nil, fmt.Errorf("could not parse OpenSea sale: %w", err)
			}
			sales = append(sales, sale)
		}
	}

	// Get all of the headers to assign timestamps to the events.
	for height := range timestamps {

		header, err := client.HeaderByNumber(ctx, big.NewInt(0).SetUint64(height))
		if err != nil {
			return nil, fmt.Errorf("could not get header for height (%d): %w", height, err)
		}

		timestamps[height] = time.Unix(int64(header.Time), 0)
	}

	// Go through all logs and assign timestamp of emission
	for _, transfer := range transfers {
		transfer.EmittedAt = timestamps[transfer.BlockNumber]
	}
	for _, sale := range sales {
		sale.EmittedAt = timestamps[sale.BlockNumber]
	}

	// Go through all transfers and convert them to mints/burns where appropriate.
	var actions []*jobs.Action
	for _, transfer := range transfers {
		switch {

		case transfer.FromAddress == ZeroAddress:

			inputs := inputs.Addition{
				NodeURL:  parsing.NodeURL,
				Standard: standards[transfer.ID],
				Owner:    transfer.ToAddress,
			}
			data, err := json.Marshal(inputs)
			if err != nil {
				return nil, fmt.Errorf("could not encode addition inputs: %w", err)
			}
			action := jobs.Action{
				ID:         uuid.New().String(),
				ChainID:    transfer.ChainID,
				Address:    transfer.CollectionAddress,
				TokenID:    transfer.TokenID,
				ActionType: jobs.ActionAddition,
				Height:     transfer.BlockNumber,
				Data:       data,
				Status:     jobs.StatusCreated,
			}
			actions = append(actions, &action)

		default:

			inputs := inputs.OwnerChange{
				NewOwner: transfer.ToAddress,
			}
			data, err := json.Marshal(inputs)
			if err != nil {
				return nil, fmt.Errorf("could not encode owner change inputs: %w", err)
			}
			action := jobs.Action{
				ID:         uuid.New().String(),
				ChainID:    transfer.ChainID,
				Address:    transfer.CollectionAddress,
				TokenID:    transfer.TokenID,
				ActionType: jobs.ActionOwnerChange,
				Height:     transfer.BlockNumber,
				Data:       data,
				Status:     jobs.StatusCreated,
			}
			actions = append(actions, &action)
		}

	}

	// Put everything together for the result.
	result := results.Parsing{
		Sales:     sales,
		Transfers: transfers,
		Actions:   actions,
	}

	return &result, nil
}
