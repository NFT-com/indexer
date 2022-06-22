import * as pulumi from '@pulumi/pulumi';
import * as aws from '@pulumi/aws';

import { getResourceName } from '../helper'
import { SharedInfraOutput } from '../defs'

const event_db = `host=${process.env.DB_EVENT_HOST} port=${process.env.DB_PORT} user=${process.env.DB_USER} password=${process.env.DB_PASSWORD} dbname=${process.env.DB_NAME}`
const job_db = `host=${process.env.DB_JOB_HOST} port=${process.env.DB_PORT} user=${process.env.DB_USER} password=${process.env.DB_PASSWORD} dbname=${process.env.DB_NAME}`
const graph_db = `host=${process.env.DB_GRAPH_HOST} port=${process.env.DB_PORT} user=${process.env.DB_USER} password=${process.env.DB_PASSWORD} dbname=${process.env.DB_NAME}`
const execRole = 'arn:aws:iam::016437323894:role/ecsTaskExecutionRole'
const taskRole = 'arn:aws:iam::016437323894:role/ECSServiceTask'
const ec2UserData =
`#!/bin/bash
echo ECS_CLUSTER=dev-indexer >> /etc/ecs/ecs.config
echo ECS_BACKEND_HOST= >> /etc/ecs/ecs.config`

export const createNsqlookupTaskDefinition = (): aws.ecs.TaskDefinition => {
    return new aws.ecs.TaskDefinition('nsqlookupd', 
    {
        containerDefinitions: JSON.stringify([
            {
                cpu: 0,
                entryPoint: '/nsqlookupd',
                environment: [],
                essential: true,
                image: 'nsqio/nsq',
                links: [],
                memoryReservation: 256,
                mountPoints: [],
                name: 'nsqlookupd',
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
        family: 'nsqlookupd',
        requiresCompatibilities: ['EC2'],
        taskRoleArn: taskRole,
    })
}

export const createNsqdTaskDefinition = (): aws.ecs.TaskDefinition => {
    return new aws.ecs.TaskDefinition('nsqd', 
    {
        containerDefinitions: JSON.stringify([
            {
                command: ['--lookupd-tcp-address=3.88.22.8:4160','--broadcast-address=3.88.22.8'],
                cpu: 0,
                entryPoint: '/nsqd',
                environment: [],
                essential: true,
                image: 'nsqio/nsq',
                links: [],
                memoryReservation: 256,
                mountPoints: [],
                name: 'nsq',
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
        family: 'nsqd',
        requiresCompatibilities: ['EC2'],
        taskRoleArn: taskRole,
    })
}

export const createParsingDispatcherTaskDefinition = (
    infraOutput: SharedInfraOutput,
): aws.ecs.TaskDefinition => {
    const resourceName = getResourceName('indexer-parsing-dispatcher')
    const ecrImage = `${process.env.ECR_REGISTRY}/${infraOutput.indexerECRRepo}:parsing-dispatcher`
    
    return new aws.ecs.TaskDefinition('indexer-parsing-dispatcher', 
    {
        containerDefinitions: JSON.stringify([
            {
                command: ['-n','parsing-worker','-k',`${process.env.EC2_PUBLIC_IP}:4161`,'-q',`${process.env.EC2_PUBLIC_IP}:4150`,'--height-range',process.env.PARSER_HEIGHT_RANGE,'--rate-limit',process.env.PARSER_RATE_LIMIT,'-j',job_db,'-e',event_db,'-g',graph_db],
                cpu: 0,
                entryPoint: '/dispatcher',
                environment: [],
                essential: true,
                image: ecrImage,
                links: [],
                memoryReservation: 2048,
                mountPoints: [],
                name: resourceName,
                portMappings: [],
                secrets: [
                    {
                        name: 'AWS_ACCESS_KEY_ID',
                        value: process.env.AWS_ACCESS_KEY_ID
                    },
                    {
                        name: 'AWS_REGION',
                        value: process.env.AWS_REGION
                    },
                    {
                        name: 'AWS_SECRET_ACCESS_KEY',
                        value: process.env.AWS_SECRET_ACCESS_KEY
                    },
                ],
                volumesFrom: []
        }]),
        executionRoleArn: execRole,
        family: resourceName,
        cpu: '512',
        memory: '2048',
        requiresCompatibilities: ['EC2'],
        taskRoleArn: taskRole,
    })
}

export const createAdditionDispatcherTaskDefinition = (
    infraOutput: SharedInfraOutput,
): aws.ecs.TaskDefinition => {
    const execRole = 'arn:aws:iam::016437323894:role/ecsTaskExecutionRole'
    const taskRole = 'arn:aws:iam::016437323894:role/ECSServiceTask'
    const resourceName = getResourceName('indexer-addition-dispatcher')
    const ecrImage = `${process.env.ECR_REGISTRY}/${infraOutput.indexerECRRepo}:addition-dispatcher`

    return new aws.ecs.TaskDefinition('indexer-addition-dispatcher', 
    {
        containerDefinitions: JSON.stringify([
            {
                command: ['-n','addition-worker','-k',`${process.env.EC2_PUBLIC_IP}:4161`,'--rate-limit',process.env.ACTION_RATE_LIMIT,'-g',graph_db,'-j',job_db],
                cpu: 0,
                entryPoint:'/dispatcher',
                environment:[],
                essential:true,
                image: ecrImage,
                links:[],
                memoryReservation: 2048,
                mountPoints:[],
                name: resourceName,
                portMappings:[],
                secrets: [
                    {
                        name: 'AWS_ACCESS_KEY_ID',
                        value: process.env.AWS_ACCESS_KEY_ID
                    },
                    {
                        name: 'AWS_REGION',
                        value: process.env.AWS_REGION
                    },
                    {
                        name: 'AWS_SECRET_ACCESS_KEY',
                        value: process.env.AWS_SECRET_ACCESS_KEY
                    },
                ],
                volumesFrom:[]
        }]),
        executionRoleArn: execRole,
        family: resourceName,
        cpu: '512',
        memory: '2048',
        requiresCompatibilities: ['EC2'],
        taskRoleArn: taskRole,
    })
}

        

export const createJobCreatorTaskDefinition = (
    infraOutput: SharedInfraOutput,
): aws.ecs.TaskDefinition => {
    const execRole = 'arn:aws:iam::016437323894:role/ecsTaskExecutionRole'
    const taskRole = 'arn:aws:iam::016437323894:role/ECSServiceTask'
    const resourceName = getResourceName('indexer-jobs-creator')
    const ecrImage = `${process.env.ECR_REGISTRY}/${infraOutput.indexerECRRepo}:jobs-creator`

    return new aws.ecs.TaskDefinition('indexer-jobs-creator', 
    {
        containerDefinitions: JSON.stringify([
            {
                command: ['-q',`${process.env.EC2_PUBLIC_IP}:4150`,'-w',process.env.ZMOK_WS_URL,'-g',graph_db,'-j',job_db],
                cpu: 0,
                entryPoint:'/creator',
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
        cpu: '512',
        memory: '2048',
        requiresCompatibilities: ['EC2'],
        taskRoleArn: taskRole,
    })
}

export const createEc2SecurityGroup = (
    infraOutput: SharedInfraOutput,
): aws.ec2.SecurityGroup => {
return new aws.ec2.SecurityGroup("ec2_ci_sg", {
    description: "ECS Allowed Ports",
    egress: [{
        cidrBlocks: ["0.0.0.0/0"],
        fromPort: 0,
        protocol: "-1",
        toPort: 0,
    }],
    ingress: [
        {
            cidrBlocks: ["0.0.0.0/0"],
            fromPort: 4160,
            protocol: "tcp",
            toPort: 4160,
        },
        {
            cidrBlocks: ["0.0.0.0/0"],
            fromPort: 4151,
            protocol: "tcp",
            toPort: 4151,
        },
        {
            cidrBlocks: ["0.0.0.0/0"],
            fromPort: 4150,
            protocol: "tcp",
            toPort: 4150,
        },
        {
            cidrBlocks: ["0.0.0.0/0"],
            fromPort: 4161,
            protocol: "tcp",
            toPort: 4161,
        },
        {
            cidrBlocks: ["0.0.0.0/0"],
            fromPort: 22,
            protocol: "tcp",
            toPort: 22,
        },
    ],
    name: "EC2ContainerService-dev-indexer-EcsSecurityGroup-19WMAN457GTFK",
    vpcId: infraOutput.vpcId,
})
}

const ecsAsgLaunchConfig = (
    infraOutput: SharedInfraOutput,
): aws.ec2.LaunchConfiguration => {
    const { id: launchConfigSG } = createEc2SecurityGroup(infraOutput)
    return new aws.ec2.LaunchConfiguration("dev-indexer-ecs-launchconfig", {
        associatePublicIpAddress: true,
        iamInstanceProfile: "arn:aws:iam::016437323894:instance-profile/ecsInstanceRole",
        imageId: "ami-0f863d7367abe5d6f",  //latest amzn linux 2 ecs-optimized ami in us-east-1
        instanceType: "m6id.large",
        keyName: "indexer_dev_key",
        name: "EC2ContainerService-dev-indexer-EcsInstanceLc-LD7MujggcNTl",
        rootBlockDevice: {
            deleteOnTermination: false,
            volumeSize: 30,
            volumeType: "gp2",
        },
        securityGroups: [launchConfigSG],
        userData: ec2UserData,
})
}

const ecsClusterASG = (
    config: pulumi.Config,
    infraOutput: SharedInfraOutput,
): aws.autoscaling.Group => {
    return new aws.autoscaling.Group("dev-indexer-ecs-asg", {
        availabilityZones: [config.require('availabilityZones')],
        defaultCooldown: 300,
        desiredCapacity: 1,
        healthCheckGracePeriod: 0,
        healthCheckType: "EC2",
        launchConfiguration: ecsAsgLaunchConfig(infraOutput),
        maxSize: 1,
        minSize: 0,
        name: "EC2ContainerService-dev-indexer-EcsInstanceAsg-5NQSZP3BD1JV",
        serviceLinkedRoleArn: "arn:aws:iam::016437323894:role/aws-service-role/autoscaling.amazonaws.com/AWSServiceRoleForAutoScaling",
        tags: [
            {
                key: "Description",
                propagateAtLaunch: true,
                value: "This instance is the part of the Auto Scaling group which was created through ECS Console",
            },
            {
                key: "AmazonECSManaged",
                propagateAtLaunch: true,
                value: "",
            },
            {
                key: "Name",
                propagateAtLaunch: true,
                value: "ECS Instance - EC2ContainerService-dev-indexer",
            },
        ],
        vpcZoneIdentifiers: infraOutput.publicSubnets,
    })
}

const ecsCapacityProvider = (
    config: pulumi.Config,
    infraOutput: SharedInfraOutput,
): string => {
    const resourceName = getResourceName('indexer-capacity-provider')
    const { arn: arn_asg } = ecsClusterASG(config, infraOutput)
    const ecp =  new aws.ecs.CapacityProvider(resourceName, {
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
    return resourceName
}

export const createEcsCluster = (
    config: pulumi.Config,
    infraOutput: SharedInfraOutput,
): aws.ecs.Cluster => {
    const resourceName = getResourceName('indexer')
    const capacityProvider = ecsCapacityProvider(config, infraOutput)
    const cluster = new aws.ecs.Cluster(resourceName, 
    {
        name: getResourceName(resourceName),
        settings: [
            {
            name: 'containerInsights',
            value: 'enabled',
        }],
        capacityProviders: [capacityProvider]
    })

    new aws.ecs.ClusterCapacityProviders('indexer-cluster-capacity-provider', {
        clusterName: cluster.name,
        capacityProviders: [capacityProvider],
        defaultCapacityProviderStrategies: [
          {
            weight: 100,
            capacityProvider: capacityProvider,
          },
        ],
    })
    return cluster 
}