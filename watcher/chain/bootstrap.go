package chain

import (
	"fmt"
	"math/big"
	"time"

	"github.com/NFT-com/indexer/jobs"
)

func (j *Watcher) bootstrap() error {
	startingBlock, ok := big.NewInt(0).SetString(j.config.StartIndex, 0)
	if !ok {
		return fmt.Errorf("could not parse block number into big.Int")
	}

	index := startingBlock
	for {
		select {
		case <-j.close:
			return nil
		case <-time.After(j.config.BatchDelay):
			if index.CmpAbs(j.latestBlock) > 0 {
				return nil
			}

			batchEnd := big.NewInt(0).Add(index, big.NewInt(j.config.Batch))

			jobList := make([]jobs.Parsing, 0, j.config.Batch)
			for ; index.CmpAbs(j.latestBlock) <= 0 && index.CmpAbs(batchEnd) < 0; index.Add(index, big.NewInt(1)) {
				job := jobs.Parsing{
					ChainURL:     j.config.ChainURL,
					ChainType:    j.config.ChainType,
					BlockNumber:  index.String(),
					Address:      j.config.Contract,
					StandardType: j.config.StandardType,
					EventType:    j.config.EventType,
				}

				jobList = append(jobList, job)
			}

			err := j.apiClient.CreateParsingJobs(jobList)
			if err != nil {
				return fmt.Errorf("could not create parsing jobs: %w", err)
			}
		}
	}
}
