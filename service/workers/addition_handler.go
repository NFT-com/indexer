package workers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ipfs/go-cid"
	"github.com/rs/zerolog"

	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/NFT-com/indexer/models/content"
	"github.com/NFT-com/indexer/models/gateway"
	"github.com/NFT-com/indexer/models/graph"
	"github.com/NFT-com/indexer/models/jobs"
	"github.com/NFT-com/indexer/models/metadata"
	"github.com/NFT-com/indexer/models/protocol"
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

	log.Debug().
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

	fetchURI := web3.NewURIFetcher(api)
	fetchMetadata := web2.NewMetadataFetcher(web2.WithDisableValidation(true))

	var tokenURI string
	switch addition.TokenStandard {

	case jobs.StandardERC721:

		tokenURI, err = fetchURI.ERC721(ctx, addition.ContractAddress, addition.TokenID)
		if err != nil {
			return nil, fmt.Errorf("could not fetch ERC721 URI: %w", err)
		}

		log.Debug().
			Str("token_uri", tokenURI).
			Msg("ERC721 token URI retrieved")

	case jobs.StandardERC1155:

		tokenURI, err = fetchURI.ERC1155(ctx, addition.ContractAddress, addition.TokenID)
		if err != nil {
			return nil, fmt.Errorf("could not fetch ERC1155 URI: %w", err)
		}

		log.Debug().
			Str("token_uri", tokenURI).
			Msg("ERC1155 token URI retrieved")

	default:

		return nil, fmt.Errorf("unknown token standard (%s)", addition.TokenStandard)
	}

	// Trim any spaces.
	tokenURI = strings.TrimSpace(tokenURI)
	if tokenURI == "" {
		return nil, fmt.Errorf("token URI empty")
	}

	// First, we check if the URI starts with a CID hash, in which case we add the IPFS prefix.
	prefixedURI := tokenURI
	parts := strings.Split(prefixedURI, "/")
	first := parts[0]
	_, err = cid.Decode(first)
	if err == nil {
		prefixedURI = protocol.IPFS + prefixedURI
		log.Debug().
			Str("prefixed_uri", prefixedURI).
			Msg("CID hash prefixed")
	}

	// Then, we substitute the known protocols with known public gateways.
	publicURI := prefixedURI
	switch {

	case strings.HasPrefix(publicURI, protocol.IPFS):
		publicURI = gateway.IPFS + strings.TrimPrefix(publicURI, protocol.IPFS)
		log.Debug().
			Str("public_uri", publicURI).
			Msg("IPFS gateway substituted")

	case strings.HasPrefix(publicURI, protocol.ARWeave):
		publicURI = gateway.ARWeave + strings.TrimPrefix(publicURI, protocol.ARWeave)
		log.Debug().
			Str("public_uri", publicURI).
			Msg("ARWeave gateway substituted")
	}

	// Then, we substitute known public gateways with our own private address.
	privateURI := publicURI
	switch {

	case strings.HasPrefix(privateURI, gateway.IPFS):
		privateURI = gateway.Immutable + strings.TrimPrefix(privateURI, gateway.IPFS)
		log.Debug().
			Str("private_uri", privateURI).
			Msg("IPFS gateway replaced")

	case strings.HasPrefix(privateURI, gateway.Pinata):
		privateURI = gateway.Immutable + strings.TrimPrefix(privateURI, gateway.Pinata)
		log.Debug().
			Str("private_uri", privateURI).
			Msg("Pinata gateway replaced")

	case strings.HasPrefix(privateURI, gateway.Parallel):
		privateURI = gateway.Immutable + strings.TrimPrefix(privateURI, gateway.Parallel)
		log.Debug().
			Str("private_uri", privateURI).
			Msg("Parallel gateway replaced")
	}

	// Finally, we check if we have a payload already, or if we need to fetch it remotely.
	var payload []byte
	var code int
	switch {

	case strings.HasPrefix(privateURI, protocol.HTTP), strings.HasPrefix(privateURI, protocol.HTTPS):
		payload, code, err = fetchMetadata.Payload(ctx, privateURI)
		if err != nil {
			return nil, fmt.Errorf("could not fetch remote metadata: %w", err)
		}
		if code == http.StatusInternalServerError && isTokenNotFound(payload) {
			return nil, results.ErrTokenNotFound
		}
		if code == http.StatusNotFound && isIPFSMissingLink(payload) {
			return nil, results.ErrTokenNotFound
		}
		log.Debug().
			Str("payload", string(payload)).
			Msg("remote payload fetched")

	case strings.HasPrefix(privateURI, content.UTF8):
		payload = []byte(strings.TrimPrefix(privateURI, content.UTF8))
		log.Debug().
			Str("payload", string(payload)).
			Msg("UTF-8 payload trimmed")

	case strings.HasPrefix(privateURI, content.ASCII):
		payload = []byte(strings.TrimPrefix(privateURI, content.ASCII))
		log.Debug().
			Str("payload", string(payload)).
			Msg("ASCII payload trimmed")

	case strings.HasPrefix(privateURI, content.Base64):
		payload, err = base64.StdEncoding.DecodeString(strings.TrimPrefix(privateURI, content.Base64))
		if err != nil {
			return nil, fmt.Errorf("could not decode base64 metadata: %w", err)
		}
		log.Debug().
			Str("payload", string(payload)).
			Msg("Base64 payload decoded")

	default:
		return nil, fmt.Errorf("unknown URI format")
	}

	var token metadata.Token
	err = json.Unmarshal(payload, &token)
	if err != nil {
		return nil, fmt.Errorf("could not decode metadata: %w", err)
	}

	log.Info().
		Str("uri", tokenURI).
		Str("prefixed", prefixedURI).
		Str("public", publicURI).
		Str("private", privateURI).
		Str("payload", string(payload)).
		Str("name", token.Name).
		Str("description", token.Description).
		Str("image", token.Image).
		Int("attributes", len(token.Attributes)).
		Msg("token metadata extracted")

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

func isTokenNotFound(payload []byte) bool {
	var reqErr *results.Error
	err := json.Unmarshal(payload, &reqErr)
	if err != nil {
		return false
	}
	if reqErr.Error() == "Token not found" {
		return true
	}
	return false
}

func isIPFSMissingLink(payload []byte) bool {
	var reqErr *results.Error
	err := json.Unmarshal(payload, &reqErr)
	if err != nil {
		return false
	}
	if strings.Contains(reqErr.Error(), "URI query for nonexistent token") ||
		strings.Contains(reqErr.Error(), "no link named") {
		return true
	}
	return false
}
