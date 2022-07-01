package results

import (
	"github.com/NFT-com/indexer/models/events"
	"github.com/NFT-com/indexer/models/jobs"
)

type Parsing struct {
	Job       *jobs.Parsing      `json:"job"`
	Transfers []*events.Transfer `json:"transfers"`
	Sales     []*events.Sale     `json:"sales"`
	Requests  uint               `json:"requests"`
}
