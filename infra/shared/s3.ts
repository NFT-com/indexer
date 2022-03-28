import * as aws from '@pulumi/aws'

import { getResourceName, getStage, joinStringsByDash } from '../helper'

export type S3Output = {
  asset: aws.s3.Bucket
  assetRole: aws.iam.Role
  deployApp: aws.s3.Bucket
}

const createAssetRole = (bucketName: string): aws.iam.Role => {
  const inlinePolicy = {
    Version: '2012-10-17',
    Statement: [
      {
        Action: 's3:PutObject',
        Effect: 'Allow',
        Resource: `arn:aws:s3:::${bucketName}/*`,
      },
    ],
  }
  return new aws.iam.Role('role_asset_bucket', {
    name: getResourceName('indexer-asset-bucket.us-east-1'),
    description: 'Role to assume to access Indexer Asset bucket',
    assumeRolePolicy: {
      Version: '2012-10-17',
      Statement: [
        {
          Action: 'sts:AssumeRole',
          Effect: 'Allow',
          Principal: {
            AWS: '*',
          },
        },
      ],
    },
    inlinePolicies: [
      {
        name: 'indexer-access-asset-bucket',
        policy: JSON.stringify(inlinePolicy),
      },
    ],
  })
}

const createAsset = (): { bucket: aws.s3.Bucket; role: aws.iam.Role } => {
  const bucketName = joinStringsByDash('nftcom','indexer', getStage(), 'assets')
  const role = createAssetRole(bucketName)
  const bucket = new aws.s3.Bucket('s3_asset', {
    bucket: bucketName,
    acl: 'private',
    policy: {
      Version: '2012-10-17',
      Statement: [
        {
          Principal: '*',
          Effect: 'Allow',
          Action: 's3:GetObject',
          Resource: `arn:aws:s3:::${bucketName}/*`,
        },
        {
          Principal: { AWS: role.arn },
          Effect: 'Allow',
          Action: 's3:PutObject',
          Resource: `arn:aws:s3:::${bucketName}/*`,
        },
      ],
    },
    corsRules: [
      {
        allowedMethods: ['HEAD', 'GET', 'PUT'],
        allowedOrigins: ['*'],
        allowedHeaders: [
          'amz-sdk-invocation-id',
          'amz-sdk-request',
          'authorization',
          'Authorization',
          'content-type',
          'Content-Type',
          'Referer',
          'User-Agent',
          'x-amz-content-sha256',
          'x-amz-date',
          'x-amz-security-token',
          'x-amz-user-agent',
        ],
        maxAgeSeconds: 3000,
      },
    ],
  })

  return { role, bucket }
}

const createAppDeploy = (): aws.s3.Bucket => {
  const bucketName = joinStringsByDash('nftcom','indexer', getStage(), 'deploy-app')
  return new aws.s3.Bucket('s3_deploy_app', {
    bucket: bucketName,
    acl: 'private',
  })
}

export const createBuckets = (): S3Output => {
  const { bucket: asset, role: assetRole } = createAsset()
  const deployApp = createAppDeploy()
  return { asset, assetRole, deployApp }
}
