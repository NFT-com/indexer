package lambdas

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/sha3"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/NFT-com/indexer/models/events"
	"github.com/NFT-com/indexer/models/graph"
	"github.com/NFT-com/indexer/models/inputs"
	"github.com/NFT-com/indexer/models/jobs"
	"github.com/NFT-com/indexer/models/results"
	"github.com/NFT-com/indexer/network/web2"
	"github.com/NFT-com/indexer/network/web3"
	"github.com/NFT-com/indexer/service/parsers"
)

type ActionHandler struct {
	log zerolog.Logger
}

func NewActionHandler(log zerolog.Logger) *ActionHandler {

	a := ActionHandler{
		log: log,
	}

	return &a
}

func (a *ActionHandler) Handle(ctx context.Context, job *jobs.Action) (*results.Action, error) {

	var err error
	var result *results.Action
	switch job.ActionType {

	case jobs.ActionAddition:

		result, err = a.handleAddition(ctx, job)

	case jobs.ActionSaleCollection:

		result, err = a.handleSaleCollection(ctx, job)

	default:
		err = fmt.Errorf("unknown action type (%s)", job.ActionType)
	}

	if err != nil {
		return nil, fmt.Errorf("could not handle job type: %w", err)
	}

	return result, nil
}

func (a *ActionHandler) handleAddition(ctx context.Context, job *jobs.Action) (*results.Action, error) {
	var addition inputs.Addition
	err := json.Unmarshal(job.InputData, &addition)
	if err != nil {
		return nil, fmt.Errorf("could not decode addition inputs: %w", err)
	}

	a.log.Debug().
		Uint64("chain_id", job.ChainID).
		Str("contract_address", addition.ContractAddress).
		Str("token_id", addition.TokenID).
		Uint64("block_height", job.BlockHeight).
		Msg("handling addition job")

	client, err := ethclient.DialContext(ctx, addition.NodeURL)
	if err != nil {
		return nil, fmt.Errorf("could not connect to node: %w", err)
	}
	defer client.Close()

	a.log.Debug().
		Str("node_url", addition.NodeURL).
		Msg("connected to node API")

	fetchURI := web3.NewURIFetcher(client)
	fetchMetadata := web2.NewMetadataFetcher()

	var tokenURI string
	switch addition.Standard {

	case jobs.StandardERC721:

		tokenURI, err = fetchURI.ERC721(ctx, addition.ContractAddress, job.BlockHeight, addition.TokenID)
		if err != nil {
			return nil, fmt.Errorf("could not fetch ERC721 URI: %w", err)
		}

		a.log.Info().
			Str("token_uri", tokenURI).
			Msg("ERC721 token URI retrieved")

	case jobs.StandardERC1155:

		tokenURI, err = fetchURI.ERC1155(ctx, addition.ContractAddress, job.BlockHeight, addition.TokenID)
		if err != nil {
			return nil, fmt.Errorf("could not fetch ERC1155 URI: %w", err)
		}

		a.log.Info().
			Str("token_uri", tokenURI).
			Msg("ERC1155 token URI retrieved")

	default:

		return nil, fmt.Errorf("unknown addition standard (%s)", addition.Standard)
	}

	token, err := fetchMetadata.Token(ctx, tokenURI)
	if err != nil {
		return nil, fmt.Errorf("could not fetch metadata: %w", err)
	}

	a.log.Info().
		Str("name", token.Name).
		Str("description", token.Description).
		Str("image", token.Image).
		Int("attributes", len(token.Attributes)).
		Msg("token metadata fetched")

	nftHash := sha3.Sum256([]byte(fmt.Sprintf("%d-%s-%s", job.ChainID, addition.ContractAddress, addition.TokenID)))
	nftID := uuid.Must(uuid.FromBytes(nftHash[:16]))

	traits := make([]*graph.Trait, 0, len(token.Attributes))
	for i, att := range token.Attributes {
		traitHash := sha3.Sum256([]byte(fmt.Sprintf("%d-%s-%s-%d", job.ChainID, addition.ContractAddress, addition.TokenID, i)))
		traitID := uuid.Must(uuid.FromBytes(traitHash[:16]))
		trait := graph.Trait{
			ID:    traitID.String(),
			NFTID: nftID.String(),
			Name:  att.TraitType,
			Type:  att.DisplayType,
			Value: fmt.Sprint(att.Value),
		}
		traits = append(traits, &trait)
	}

	nft := graph.NFT{
		ID: nftID.String(),
		// CollectionID is populated after parsing
		TokenID:     addition.TokenID,
		Name:        token.Name,
		URI:         tokenURI,
		Image:       token.Image,
		Description: token.Description,
		Owner:       addition.Owner,
		Number:      addition.Number,
	}

	result := results.Action{
		NFT:    &nft,
		Traits: traits,
	}

	return &result, nil
}

func (a *ActionHandler) handleSaleCollection(ctx context.Context, job *jobs.Action) (*results.Action, error) {
	var input inputs.SaleCollection
	err := json.Unmarshal(job.InputData, &input)
	if err != nil {
		return nil, fmt.Errorf("could not decode sale collection inputs: %w", err)
	}

	a.log.Debug().
		Uint64("chain_id", job.ChainID).
		Str("sale_id", input.SaleID).
		Str("transaction_hash", input.TransactionHash).
		Uint64("block_height", job.BlockHeight).
		Msg("handling sale collection job")

	client, err := ethclient.DialContext(ctx, input.NodeURL)
	if err != nil {
		return nil, fmt.Errorf("could not connect to node: %w", err)
	}
	defer client.Close()

	a.log.Debug().
		Str("node_url", input.NodeURL).
		Msg("connected to node API")

	fetchLogs := web3.NewLogsFetcher(client)

	hashes := []string{ERC721TransferHash, ERC1155TransferHash}
	logs, err := fetchLogs.Logs(ctx, nil, hashes, job.BlockHeight, job.BlockHeight)
	if err != nil {
		return nil, fmt.Errorf("could not get block logs: %w", err)
	}

	logs = FilterForTransactionHash(logs, input.TransactionHash)
	if len(logs) == 0 {
		return nil, fmt.Errorf("could not found any trasnfer log for transaction hash: %s", input.TransactionHash)
	}

	if len(logs) > 1 {
		return nil, fmt.Errorf("could found multiple trasnfer logs for transaction hash: %s", input.TransactionHash)
	}

	var transfer *events.Transfer

	topic := logs[0].Topics[0]
	switch topic.String() {

	case ERC721TransferHash:

		transfer, err = parsers.ERC721Transfer(logs[0])

	case ERC1155TransferHash:

		transfer, err = parsers.ERC1155Transfer(logs[0])

	default:
		err = fmt.Errorf("unknow log transfer hash: (%s)", topic.String())

	}

	if err != nil {
		return nil, fmt.Errorf("could not proccess transfer event: %w", err)
	}

	sale := events.Sale{
		ID:                input.SaleID,
		CollectionAddress: transfer.CollectionAddress,
		TokenID:           transfer.TokenID,
	}

	result := results.Action{
		Sale: &sale,
	}

	return &result, nil
}
