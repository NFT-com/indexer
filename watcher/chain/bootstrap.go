package chain

import (
	"fmt"
	"math/big"
	"time"

	"github.com/NFT-com/indexer/jobs"
)

func (j *Watcher) bootstrap() error {
	startingBlocks, lowestBlock := j.startingBlocks()

	index := lowestBlock
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
				contracts := make([]string, 0, len(j.config.Contracts))
				for _, contract := range j.config.Contracts {
					startingBlock, ok := startingBlocks[contract]
					if !ok {
						j.log.Error().Str("contract", contract).Msg("could not check contract starting block")
						continue
					}

					// means that the current block index is lower that the starting block for this contract
					if index.CmpAbs(startingBlock) < 0 {
						continue
					}

					contracts = append(contracts, contract)
				}

				if len(contracts) == 0 {
					j.log.Info().Msg("could not build contract array to insert in job")
					continue
				}

				job := jobs.Parsing{
					ChainURL:     j.config.ChainURL,
					ChainID:      j.config.ChainID,
					ChainType:    j.config.ChainType,
					BlockNumber:  index.String(),
					Addresses:    contracts,
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

func (j *Watcher) startingBlocks() (map[string]*big.Int, *big.Int) {
	startingBlocks := make(map[string]*big.Int, len(j.config.Contracts))
	lowestBlock := big.NewInt(1)

	for i, contract := range j.config.Contracts {
		blockHeight, ok := j.config.ContractBlockHeights[contract]
		if !ok {
			j.log.Error().Str("contract", contract).Msg("could not create starting height for contract")
			continue
		}

		value, ok := big.NewInt(0).SetString(blockHeight, 0)
		if !ok {
			j.log.Error().Str("contract", contract).Msg("could not convert block height to big.Int")
			continue
		}

		if lowestBlock.CmpAbs(value) > 0 || i == 0 {
			lowestBlock.SetBytes(value.Bytes())
		}

		startingBlocks[contract] = value
	}

	return startingBlocks, lowestBlock
}
