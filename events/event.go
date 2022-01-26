package events

import "github.com/ethereum/go-ethereum/common"

type Event struct {
	ID              string         `json:"id"`
	Chain           string         `json:"chain"`
	Network         string         `json:"network"`
	Block           uint64         `json:"block"`
	TransactionHash common.Hash    `json:"transaction_hash"`
	Address         common.Address `json:"address"`
	Topic           common.Hash    `json:"topic"`
	IndexedData     []common.Hash  `json:"indexed_data"`
	Data            []byte         `json:"data"`
}
