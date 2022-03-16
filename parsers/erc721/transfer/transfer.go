package transfer

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/NFT-com/indexer/events"
)

const (
	ZeroValueHash = "0x0000000000000000000000000000000000000000"
)

type Parser struct {
}

func NewParser() *Parser {
	p := Parser{}

	return &p
}

func (p *Parser) ParseRawEvent(rawEvent events.RawEvent) (events.Event, error) {
	switch {
	case rawEvent.IndexData[0] == ZeroValueHash:
		nftID := rawEvent.IndexData[2]
		owner := rawEvent.IndexData[1]

		m := events.Event{
			Type:      events.EventTypeMint,
			ChainID:   rawEvent.ChainID,
			NetworkID: rawEvent.NetworkID,
			NftID:     nftID,
			Contract:  rawEvent.Address,
			ToAddress: owner,
		}

		return m, nil
	case rawEvent.IndexData[1] == ZeroValueHash:
		nftID := rawEvent.IndexData[2]

		m := events.Event{
			Type:      events.EventTypeBurn,
			ChainID:   rawEvent.ChainID,
			NetworkID: rawEvent.NetworkID,
			NftID:     nftID,
			Contract:  rawEvent.Address,
			ToAddress: ZeroValueHash,
		}

		return m, nil
	default:
		nftID := common.HexToHash(rawEvent.IndexData[2]).Big().String()
		owner := common.HexToAddress(rawEvent.IndexData[1]).String()

		m := events.Event{
			Type:      events.EventTypeUpdate,
			ChainID:   rawEvent.ChainID,
			NetworkID: rawEvent.NetworkID,
			NftID:     nftID,
			Contract:  rawEvent.Address,
			ToAddress: owner,
		}

		return m, nil
	}
}
