package lambdas

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
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
	err := json.Unmarshal(action.Data, &inputs)
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
	switch inputs.EventType {

	case ERC721TransferHash:

		uri, err = fetchURI.ERC721(ctx, action.Address, action.TokenID)
		if err != nil {
			return nil, fmt.Errorf("could not fetch ERC721 URI: %w", err)
		}

	case ERC1155TransferHash, ERC1155BatchHash:

		uri, err = fetchURI.ERC1155(ctx, action.Address, action.TokenID)
		if err != nil {
			return nil, fmt.Errorf("could not fetch ERC1155 URI: %w", err)
		}

	default:

		return nil, fmt.Errorf("unknown addition event type (%s)", inputs.EventType)
	}

	token, err := fetchMetadata.Token(ctx, uri)
	if err != nil {
		return nil, fmt.Errorf("could not fetch metadata: %w", err)
	}

	nftHash := sha3.Sum256([]byte(fmt.Sprintf("%s-%s-%s", action.NetworkID, action.Address, action.TokenID)))
	nftID := hex.EncodeToString(nftHash[:])

	traits := make([]*graph.Trait, 0, len(token.Attributes))
	for i, att := range token.Attributes {
		traitHash := sha3.Sum256([]byte(fmt.Sprintf("%s-%s-%s-%d", action.NetworkID, action.Address, action.TokenID, i)))
		trait := graph.Trait{
			ID:    hex.EncodeToString(traitHash[:]),
			NFTID: nftID,
			Name:  att.TraitType,
			Type:  att.DisplayType,
			Value: fmt.Sprint(att.Value),
		}
		traits = append(traits, &trait)
	}

	nft := graph.NFT{
		ID:           nftID,
		CollectionID: inputs.CollectionID,
		TokenID:      action.TokenID,
		Name:         token.Name,
		URI:          uri,
		Image:        token.Image,
		Description:  token.Description,
		Owner:        inputs.Owner,
	}

	result := results.Addition{
		NFT:    &nft,
		Traits: traits,
	}

	return &result, nil
}
