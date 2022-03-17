import * as aws from '@pulumi/aws'

import { getResourceName } from '../helper'

export type RepositoryOut = {
  indexer: aws.ecr.Repository
}

export const createIndexerRepository = (): aws.ecr.Repository => {
  return new aws.ecr.Repository('ecr_indexer', {
    name: getResourceName('indexer'),
    imageScanningConfiguration: {
      scanOnPush: true,
    },
  })
}

export const createRepositories = (): RepositoryOut => {
  const IndexerRepo = createIndexerRepository()
  return {
    indexer: IndexerRepo,
  }
}
