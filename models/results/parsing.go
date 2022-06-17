package results

import (
	"github.com/NFT-com/indexer/models/events"
	"github.com/NFT-com/indexer/models/jobs"
)

type Parsing struct {
	Transfers     []*events.Transfer   `json:"transfers"`
	Sales         []*events.Sale       `json:"sales"`
	Additions     []*jobs.Addition     `json:"additions"`
	Modifications []*jobs.Modification `json:"modifications"`
	Requests      uint                 `json:"requests"`
}
