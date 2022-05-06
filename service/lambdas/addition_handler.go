package lambdas

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/sha3"

	"github.com/NFT-com/indexer/models/graph"
	"github.com/NFT-com/indexer/models/inputs"
	"github.com/NFT-com/indexer/models/jobs"
	"github.com/NFT-com/indexer/models/results"
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

func (a *AdditionHandler) Handle(ctx context.Context, action *jobs.Action) (*results.Addition, error) {

	var inputs inputs.Addition
	err := json.Unmarshal(action.InputData, &inputs)
	if err != nil {
		return nil, fmt.Errorf("could not decode addition inputs: %w", err)
	}

	client, err := ethclient.DialContext(ctx, inputs.NodeURL)
	if err != nil {
		return nil, fmt.Errorf("could not connect to node: %w", err)
	}

	fetchURI := web3.NewURIFetcher(client)
	fetchMetadata := web2.NewMetadataFetcher()

	var uri string
	switch inputs.Standard {

	case jobs.StandardERC721:

		uri, err = fetchURI.ERC721(ctx, action.ContractAddress, action.TokenID)
		if err != nil {
			return nil, fmt.Errorf("could not fetch ERC721 URI: %w", err)
		}

	case jobs.StandardERC1155:

		uri, err = fetchURI.ERC1155(ctx, action.ContractAddress, action.TokenID)
		if err != nil {
			return nil, fmt.Errorf("could not fetch ERC1155 URI: %w", err)
		}

	default:

		return nil, fmt.Errorf("unknown addition standard (%s)", inputs.Standard)
	}

	token, err := fetchMetadata.Token(ctx, uri)
	if err != nil {
		return nil, fmt.Errorf("could not fetch metadata: %w", err)
	}

	nftHash := sha3.Sum256([]byte(fmt.Sprintf("%d-%s-%s", action.ChainID, action.ContractAddress, action.TokenID)))
	nftID := uuid.Must(uuid.FromBytes(nftHash[:16]))

	traits := make([]*graph.Trait, 0, len(token.Attributes))
	for i, att := range token.Attributes {
		traitHash := sha3.Sum256([]byte(fmt.Sprintf("%d-%s-%s-%d", action.ChainID, action.ContractAddress, action.TokenID, i)))
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
		ID:          nftID.String(),
		TokenID:     action.TokenID,
		Name:        token.Name,
		URI:         uri,
		Image:       token.Image,
		Description: token.Description,
		Owner:       inputs.Owner,
	}

	result := results.Addition{
		NFT:    &nft,
		Traits: traits,
	}

	return &result, nil
}
