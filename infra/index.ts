import * as console from 'console'
import fs from 'fs'
import * as process from 'process'
import * as upath from 'upath'
import * as pulumi from '@pulumi/pulumi'
import { SharedInfraOutput, sharedOutputFileName } from './defs'
import { createIndexerServer, updateIndexerEnvFile } from './indexer'
import { createSharedInfra } from './shared'

export const sharedOutToJSONFile = (outMap: pulumi.automation.OutputMap): void => {
  const assetBucket = outMap.assetBucket.value
  const assetBucketRole = outMap.assetBucketRole.value
  const dbHost = outMap.dbHost.value
  const deployIndexerAppBucket = outMap.deployIndexerAppBucket.value
  const indexerECRRepo = outMap.indexerECRRepo.value
  const redisHost = outMap.redisHost.value
  const publicSubnets = outMap.publicSubnetIds.value
  const vpcId = outMap.vpcId.value
  const webSGId = outMap.webSGId.value
  const sharedOutput: SharedInfraOutput = {
    assetBucket,
    assetBucketRole,
    dbHost,
    deployIndexerAppBucket,
    indexerECRRepo,
    redisHost,
    publicSubnets,
    vpcId,
    webSGId,
  }
  const file = upath.joinSafe(__dirname, sharedOutputFileName)
  fs.writeFileSync(file, JSON.stringify(sharedOutput))
}

const main = async (): Promise<any> => {
  const args = process.argv.slice(2)
  const deployShared = args?.[0] === 'deploy:shared' || false
  const deployIndexer = args?.[0] === 'deploy:indexer' || false
  const buildIndexerEnv = args?.[0] === 'indexer:env' || false

  if (deployShared) {
    return createSharedInfra()
      .then(sharedOutToJSONFile)
  }

  if (buildIndexerEnv) {
    updateIndexerEnvFile()
    return
  }

  if (deployIndexer) {
    return createIndexerServer()
  }
}

main()
  .catch((err) => {
    console.error(err)
    process.exit(1)
  })

