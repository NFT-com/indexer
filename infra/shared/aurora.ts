import * as aws from '@pulumi/aws'
import { EngineType } from '@pulumi/aws/types/enums/rds'
import { ec2 } from '@pulumi/awsx'
import * as pulumi from '@pulumi/pulumi'
import { getResourceName, isFalse, isProduction } from '../helper'

export type AuroraOutput = {
  main: aws.rds.Cluster
}

const getSubnetGroupJob = (vpc: ec2.Vpc): aws.rds.SubnetGroup => {
  return new aws.rds.SubnetGroup('aurora_subnet_group_job', {
    name: getResourceName('indexer-job-aurora'),
    subnetIds: isProduction() ? vpc.privateSubnetIds : vpc.publicSubnetIds,
  })
}

const getSubnetGroupEvent = (vpc: ec2.Vpc): aws.rds.SubnetGroup => {
  return new aws.rds.SubnetGroup('aurora_subnet_group_event', {
    name: getResourceName('indexer-graph-aurora'),
    subnetIds: isProduction() ? vpc.privateSubnetIds : vpc.publicSubnetIds,
  })
}

const getSubnetGroupGraph = (vpc: ec2.Vpc): aws.rds.SubnetGroup => {
  return new aws.rds.SubnetGroup('aurora_subnet_group_graph', {
    name: getResourceName('indexer-event-aurora'),
    subnetIds: isProduction() ? vpc.privateSubnetIds : vpc.publicSubnetIds,
  })
}

// https://github.com/pulumi/pulumi-aws-quickstart-aurora-postgres/blob/master/provider/pkg/provider/postgresql.go
const createMainJobDb = (
  config: pulumi.Config,
  vpc: ec2.Vpc,
  sg: aws.ec2.SecurityGroup,
  zones: string[],
): aws.rds.Cluster => {
  const paramFamily = 'aurora-postgresql13'
  const clusterParameterGroup = new aws.rds.ClusterParameterGroup('aurora_job_cluster_param_group', {
    name: getResourceName('indexer-job-cluster'),
    family: paramFamily,
    parameters: [
      {
        name: 'rds.force_ssl',
        value: '1',
        applyMethod: 'pending-reboot',
      },
    ],
  })

  const subnetGroup = getSubnetGroupJob(vpc)
  const engineType = EngineType.AuroraPostgresql
  const cluster = new aws.rds.Cluster('aurora_job_cluster', {
    engine: engineType,
    engineVersion: '13.4',
    availabilityZones: zones,
    vpcSecurityGroupIds: [sg.id],
    dbSubnetGroupName: subnetGroup.name,
    storageEncrypted: true,
    clusterIdentifier: getResourceName('indexer-job'),
    dbClusterParameterGroupName: clusterParameterGroup.name,

    masterUsername: 'app',
    masterPassword: process.env.DB_PASSWORD,
    databaseName: 'app',

    skipFinalSnapshot: true,
    backupRetentionPeriod: isProduction() ? 7 : 1,
    preferredBackupWindow: '07:00-09:00',
  })

  const dbParameterGroup = new aws.rds.ParameterGroup('aurora_job_instance_param_group', {
    name: getResourceName('indexer-job-instance'),
    family: paramFamily,
    parameters: [
      {
        name: 'log_rotation_age',
        value: '1440',
      },
      {
        name: 'log_rotation_size',
        value: '102400',
      },
    ],
  })
  const instance = config.require('auroraJobInstance')
  const numInstances = parseInt(config.require('auroraJobInstances')) || 1
  const clusterInstances: aws.rds.ClusterInstance[] = []
  for (let i = 0; i < numInstances; i++) {
    clusterInstances.push(new aws.rds.ClusterInstance(`aurora_job_instance_${i + 1}`, {
      identifier: getResourceName(`indexer-job-${i+1}`),
      clusterIdentifier: cluster.id,
      instanceClass: instance,
      engine: engineType,
      engineVersion: cluster.engineVersion,
      dbParameterGroupName: dbParameterGroup.name,
      dbSubnetGroupName: subnetGroup.name,
      availabilityZone: zones[0],
      autoMinorVersionUpgrade: true,
      publiclyAccessible: isFalse(isProduction()),
    }))
  }

  return cluster
}

const createMainEventDb = (
  config: pulumi.Config,
  vpc: ec2.Vpc,
  sg: aws.ec2.SecurityGroup,
  zones: string[],
): aws.rds.Cluster => {
  const paramFamily = 'aurora-postgresql13'
  const clusterParameterGroup = new aws.rds.ClusterParameterGroup('aurora_event_cluster_param_group', {
    name: getResourceName('indexer-event-cluster'),
    family: paramFamily,
    parameters: [
      {
        name: 'rds.force_ssl',
        value: '1',
        applyMethod: 'pending-reboot',
      },
    ],
  })

  const subnetGroup = getSubnetGroupEvent(vpc)
  const engineType = EngineType.AuroraPostgresql
  const cluster = new aws.rds.Cluster('aurora_event_cluster', {
    engine: engineType,
    engineVersion: '13.4',
    availabilityZones: zones,
    vpcSecurityGroupIds: [sg.id],
    dbSubnetGroupName: subnetGroup.name,
    storageEncrypted: true,
    clusterIdentifier: getResourceName('indexer-event'),
    dbClusterParameterGroupName: clusterParameterGroup.name,

    masterUsername: 'app',
    masterPassword: process.env.DB_PASSWORD,
    databaseName: 'app',

    skipFinalSnapshot: true,
    backupRetentionPeriod: isProduction() ? 7 : 1,
    preferredBackupWindow: '07:00-09:00',
  })

  const dbParameterGroup = new aws.rds.ParameterGroup('aurora_event_instance_param_group', {
    name: getResourceName('indexer-event-instance'),
    family: paramFamily,
    parameters: [
      {
        name: 'log_rotation_age',
        value: '1440',
      },
      {
        name: 'log_rotation_size',
        value: '102400',
      },
    ],
  })
  const instance = config.require('auroraEventInstance')
  const numInstances = parseInt(config.require('auroraEventInstances')) || 1
  const clusterInstances: aws.rds.ClusterInstance[] = []
  for (let i = 0; i < numInstances; i++) {
    clusterInstances.push(new aws.rds.ClusterInstance(`aurora_event_instance_${i + 1}`, {
      identifier: getResourceName(`indexer-event-${i+1}`),
      clusterIdentifier: cluster.id,
      instanceClass: instance,
      engine: engineType,
      engineVersion: cluster.engineVersion,
      dbParameterGroupName: dbParameterGroup.name,
      dbSubnetGroupName: subnetGroup.name,
      availabilityZone: zones[0],
      autoMinorVersionUpgrade: true,
      publiclyAccessible: isFalse(isProduction()),
    }))
  }

  return cluster
}

const createMainGraphDb = (
  config: pulumi.Config,
  vpc: ec2.Vpc,
  sg: aws.ec2.SecurityGroup,
  zones: string[],
): aws.rds.Cluster => {
  const paramFamily = 'aurora-postgresql13'
  const clusterParameterGroup = new aws.rds.ClusterParameterGroup('aurora_graph_cluster_param_group', {
    name: getResourceName('indexer-graph-cluster'),
    family: paramFamily,
    parameters: [
      {
        name: 'rds.force_ssl',
        value: '1',
        applyMethod: 'pending-reboot',
      },
    ],
  })

  const subnetGroup = getSubnetGroupGraph(vpc)
  const engineType = EngineType.AuroraPostgresql
  const cluster = new aws.rds.Cluster('aurora_graph_cluster', {
    engine: engineType,
    engineVersion: '13.4',
    availabilityZones: zones,
    vpcSecurityGroupIds: [sg.id],
    dbSubnetGroupName: subnetGroup.name,
    storageEncrypted: true,
    clusterIdentifier: getResourceName('indexer-graph'),
    dbClusterParameterGroupName: clusterParameterGroup.name,

    masterUsername: 'app',
    masterPassword: process.env.DB_PASSWORD,
    databaseName: 'app',

    skipFinalSnapshot: true,
    backupRetentionPeriod: isProduction() ? 7 : 1,
    preferredBackupWindow: '07:00-09:00',
  })

  const dbParameterGroup = new aws.rds.ParameterGroup('aurora_graph_instance_param_group', {
    name: getResourceName('indexer-graph-instance'),
    family: paramFamily,
    parameters: [
      {
        name: 'log_rotation_age',
        value: '1440',
      },
      {
        name: 'log_rotation_size',
        value: '102400',
      },
    ],
  })
  const instance = config.require('auroraGraphInstance')
  const numInstances = parseInt(config.require('auroraGraphInstances')) || 1
  const clusterInstances: aws.rds.ClusterInstance[] = []
  for (let i = 0; i < numInstances; i++) {
    clusterInstances.push(new aws.rds.ClusterInstance(`aurora_graph_instance_${i + 1}`, {
      identifier: getResourceName(`indexer-graph-${i+1}`),
      clusterIdentifier: cluster.id,
      instanceClass: instance,
      engine: engineType,
      engineVersion: cluster.engineVersion,
      dbParameterGroupName: dbParameterGroup.name,
      dbSubnetGroupName: subnetGroup.name,
      availabilityZone: zones[0],
      autoMinorVersionUpgrade: true,
      publiclyAccessible: isFalse(isProduction()),
    }))
  }

  return cluster
}

// spin up 3x databases for indexer (jobs, events, graph). All same configs 
export const createAuroraClustersJob = (
  config: pulumi.Config,
  vpc: ec2.Vpc,
  sg: aws.ec2.SecurityGroup,
  zones: string[],
): AuroraOutput => {
  const main = createMainJobDb(config, vpc, sg, zones)
  return { main }
}

export const createAuroraClustersEvent = (
  config: pulumi.Config,
  vpc: ec2.Vpc,
  sg: aws.ec2.SecurityGroup,
  zones: string[],
): AuroraOutput => {
  const main = createMainEventDb(config, vpc, sg, zones)
  return { main }
}

export const createAuroraClustersGraph = (
  config: pulumi.Config,
  vpc: ec2.Vpc,
  sg: aws.ec2.SecurityGroup,
  zones: string[],
): AuroraOutput => {
  const main = createMainGraphDb(config, vpc, sg, zones)
  return { main }
}