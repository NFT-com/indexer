package results

import (
	"github.com/NFT-com/indexer/models/events"
	"github.com/NFT-com/indexer/models/jobs"
)

type Completion struct {
	Job      *jobs.Completion `json:"job"`
	Sale     *events.Sale     `json:"sale"`
	Requests uint             `json:"requests"`
}
