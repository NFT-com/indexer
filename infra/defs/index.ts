export const sharedOutputFileName = 'shared-out.json'

export type SharedInfraOutput = {
  jobDbHost: string
  eventDbHost: string
  graphDbHost: string
  publicSubnets: string[]
  redisHost: string
  vpcId: string
  webSGId: string
  indexerECRRepo: string
  parserFunctionId: string
  additionFunctionId: string
  completionFunctionId: string
}
