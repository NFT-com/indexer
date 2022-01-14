package events

import "github.com/ethereum/go-ethereum/common"

type Event interface {
	ID() string
	Chain() string
	Network() string
	Topic() string
}

// FIXME: Clean up code below.

type OrdersMatched struct {
	id       string         // Hash(block hash + transaction hash + log index)
	chain    string         // Ethereum
	network  string         // Mainnet
	topic    string         // Transfer
	Address  common.Address // Event Contract Address
	BuyHash  common.Hash
	SellHash common.Hash
	Maker    common.Address
	Taker    common.Address
	Price    uint64
	Metadata common.Hash
}

func NewOrdersMatched(
	id, chain, network, topic string,
	address common.Address,
	buyHash, sellHash common.Hash,
	maker, taker common.Address,
	price uint64,
	metadata common.Hash,
) Event {
	return &OrdersMatched{
		id:       id,
		chain:    chain,
		network:  network,
		topic:    topic,
		Address:  address,
		BuyHash:  buyHash,
		SellHash: sellHash,
		Maker:    maker,
		Taker:    taker,
		Price:    price,
		Metadata: metadata,
	}
}

func (t *OrdersMatched) ID() string {
	return t.id
}

func (t *OrdersMatched) Chain() string {
	return t.chain
}

func (t *OrdersMatched) Network() string {
	return t.network
}

func (t *OrdersMatched) Topic() string {
	return t.topic
}

type Transfer struct {
	id      string         // Hash(block hash + transaction hash + log index)
	chain   string         // Ethereum
	network string         // Mainnet
	topic   string         // Transfer
	Address common.Address // Event Contract Address
	From    common.Address // From address
	To      common.Address // To address
	NftID   uint64         // IF of the NFT
}

func NewTransfer(
	id, chain, network, topic string,
	address, from, to common.Address,
	nftID uint64,
) Event {
	t := Transfer{
		id:      id,
		chain:   chain,
		network: network,
		topic:   topic,
		Address: address,
		From:    from,
		To:      to,
		NftID:   nftID,
	}

	return &t
}

func (t *Transfer) ID() string {
	return t.id
}

func (t *Transfer) Chain() string {
	return t.chain
}

func (t *Transfer) Network() string {
	return t.network
}

func (t *Transfer) Topic() string {
	return t.topic
}
