import * as upath from 'upath'
import * as pulumi from '@pulumi/pulumi'
import { deployInfra } from '../helper'
import { createAuroraClustersJob, createAuroraClustersEvent, createAuroraClustersGraph } from './aurora'
import { createRepositories } from './ecr'
import { createCacheClusters } from './elasticache'
import { createBuckets } from './s3'
import { createSecurityGroups } from './security-group'
import { createVPC } from './vpc'
import { createParsingWorker, createAdditionWorker, createCompletionWorker } from './lambda'

const pulumiProgram = async (): Promise<Record<string, any> | void> => {
  const config = new pulumi.Config()
  const zones = config.require('availabilityZones').split(',')

  const vpc = createVPC()
  const sgs = await createSecurityGroups(config, vpc)
  const { main: dbJob } = createAuroraClustersJob(config, vpc, sgs.aurora, zones)
  const { main: dbEvent } = createAuroraClustersEvent(config, vpc, sgs.aurora, zones)
  const { main: dbGraph } = createAuroraClustersGraph(config, vpc, sgs.aurora, zones)
  const { main: cacheMain } = createCacheClusters(config, vpc, sgs.redis, zones)
  const { asset, assetRole, deployApp } = createBuckets()
  const { indexer } = createRepositories()
  const lambdaParser = createParsingWorker()
  const lambdaAddition = createAdditionWorker()
  const lambdaCompletion = createCompletionWorker()

  return {
    assetBucket: asset.bucket,
    assetBucketRole: assetRole.arn,
    jobDbHost: dbJob.endpoint,
    eventDbHost: dbEvent.endpoint,
    graphDbHost: dbGraph.endpoint,
    deployIndexerAppBucket: deployApp.bucket,
    indexerECRRepo: indexer.name,
    redisHost: cacheMain.cacheNodes[0].address,
    publicSubnetIds: vpc.publicSubnetIds,
    vpcId: vpc.id,
    webSGId: sgs.web.id,
    parserFunctionId: lambdaParser.id,
    additionFunctionId: lambdaAddition.id,
    completionFunctionId: lambdaCompletion.id,
  }
}

export const createSharedInfra = (
  preview?: boolean,
): Promise<pulumi.automation.OutputMap> => {
  const stackName = `${process.env.STAGE}.indexer.shared.${process.env.AWS_REGION}`
  const workDir = upath.joinSafe(__dirname, 'stack')
  return deployInfra(stackName, workDir, pulumiProgram, preview)
}
