package results

import (
	"github.com/NFT-com/indexer/models/jobs"
)

type Completion struct {
	Job *jobs.Completion `json:"job"`
}
