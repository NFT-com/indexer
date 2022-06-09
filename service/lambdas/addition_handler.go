package lambdas

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/sha3"

	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/NFT-com/indexer/models/graph"
	"github.com/NFT-com/indexer/models/inputs"
	"github.com/NFT-com/indexer/models/jobs"
	"github.com/NFT-com/indexer/models/results"
	"github.com/NFT-com/indexer/network/ethereum"
	"github.com/NFT-com/indexer/network/web2"
	"github.com/NFT-com/indexer/network/web3"
)

type AdditionHandler struct {
	log zerolog.Logger
}

func NewAdditionHandler(log zerolog.Logger) *AdditionHandler {

	a := AdditionHandler{
		log: log,
	}

	return &a
}

func (a *AdditionHandler) Handle(ctx context.Context, job *jobs.Action) (*results.Addition, error) {

	var addition inputs.Addition
	err := json.Unmarshal(job.InputData, &addition)
	if err != nil {
		return nil, fmt.Errorf("could not decode addition inputs: %w", err)
	}

	a.log.Debug().
		Uint64("chain_id", job.ChainID).
		Str("contract_address", job.ContractAddress).
		Str("token_id", job.TokenID).
		Uint64("block_height", job.BlockHeight).
		Msg("handling addition job")

	var api *ethclient.Client
	close := func() {}
	if strings.Contains(addition.NodeURL, "ethereum.managedblockchain") {
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not load AWS configuration: %w", err)
		}
		api, close, err = ethereum.NewSigningClient(ctx, addition.NodeURL, cfg)
		if err != nil {
			return nil, fmt.Errorf("could not create signing client (url: %s): %w", addition.NodeURL, err)
		}
	} else {
		api, err = ethclient.DialContext(ctx, addition.NodeURL)
		if err != nil {
			return nil, fmt.Errorf("could not create default client (url: %s): %w", addition.NodeURL, err)
		}
	}
	defer api.Close()
	defer close()

	a.log.Debug().
		Str("node_url", addition.NodeURL).
		Msg("connected to node API")

	fetchURI := web3.NewURIFetcher(api)
	fetchMetadata := web2.NewMetadataFetcher()

	requests := uint(0)
	var tokenURI string
	switch addition.Standard {

	case jobs.StandardERC721:

		requests++
		tokenURI, err = fetchURI.ERC721(ctx, job.ContractAddress, job.TokenID)
		if err != nil && strings.Contains(err.Error(), "nonexistent token") {
			requests++
			tokenURI, err = fetchURI.ERC721Archive(ctx, job.ContractAddress, job.BlockHeight, job.TokenID)
		}
		if err != nil && strings.Contains(err.Error(), "missing trie node") {
			return nil, fmt.Errorf("node does not support archive mode (missing trie node)")
		}
		if err != nil {
			return nil, fmt.Errorf("could not fetch ERC721 URI: %w", err)
		}

		a.log.Info().
			Str("token_uri", tokenURI).
			Msg("ERC721 token URI retrieved")

	case jobs.StandardERC1155:

		requests++
		tokenURI, err = fetchURI.ERC1155(ctx, job.ContractAddress, job.TokenID)
		if err != nil && strings.Contains(err.Error(), "nonexistent token") {
			requests++
			tokenURI, err = fetchURI.ERC1155Archive(ctx, job.ContractAddress, job.BlockHeight, job.TokenID)
		}
		if err != nil && strings.Contains(err.Error(), "missing trie node") {
			return nil, fmt.Errorf("node does not support archive mode (missing trie node)")
		}
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

	nftHash := sha3.Sum256([]byte(fmt.Sprintf("%d-%s-%s", job.ChainID, job.ContractAddress, job.TokenID)))
	nftID := uuid.Must(uuid.FromBytes(nftHash[:16]))

	traits := make([]*graph.Trait, 0, len(token.Attributes))
	for i, att := range token.Attributes {
		traitHash := sha3.Sum256([]byte(fmt.Sprintf("%d-%s-%s-%d", job.ChainID, job.ContractAddress, job.TokenID, i)))
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
		TokenID:     job.TokenID,
		Name:        token.Name,
		URI:         tokenURI,
		Image:       token.Image,
		Description: token.Description,
		Owner:       addition.Owner,
		Number:      addition.Number,
	}

	result := results.Addition{
		NFT:      &nft,
		Traits:   traits,
		Requests: requests,
	}

	return &result, nil
}
