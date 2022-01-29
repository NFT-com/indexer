package postgres

import (
	"encoding/json"
	"github.com/NFT-com/indexer/event"
)

type Event struct {
	ID              string          `json:"id" bson:"id"`
	Network         string          `json:"network" bson:"network"`
	Chain           string          `json:"chain" bson:"chain"`
	Block           uint64          `json:"block" bson:"block"`
	TransactionHash string          `json:"transaction_hash" bson:"transaction_hash"`
	Address         string          `json:"address" bson:"address"`
	Type            string          `json:"type" bson:"type"`
	Data            json.RawMessage `json:"data" bson:"data"`
}

func FromExternalEvent(event *event.ParsedEvent) (*Event, error) {
	internalEvent := Event{
		ID:              event.ID,
		Network:         event.Network,
		Chain:           event.Chain,
		Block:           event.Block,
		TransactionHash: event.TransactionHash,
		Address:         event.Address,
		Type:            event.Type,
	}

	data, err := json.Marshal(event.Data)
	if err != nil {
		return nil, err
	}

	internalEvent.Data = data
	return &internalEvent, nil
}
