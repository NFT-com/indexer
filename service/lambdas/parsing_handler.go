package lambdas

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/NFT-com/indexer/models/events"
	"github.com/NFT-com/indexer/models/inputs"
	"github.com/NFT-com/indexer/models/jobs"
	"github.com/NFT-com/indexer/models/results"
	"github.com/NFT-com/indexer/network/web3"
	"github.com/NFT-com/indexer/service/parsers"
)

const (
	zeroAddress = "0x0000000000000000000000000000000000000000"
)

type ParsingHandler struct {
	log    zerolog.Logger
	client *http.Client
}

func NewParsingHandler(log zerolog.Logger, client *http.Client) *ParsingHandler {

	e := ParsingHandler{
		log:    log,
		client: client,
	}

	return &e
}

func (p *ParsingHandler) Handle(ctx context.Context, job *jobs.Parsing) (*results.Parsing, error) {

	var parsing inputs.Parsing
	err := json.Unmarshal(job.InputData, &parsing)
	if err != nil {
		return nil, fmt.Errorf("could not decode parsing inputs: %w", err)
	}

	p.log.Debug().
		Uint64("chain_id", job.ChainID).
		Strs("contract_addresses", job.ContractAddresses).
		Strs("event_hashes", job.EventHashes).
		Uint64("start_height", job.StartHeight).
		Uint64("end_height", job.EndHeight).
		Msg("handling parsing job")

	rpc, err := rpc.DialHTTPWithClient(parsing.NodeURL, p.client)
	if err != nil {
		return nil, fmt.Errorf("could not connect to node: %w", err)
	}

	client := ethclient.NewClient(rpc)
	defer client.Close()

	p.log.Debug().
		Str("node_url", parsing.NodeURL).
		Msg("connected to node API")

	fetch := web3.NewLogsFetcher(client)

	// Retrieve the logs for all of the addresses and event types for the given block range.
	logs, err := fetch.Logs(ctx, job.ContractAddresses, job.EventHashes, job.StartHeight, job.EndHeight)
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
	standards := make(map[string]string)
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

		case ERC721TransferHash:

			transfer, err := parsers.ERC721Transfer(log)
			if err != nil {
				return nil, fmt.Errorf("could not parse ERC721 transfer: %w", err)
			}
			transfers = append(transfers, transfer)
			standards[transfer.ID] = jobs.StandardERC721

			p.log.Trace().
				Str("transaction", log.TxHash.Hex()).
				Uint("index", log.Index).
				Str("collection_address", transfer.CollectionAddress).
				Str("token_id", transfer.TokenID).
				Str("sender_address", transfer.SenderAddress).
				Str("receiver_address", transfer.ReceiverAddress).
				Uint("token_count", transfer.TokenCount).
				Msg("ERC721 transfer parsed")

		case ERC1155TransferHash:

			transfer, err := parsers.ERC1155Transfer(log)
			if err != nil {
				return nil, fmt.Errorf("could not parse ERC1155 transfer: %w", err)
			}
			transfers = append(transfers, transfer)
			standards[transfer.ID] = jobs.StandardERC1155

			p.log.Trace().
				Str("transaction", log.TxHash.Hex()).
				Uint("index", log.Index).
				Str("collection_address", transfer.CollectionAddress).
				Str("token_id", transfer.TokenID).
				Str("sender_address", transfer.SenderAddress).
				Str("receiver_address", transfer.ReceiverAddress).
				Uint("token_count", transfer.TokenCount).
				Msg("ERC1155 transfer parsed")

		case ERC1155BatchHash:

			batch, err := parsers.ERC1155Batch(log)
			if err != nil {
				return nil, fmt.Errorf("could not parse ERC1155 batch: %w", err)
			}
			transfers = append(transfers, batch...)
			for _, transfer := range batch {

				standards[transfer.ID] = jobs.StandardERC1155

				p.log.Trace().
					Str("transaction", log.TxHash.Hex()).
					Uint("index", log.Index).
					Str("collection_address", transfer.CollectionAddress).
					Str("token_id", transfer.TokenID).
					Str("sender_address", transfer.SenderAddress).
					Str("receiver_address", transfer.ReceiverAddress).
					Uint("token_count", transfer.TokenCount).
					Msg("ERC115 batch parsed")
			}

		case OpenSeaTradeHash:

			sale, err := parsers.OpenSeaSale(log)
			if err != nil {
				return nil, fmt.Errorf("could not parse OpenSea sale: %w", err)
			}
			sales = append(sales, sale)

			p.log.Trace().
				Str("transaction", log.TxHash.Hex()).
				Uint("index", log.Index).
				Msg("OpenSea sale parsed")
		}
	}

	p.log.Info().
		Int("logs", len(logs)).
		Int("transfers", len(transfers)).
		Int("sales", len(sales)).
		Msg("all logs parsed")

	// Get all the headers to assign timestamps to the events.
	for height := range timestamps {

		header, err := client.HeaderByNumber(ctx, big.NewInt(0).SetUint64(height))
		if err != nil {
			return nil, fmt.Errorf("could not get header for height (%d): %w", height, err)
		}

		timestamps[height] = time.Unix(int64(header.Time), 0)
	}

	p.log.Info().
		Int("heights", len(timestamps)).
		Msg("block heights retrieved")

	// Go through all logs and assign timestamp of emission
	for _, transfer := range transfers {
		transfer.ChainID = job.ChainID
		transfer.EmittedAt = timestamps[transfer.BlockNumber]
	}
	for _, sale := range sales {
		sale.ChainID = job.ChainID
		sale.EmittedAt = timestamps[sale.BlockNumber]
	}

	// Go through all transfers and convert them to mints/burns where appropriate.
	var actions []*jobs.Action
	for _, transfer := range transfers {
		switch {

		case transfer.SenderAddress == zeroAddress:

			inputs := inputs.Addition{
				NodeURL:  parsing.NodeURL,
				Standard: standards[transfer.ID],
				Owner:    transfer.ReceiverAddress,
				Number:   transfer.TokenCount,
			}
			data, err := json.Marshal(inputs)
			if err != nil {
				return nil, fmt.Errorf("could not encode addition inputs: %w", err)
			}
			action := jobs.Action{
				ID:              uuid.NewString(),
				ChainID:         transfer.ChainID,
				ContractAddress: transfer.CollectionAddress,
				TokenID:         transfer.TokenID,
				ActionType:      jobs.ActionAddition,
				BlockHeight:     transfer.BlockNumber,
				JobStatus:       jobs.StatusCreated,
				InputData:       data,
			}
			actions = append(actions, &action)

		default:

			inputs := inputs.OwnerChange{
				PrevOwner: transfer.SenderAddress,
				NewOwner:  transfer.ReceiverAddress,
				Number:    transfer.TokenCount,
			}
			data, err := json.Marshal(inputs)
			if err != nil {
				return nil, fmt.Errorf("could not encode owner change inputs: %w", err)
			}
			action := jobs.Action{
				ID:              uuid.NewString(),
				ChainID:         transfer.ChainID,
				ContractAddress: transfer.CollectionAddress,
				TokenID:         transfer.TokenID,
				ActionType:      jobs.ActionOwnerChange,
				BlockHeight:     transfer.BlockNumber,
				JobStatus:       jobs.StatusCreated,
				InputData:       data,
			}
			actions = append(actions, &action)
		}
	}

	p.log.Info().
		Int("actions", len(actions)).
		Msg("downstream actions created")

	// Put everything together for the result.
	result := results.Parsing{
		Sales:     sales,
		Transfers: transfers,
		Actions:   actions,
	}

	return &result, nil
}
