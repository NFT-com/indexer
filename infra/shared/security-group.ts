import { ec2 as awsEC2 } from '@pulumi/aws'
import { ec2 } from '@pulumi/awsx'
import * as pulumi from '@pulumi/pulumi'
import { getResourceName, isNotEmpty, isProduction } from '../helper'

export type SGOutput = {
  aurora: awsEC2.SecurityGroup
  redis: awsEC2.SecurityGroup
  web: awsEC2.SecurityGroup
}

const buildIngressRule = (
  port: number,
  protocol = 'tcp',
  sourceSecurityGroupId?: pulumi.Output<string>[],
): any => {
  const rule = {
    protocol,
    fromPort: port,
    toPort: port,
  }
  if (isNotEmpty(sourceSecurityGroupId)) {
    return {
      ...rule,
      securityGroups: sourceSecurityGroupId,
    }
  }

  return {
    ...rule,
    cidrBlocks: new ec2.AnyIPv4Location().cidrBlocks,
    ipv6CidrBlocks: new ec2.AnyIPv6Location().ipv6CidrBlocks,
  }
}

const buildEgressRule = (
  port: number,
  protocol = 'tcp',
): any => ({
  protocol,
  fromPort: port,
  toPort: port,
  cidrBlocks: new ec2.AnyIPv4Location().cidrBlocks,
})

export const createSecurityGroups = (
  config: pulumi.Config, 
  vpc: ec2.Vpc
  ): SGOutput => {
  const web = new awsEC2.SecurityGroup('sg_web', {
    description: 'Allow traffic from/to web',
    name: getResourceName('indexer-web'),
    vpcId: vpc.id,
    ingress: [
      buildIngressRule(22),
      buildIngressRule(4150),
      buildIngressRule(4151),
      buildIngressRule(4160),
      buildIngressRule(4161),
      buildIngressRule(4171),
      buildIngressRule(8080),
    ],
    egress: [
      buildEgressRule(0, '-1'),
    ],
  })

  const aurora = new awsEC2.SecurityGroup('sg_aurora_main', {
    name: getResourceName('indexer-aurora-main'),
    description: 'Allow traffic to Aurora (Postgres) main instance',
    vpcId: vpc.id,
    ingress: [
      isProduction()
        ? buildIngressRule(5432, 'tcp', [web.id])
        : buildIngressRule(5432),
    ],
    egress: [
      buildEgressRule(5432),
    ],
  })

  const redis = new awsEC2.SecurityGroup('sg_redis_main', {
    name: getResourceName('indexer-redis-main'),
    description: 'Allow traffic to Elasticache (Redis) main instance',
    vpcId: vpc.id,
    ingress: [
      isProduction()
        ? buildIngressRule(6379, 'tcp', [web.id])
        : buildIngressRule(6379),
    ],
    egress: [
      buildEgressRule(6379),
    ],
  })

  return {
    aurora,
    redis,
    web,
  }
}
