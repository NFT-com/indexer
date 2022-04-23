package chain

import (
	"fmt"
	"math/big"
	"time"

	"github.com/NFT-com/indexer/jobs"
)

func (j *Watcher) bootstrap() error {

	j.log.Info().Uint64("start", j.config.StartingBlock.Uint64()).Msg("bootstrapping jobs")

	index := j.config.StartingBlock
	for {
		select {

		case <-j.close:
			j.log.Debug().Msg("bootstrapping aborted")
			return nil

		case <-time.After(j.config.BatchDelay):

			if index.CmpAbs(j.latestBlock) > 0 {
				j.log.Debug().Msg("bootstrapping done")
				return nil
			}

			batchEnd := big.NewInt(0).Add(index, big.NewInt(j.config.Batch))

			jobs := make([]jobs.Parsing, 0, j.config.Batch)
			for ; index.CmpAbs(j.latestBlock) <= 0 && index.CmpAbs(batchEnd) < 0; index.Add(index, big.NewInt(1)) {
				for _, contract := range j.config.Contracts {
					jobs = append(jobs, j.createJobsForContract(contract, index)...)
				}
			}

			j.log.Debug().
				Uint64("start", index.Uint64()).
				Uint64("end", batchEnd.Uint64()).
				Int("jobs", len(jobs)).
				Msg("processing batch")

			err := j.publishJobs(jobs)
			if err != nil {
				return fmt.Errorf("could not create parsing jobs: %w", err)
			}
		}
	}
}
