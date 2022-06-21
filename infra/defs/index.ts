export const sharedOutputFileName = 'shared-out.json'

export type SharedInfraOutput = {
  assetBucket: string
  assetBucketRole: string
  dbHost: string
  deployIndexerAppBucket: string
  publicSubnets: string[]
  redisHost: string
  vpcId: string
  webSGId: string
  indexerECRRepo: string
  parserFunctionId: string
  additionFunctionId: string
}
