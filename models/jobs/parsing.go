package jobs

import (
	"github.com/gammazero/deque"
)

// Parsing is a job that parses an NFT's data from block data.
type Parsing struct {
	ID                string   `json:"id"`
	ChainID           uint64   `json:"chain_id"`
	StartHeight       uint64   `json:"start_height"`
	EndHeight         uint64   `json:"end_height"`
	ContractAddresses []string `json:"contract_addresses"`
	EventHashes       []string `json:"event_hashes"`
}

func (p Parsing) Heights() uint {
	return uint(p.StartHeight - p.EndHeight + 1)
}

func (p Parsing) Addresses() uint {
	return uint(len(p.ContractAddresses))
}

func (p *Parsing) Split(heights uint, addresses uint) []*Parsing {

	var parsings []*Parsing
	var queue deque.Deque[*Parsing]
	queue.PushBack(p)
	for queue.Len() != 0 {
		parsing := queue.PopFront()
		if parsing.Heights() > heights {
			pivot := (parsing.StartHeight + parsing.EndHeight) / 2
			left, right := *parsing, *parsing
			left.EndHeight, right.StartHeight = pivot, pivot+1
			queue.PushBack(&left)
			queue.PushBack(&right)
			continue
		}
		if parsing.Addresses() > addresses {
			length := parsing.Addresses()
			pivot := length / 2
			left, right := *parsing, *parsing
			left.ContractAddresses, right.ContractAddresses = left.ContractAddresses[0:pivot], right.ContractAddresses[pivot:length]
			queue.PushBack(&left)
			queue.PushBack(&right)
			continue
		}
		parsings = append(parsings, parsing)
	}
	return parsings
}
