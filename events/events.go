package events

import "github.com/ethereum/go-ethereum/common"

type Event interface {
	String() string
}

// FIXME: Clean up code below.

type Transfer struct {
	ID      string         // Hash(block hash + transaction hash + log index)
	Chain   string         // Ethereum
	Network string         // Mainnet
	Topic   string         // Transfer
	Address common.Address // Event Contract Address
	From    common.Address // From address
	To      common.Address // To address
	NftID   uint64         // ID of the NFT
}

func (t *Transfer) String() string {
	return t.ID
}

type TransferSingle struct {
	ID       string         // Hash(block hash + transaction hash + log index)
	Chain    string         // Ethereum
	Network  string         // Mainnet
	Topic    string         // Transfer
	Address  common.Address // Event Contract Address
	Operator common.Address // The one who started this transaction.
	From     common.Address // From address
	To       common.Address // To address
	NftID    uint64         // ID of the NFT
	Value    uint64         // ID of the NFT
}

func (t *TransferSingle) String() string {
	return t.ID
}

type TransferBatch struct {
	ID       string         // Hash(block hash + transaction hash + log index)
	Chain    string         // Ethereum
	Network  string         // Mainnet
	Topic    string         // Transfer
	Address  common.Address // Event Contract Address
	Operator common.Address // The one who started this transaction.
	From     common.Address // From address
	To       common.Address // To address
	NftIDs   []uint64       // ID of the NFT
	Values   []uint64       // ID of the NFT
}

func (t *TransferBatch) String() string {
	return t.ID
}

type URI struct {
	ID      string         // Hash(block hash + transaction hash + log index)
	Chain   string         // Ethereum
	Network string         // Mainnet
	Topic   string         // Transfer
	Address common.Address // Event Contract Address
	NftID   uint64         // ID of the NFT
	URI     string         // New token URI
}

func (t *URI) String() string {
	return t.ID
}
