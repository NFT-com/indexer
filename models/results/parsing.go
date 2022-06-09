package results

import (
	"github.com/NFT-com/indexer/models/events"
	"github.com/NFT-com/indexer/models/jobs"
)

type Parsing struct {
	Transfers []*events.Transfer `json:"transfers"`
	Sales     []*events.Sale     `json:"sales"`
	Actions   []*jobs.Action     `json:"actions"`
	Requests  uint               `json:"requests"`
}
