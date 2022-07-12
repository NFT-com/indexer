import * as pulumi from '@pulumi/pulumi';
import * as aws from '@pulumi/aws';
import { getResourceName } from '../helper'
import { SharedInfraOutput } from '../defs'

const event_db = `host=${process.env.DB_EVENT_HOST} port=${process.env.DB_PORT} user=${process.env.DB_USER} password=${process.env.DB_PASSWORD} dbname=${process.env.DB_NAME}`
const job_db = `host=${process.env.DB_JOB_HOST} port=${process.env.DB_PORT} user=${process.env.DB_USER} password=${process.env.DB_PASSWORD} dbname=${process.env.DB_NAME}`
const graph_db = `host=${process.env.DB_GRAPH_HOST} port=${process.env.DB_PORT} user=${process.env.DB_USER} password=${process.env.DB_PASSWORD} dbname=${process.env.DB_NAME}`
const execRole = 'arn:aws:iam::016437323894:role/ecsTaskExecutionRole'
const taskRole = 'arn:aws:iam::016437323894:role/ECSServiceTask'

export const createNsqlookupTaskDefinition = (): aws.ecs.TaskDefinition => {
    const resourceName = 'nsqlookupd'
    return new aws.ecs.TaskDefinition(resourceName, 
    {
        containerDefinitions: JSON.stringify([
            {
                cpu: 0,
                entryPoint: ['/nsqlookupd'],
                environment: [],
                essential: true,
                image: 'nsqio/nsq',
                links: [],
                memoryReservation: 256,
                mountPoints: [],
                name: resourceName,
                portMappings: [
                    { 
                        containerPort: 4160,
                        hostPort: 4160,
                        protocol: 'tcp'
                    },
                    { 
                        containerPort: 4161,
                        hostPort: 4161,
                        protocol: 'tcp'
                    }
                ],
                volumesFrom: []
        }]),
        executionRoleArn: execRole,
        family: resourceName,
        requiresCompatibilities: ['EC2'],
        taskRoleArn: taskRole,
    })
}

export const createNsqdTaskDefinition = (): aws.ecs.TaskDefinition => {
    const resourceName = 'nsqd'
    return new aws.ecs.TaskDefinition(resourceName, 
    {
        containerDefinitions: JSON.stringify([
            {
                command: [`--lookupd-tcp-address=${process.env.EC2_PUBLIC_IP}:4160`,`--broadcast-address=${process.env.EC2_PUBLIC_IP}`,`--msg-timeout=15m`,`--max-msg-timeout=15m`],
                cpu: 0,
                entryPoint: ['/nsqd'],
                environment: [],
                essential: true,
                image: 'nsqio/nsq',
                links: [],
                memoryReservation: 256,
                mountPoints: [],
                name: resourceName,
                portMappings: [
                    { 
                        containerPort: 4150,
                        hostPort: 4150,
                        protocol: 'tcp'
                    },
                    { 
                        containerPort: 4151,
                        hostPort: 4151,
                        protocol: 'tcp'
                    }
                ],
                volumesFrom: []
        }]),
        executionRoleArn: execRole,
        family: resourceName,
        requiresCompatibilities: ['EC2'],
        taskRoleArn: taskRole,
    })
}

export const createNsqadminTaskDefinition = (): aws.ecs.TaskDefinition => {
    const resourceName = 'nsqadmin'
    return new aws.ecs.TaskDefinition(resourceName, 
    {
        containerDefinitions: JSON.stringify([
            {
                command: [`--lookupd-http-address=${process.env.EC2_PUBLIC_IP}:4161`],
                cpu: 0,
                entryPoint: ['/nsqadmin'],
                environment: [],
                essential: true,
                image: 'nsqio/nsq',
                links: [],
                memoryReservation: 256,
                mountPoints: [],
                name: resourceName,
                portMappings: [
                    { 
                        containerPort: 4171,
                        hostPort: 4171,
                        protocol: 'tcp'
                    }
                ],
                volumesFrom: []
        }]),
        executionRoleArn: execRole,
        family: resourceName,
        requiresCompatibilities: ['EC2'],
        taskRoleArn: taskRole,
    })
}

export const createParsingDispatcherTaskDefinition = (
    infraOutput: SharedInfraOutput,
): aws.ecs.TaskDefinition => {
    const resourceName = getResourceName('indexer-td-parsing-dispatcher')
    const ecrImage = `${process.env.ECR_REGISTRY}/${infraOutput.indexerECRRepo}:parsing-dispatcher`
    
    return new aws.ecs.TaskDefinition(resourceName, 
    {
        containerDefinitions: JSON.stringify([
            {
                command: ['-n',process.env.PARSING_LAMBDA_NAME,'-k',`${process.env.EC2_PUBLIC_IP}:4161`,'-q',`${process.env.EC2_PUBLIC_IP}:4150`,'--height-range',process.env.PARSER_HEIGHT_RANGE,'--rate-limit',process.env.PARSER_RATE_LIMIT,'-j',job_db,'-e',event_db,'-g',graph_db,'-l',process.env.INDEXER_LOG_LEVEL],
                cpu: 0,
                entryPoint: ['/dispatcher'],
                essential: true,
                image: ecrImage,
                links: [],
                memoryReservation: 2048,
                mountPoints: [],
                name: resourceName,
                portMappings: [],
                environment: [
                    {
                        Name: 'AWS_ACCESS_KEY_ID',
                        Value: process.env.AWS_ACCESS_KEY_ID
                    },
                    {
                        Name: 'AWS_REGION',
                        Value: process.env.AWS_REGION
                    },
                    {
                        Name: 'AWS_SECRET_ACCESS_KEY',
                        Value: process.env.AWS_SECRET_ACCESS_KEY
                    },
                ],
                volumesFrom: []
        }]),
        executionRoleArn: execRole,
        family: resourceName,
        cpu: '1024',
        memory: '2048',
        requiresCompatibilities: ['EC2'],
        taskRoleArn: taskRole,
    })
}

export const createAdditionDispatcherTaskDefinition = (
    infraOutput: SharedInfraOutput,
): aws.ecs.TaskDefinition => {
    const resourceName = getResourceName('indexer-td-addition-dispatcher')
    const ecrImage = `${process.env.ECR_REGISTRY}/${infraOutput.indexerECRRepo}:addition-dispatcher`

    return new aws.ecs.TaskDefinition(resourceName, 
    {
        containerDefinitions: JSON.stringify([
            {
                command: ['-n',process.env.ADDITION_LAMBDA_NAME,'-k',`${process.env.EC2_PUBLIC_IP}:4161`,'--rate-limit',process.env.ADDITION_RATE_LIMIT,'-g',graph_db,'-j',job_db,'-l',process.env.INDEXER_LOG_LEVEL],
                cpu: 0,
                entryPoint: ['/dispatcher'],
                essential: true,
                image: ecrImage,
                links: [],
                memoryReservation: 2048,
                mountPoints: [],
                name: resourceName,
                portMappings: [],
                environment: [
                    {
                        Name: 'AWS_ACCESS_KEY_ID',
                        Value: process.env.AWS_ACCESS_KEY_ID
                    },
                    {
                        Name: 'AWS_REGION',
                        Value: process.env.AWS_REGION
                    },
                    {
                        Name: 'AWS_SECRET_ACCESS_KEY',
                        Value: process.env.AWS_SECRET_ACCESS_KEY
                    },
                ],
                volumesFrom:[]
        }]),
        executionRoleArn: execRole,
        family: resourceName,
        cpu: '1024',
        memory: '2048',
        requiresCompatibilities: ['EC2'],
        taskRoleArn: taskRole,
    })
}

export const createCompletionDispatcherTaskDefinition = (
    infraOutput: SharedInfraOutput,
): aws.ecs.TaskDefinition => {
    const resourceName = getResourceName('indexer-td-completion-dispatcher')
    const ecrImage = `${process.env.ECR_REGISTRY}/${infraOutput.indexerECRRepo}:completion-dispatcher`

    return new aws.ecs.TaskDefinition(resourceName,
    {
        containerDefinitions: JSON.stringify([
            {
                command: ['-n',process.env.COMPLETION_LAMBDA_NAME,'-k',`${process.env.EC2_PUBLIC_IP}:4161`,'--rate-limit',process.env.COMPLETION_RATE_LIMIT,'-g',graph_db,'-e',event_db,'-j',job_db,'-l',process.env.INDEXER_LOG_LEVEL],
                cpu: 0,
                entryPoint: ['/dispatcher'],
                essential: true,
                image: ecrImage,
                links: [],
                memoryReservation: 2048,
                mountPoints: [],
                name: resourceName,
                portMappings: [],
                environment: [
                    {
                        Name: 'AWS_ACCESS_KEY_ID',
                        Value: process.env.AWS_ACCESS_KEY_ID
                    },
                    {
                        Name: 'AWS_REGION',
                        Value: process.env.AWS_REGION
                    },
                    {
                        Name: 'AWS_SECRET_ACCESS_KEY',
                        Value: process.env.AWS_SECRET_ACCESS_KEY
                    },
                ],
                volumesFrom:[]
        }]),
        executionRoleArn: execRole,
        family: resourceName,
        cpu: '1024',
        memory: '2048',
        requiresCompatibilities: ['EC2'],
        taskRoleArn: taskRole,
    })
}

export const createJobCreatorTaskDefinition = (
    infraOutput: SharedInfraOutput,
): aws.ecs.TaskDefinition => {
    const resourceName = getResourceName('indexer-td-jobs-creator')
    const ecrImage = `${process.env.ECR_REGISTRY}/${infraOutput.indexerECRRepo}:jobs-creator`

    return new aws.ecs.TaskDefinition(resourceName, 
    {
        containerDefinitions: JSON.stringify([
            {
                command: ['-q',`${process.env.EC2_PUBLIC_IP}:4150`,'-w',process.env.ZMOK_WS_URL,'-g',graph_db,'-j',job_db,'-l',process.env.INDEXER_LOG_LEVEL],
                cpu: 0,
                entryPoint: ['/creator'],
                environment: [],
                essential: true,
                image: ecrImage,
                links: [],
                memoryReservation: 2048,
                mountPoints: [],
                name: resourceName,
                portMappings: [],
                volumesFrom: []
        }]),
        executionRoleArn: execRole,
        family: resourceName,
        cpu: '1024',
        memory: '2048',
        requiresCompatibilities: ['EC2'],
        taskRoleArn: taskRole,
    })
}

const createEcsAsgLaunchConfig = (
    infraOutput: SharedInfraOutput,
): aws.ec2.LaunchConfiguration => {
    const launchConfigSG = infraOutput.webSGId
    const clusterName = getResourceName('indexer')
    const resourceName = getResourceName('indexer-asg-launchconfig')
    const ec2UserData =
    `#!/bin/bash
    echo ECS_CLUSTER=${clusterName} >> /etc/ecs/ecs.config
    echo ECS_BACKEND_HOST= >> /etc/ecs/ecs.config`

    return new aws.ec2.LaunchConfiguration(resourceName, {
        associatePublicIpAddress: true,
        iamInstanceProfile: 'arn:aws:iam::016437323894:instance-profile/ecsInstanceRole',
        imageId: 'ami-0f863d7367abe5d6f',  //latest amzn linux 2 ecs-optimized ami in us-east-1
        instanceType: 'm6i.xlarge',
        keyName: 'indexer_dev_key',
        name: resourceName,
        rootBlockDevice: {
            deleteOnTermination: false,
            volumeSize: 30,
            volumeType: 'gp2',
        },
        securityGroups: [launchConfigSG],
        userData: ec2UserData,
    })
}

const createEcsASG = (
    config: pulumi.Config,
    infraOutput: SharedInfraOutput,
): aws.autoscaling.Group => {
    const resourceName = getResourceName('indexer-ec2')
    return new aws.autoscaling.Group(resourceName, {
        defaultCooldown: 300,
        desiredCapacity: 1,
        healthCheckGracePeriod: 0,
        healthCheckType: 'EC2',
        launchConfiguration: createEcsAsgLaunchConfig(infraOutput),
        maxSize: 1,
        minSize: 1,
        name: resourceName,
        serviceLinkedRoleArn: 'arn:aws:iam::016437323894:role/aws-service-role/autoscaling.amazonaws.com/AWSServiceRoleForAutoScaling',
        tags: [
            {
                key: 'Description',
                propagateAtLaunch: true,
                value: 'This instance is the part of the Auto Scaling group which was created through ECS Console',
            },
            {
                key: 'AmazonECSManaged',
                propagateAtLaunch: true,
                value: '',
            },
            {
                key: 'Name',
                propagateAtLaunch: true,
                value: resourceName,
            },
        ],
        vpcZoneIdentifiers: infraOutput.publicSubnets,
    })
}

const createEcsCapacityProvider = (
    config: pulumi.Config,
    infraOutput: SharedInfraOutput,
): aws.ecs.CapacityProvider => {
    const resourceName = getResourceName('indexer-cp')
    const { arn: arn_asg } = createEcsASG(config, infraOutput)
    return new aws.ecs.CapacityProvider(resourceName, {
        autoScalingGroupProvider: {
            autoScalingGroupArn: arn_asg,
            managedScaling: {
                instanceWarmupPeriod: 300,
                maximumScalingStepSize: 1,
                minimumScalingStepSize: 1,
                status: 'DISABLED',
                targetCapacity: 100,
            },
            managedTerminationProtection: 'DISABLED',
        },
        name: resourceName,
    })
}

export const createEcsCluster = (
    config: pulumi.Config,
    infraOutput: SharedInfraOutput,
): aws.ecs.Cluster => {
    const resourceName = getResourceName('indexer')
    const { name: capacityProvider } = createEcsCapacityProvider(config, infraOutput)
    const cluster = new aws.ecs.Cluster(resourceName, 
    {
        name: resourceName,
        settings: [
            {
            name: 'containerInsights',
            value: 'enabled',
        }],
        capacityProviders: [capacityProvider]
    })

    return cluster 
}
