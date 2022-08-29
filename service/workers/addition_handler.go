package workers

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog"

	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/NFT-com/indexer/models/graph"
	"github.com/NFT-com/indexer/models/jobs"
	"github.com/NFT-com/indexer/models/results"
	"github.com/NFT-com/indexer/network/ethereum"
	"github.com/NFT-com/indexer/network/web2"
	"github.com/NFT-com/indexer/network/web3"
)

type AdditionHandler struct {
	log zerolog.Logger
	url string
}

func NewAdditionHandler(log zerolog.Logger, url string) *AdditionHandler {

	a := AdditionHandler{
		log: log,
		url: url,
	}

	return &a
}

func (a *AdditionHandler) Handle(ctx context.Context, addition *jobs.Addition) (*results.Addition, error) {

	log := a.log.With().
		Str("job_id", addition.ID).
		Uint64("chain_id", addition.ChainID).
		Str("contract_address", addition.ContractAddress).
		Str("token_id", addition.TokenID).
		Str("token_standard", addition.TokenStandard).
		Logger()

	log.Info().
		Str("node_url", a.url).
		Msg("initiating connection to Ethereum node")

	var err error
	var api *ethclient.Client
	close := func() {}
	if strings.Contains(a.url, "ethereum.managedblockchain") {
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not load AWS configuration: %w", err)
		}
		api, close, err = ethereum.NewSigningClient(ctx, a.url, cfg)
		if err != nil {
			return nil, fmt.Errorf("could not create signing client (url: %s): %w", a.url, err)
		}
	} else {
		api, err = ethclient.DialContext(ctx, a.url)
		if err != nil {
			return nil, fmt.Errorf("could not create default client (url: %s): %w", a.url, err)
		}
	}
	defer api.Close()
	defer close()

	log.Info().
		Str("node_url", a.url).
		Msg("connection to Ethereum node established")

	fetchURI := web3.NewURIFetcher(api)
	fetchMetadata := web2.NewMetadataFetcher()

	var tokenURI string
	switch addition.TokenStandard {

	case jobs.StandardERC721:

		tokenURI, err = fetchURI.ERC721(ctx, addition.ContractAddress, addition.TokenID)
		if err != nil {
			return nil, fmt.Errorf("could not fetch ERC721 URI: %w", err)
		}

		a.log.Info().
			Str("token_uri", tokenURI).
			Msg("ERC721 token URI retrieved")

	case jobs.StandardERC1155:

		tokenURI, err = fetchURI.ERC1155(ctx, addition.ContractAddress, addition.TokenID)
		if err != nil {
			return nil, fmt.Errorf("could not fetch ERC1155 URI: %w", err)
		}

		a.log.Info().
			Str("token_uri", tokenURI).
			Msg("ERC1155 token URI retrieved")

	default:

		return nil, fmt.Errorf("unknown token standard (%s)", addition.TokenStandard)
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

	nftID := addition.NFTID()
	traits := make([]*graph.Trait, 0, len(token.Attributes))
	for i, att := range token.Attributes {

		var value string
		switch {
		case att.Value != nil:
			value = fmt.Sprint(att.Value)
		case att.TraitValue != nil:
			value = fmt.Sprint(att.TraitValue)
		}

		traitID := addition.TraitID(uint(i))
		trait := graph.Trait{
			ID:    traitID,
			NFTID: nftID,
			Name:  att.TraitType,
			Type:  att.DisplayType,
			Value: value,
		}
		traits = append(traits, &trait)
	}

	nft := graph.NFT{
		ID:           nftID,
		CollectionID: addition.CollectionID,
		TokenID:      addition.TokenID,
		Name:         token.Name,
		URI:          tokenURI,
		Image:        token.Image,
		Description:  token.Description,
	}

	result := results.Addition{
		Job:    addition,
		NFT:    &nft,
		Traits: traits,
	}

	return &result, nil
}
