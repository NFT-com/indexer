package jobs

import (
	"github.com/gammazero/deque"
	"github.com/google/uuid"
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
	return uint(p.EndHeight - p.StartHeight + 1)
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
		if parsing.Heights() <= heights && parsing.Addresses() <= addresses {
			parsings = append(parsings, parsing)
			continue
		}
		left, right := *parsing, *parsing
		left.ID, right.ID = uuid.NewString(), uuid.NewString()
		switch {
		case parsing.Heights() > heights:
			pivot := (parsing.StartHeight + parsing.EndHeight) / 2
			left.EndHeight = pivot
			right.StartHeight = pivot + 1
		case parsing.Addresses() > addresses:
			end := parsing.Addresses()
			pivot := end / 2
			left.ContractAddresses = left.ContractAddresses[0:pivot]
			right.ContractAddresses = right.ContractAddresses[pivot:end]
		}
		queue.PushBack(&left)
		queue.PushBack(&right)
	}
	return parsings
}
