package event

import "github.com/ethereum/go-ethereum/common"

type Event struct {
	ID              string         `json:"id"`
	Network         string         `json:"network"`
	Chain           string         `json:"chain"`
	Block           uint64         `json:"block"`
	TransactionHash common.Hash    `json:"transaction_hash"`
	Address         common.Address `json:"address"`
	Topic           common.Hash    `json:"topic"`
	IndexedData     []common.Hash  `json:"indexed_data"`
	Data            []byte         `json:"data"`
}

type ParsedEvent struct {
	ID              string `json:"id"`
	Network         string `json:"network"`
	Chain           string `json:"chain"`
	Block           uint64 `json:"block"`
	TransactionHash string `json:"transaction_hash"`
	Address         string `json:"address"`
	Type            string `json:"type"`
	Data            map[string]interface{}
}
