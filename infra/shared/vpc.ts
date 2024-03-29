import { ec2 } from '@pulumi/awsx'
import { getStage } from '../helper'

export const createVPC = (): ec2.Vpc => {
  const stage = getStage()
  return new ec2.Vpc('vpc', {
    cidrBlock: process.env.VPC_CIDR,
    numberOfAvailabilityZones: 3,
    numberOfNatGateways: 0,
    subnets: [
      { type: 'public', name: stage },
      { type: 'private', name: stage },
    ],
  })
}
