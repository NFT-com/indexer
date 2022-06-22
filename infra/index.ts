import * as console from 'console'
import fs from 'fs'
import * as process from 'process'
import * as upath from 'upath'
import * as pulumi from '@pulumi/pulumi'
import { SharedInfraOutput, sharedOutputFileName } from './defs'
import { createSharedInfra } from './shared'
import { createIndexerEcsCluster } from './indexer'

export const sharedOutToJSONFile = (outMap: pulumi.automation.OutputMap): void => {
  const assetBucket = outMap.assetBucket.value
  const assetBucketRole = outMap.assetBucketRole.value
  const jobDbHost = outMap.jobDbHost.value
  const eventDbHost = outMap.eventDbHost.value
  const graphDbHost = outMap.graphDbHost.value
  const deployIndexerAppBucket = outMap.deployIndexerAppBucket.value
  const indexerECRRepo = outMap.indexerECRRepo.value
  const redisHost = outMap.redisHost.value
  const publicSubnets = outMap.publicSubnetIds.value
  const vpcId = outMap.vpcId.value
  const webSGId = outMap.webSGId.value
  const parserFunctionId = outMap.parserFunctionId.value
  const additionFunctionId = outMap.additionFunctionId.value
  const sharedOutput: SharedInfraOutput = {
    assetBucket,
    assetBucketRole,
    jobDbHost,
    eventDbHost,
    graphDbHost,
    deployIndexerAppBucket,
    indexerECRRepo,
    redisHost,
    publicSubnets,
    vpcId,
    webSGId,
    parserFunctionId,
    additionFunctionId,
  }
  const file = upath.joinSafe(__dirname, sharedOutputFileName)
  fs.writeFileSync(file, JSON.stringify(sharedOutput))
}

const main = async (): Promise<any> => {
  const args = process.argv.slice(2)
  const deployShared = args?.[0] === 'deploy:shared' || false
  const deployIndexer = args?.[0] === 'deploy:indexer' || false

  if (deployShared) {
    return createSharedInfra(true)
      .then(sharedOutToJSONFile)
  }
  
  if (deployIndexer) {
    return createIndexerEcsCluster()
  }
}

main()
  .catch((err) => {
    console.error(err)
    process.exit(1)
  })

