package events

import (
	"time"

	"github.com/NFT-com/indexer/models/id"
	"golang.org/x/crypto/sha3"
)

type Transfer struct {
	ID                string    `json:"id"`
	ChainID           uint64    `json:"chain_id"`
	TokenStandard     string    `json:"token_standard"`
	CollectionAddress string    `json:"collection_address"`
	TokenID           string    `json:"token_id"`
	BlockNumber       uint64    `json:"block_number"`
	TransactionHash   string    `json:"transaction_hash"`
	EventIndex        uint      `json:"event_index"`
	SenderAddress     string    `json:"sender_address"`
	ReceiverAddress   string    `json:"receiver_address"`
	TokenCount        string    `json:"token_count"`
	EmittedAt         time.Time `json:"emitted_at"`
}

func (t Transfer) NFTID() string {
	return id.NFT(t.ChainID, t.CollectionAddress, t.TokenID)
}

func (t Transfer) EventID() string {
	return id.Event(t.TransactionHash, t.EventIndex)
}

func (t Transfer) Hash() string {
	hash := sha3.Sum256([]byte(t.TransactionHash + t.SenderAddress + t.ReceiverAddress))
	return string(hash[:])
}
