package chain

import (
	"fmt"
	"math/big"
	"time"

	"github.com/NFT-com/indexer/jobs"
)

func (j *Watcher) bootstrap() error {
	index := j.config.StartingBlock
	for {
		select {
		case <-j.close:
			return nil
		case <-time.After(j.config.BatchDelay):
			if index.CmpAbs(j.latestBlock) > 0 {
				return nil
			}

			batchEnd := big.NewInt(0).Add(index, big.NewInt(j.config.Batch))

			jobs := make([]jobs.Parsing, 0, j.config.Batch)
			for ; index.CmpAbs(j.latestBlock) <= 0 && index.CmpAbs(batchEnd) < 0; index.Add(index, big.NewInt(1)) {
				for _, contract := range j.config.Contracts {
					jobs = append(jobs, j.createJobsForContract(contract, index)...)
				}
			}

			err := j.publishJobs(jobs)
			if err != nil {
				return fmt.Errorf("could not create parsing jobs: %w", err)
			}
		}
	}
}
