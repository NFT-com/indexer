import * as pulumi from "@pulumi/pulumi";
import * as aws from "@pulumi/aws";
import { getResourceName } from "../helper";

const ZMOK_HTTP_URL = 'https://api.zmok.io/custom1/qkher8p6hmchaxni'

export const createParsingWorker = (): aws.lambda.Function => {
    const resourceName = getResourceName('parsing-worker')
    return new aws.lambda.Function(resourceName, {
        architectures: ["x86_64"],
        description: "parsing-web3",
        handler: "worker",
        memorySize: 128,
        code: new pulumi.asset.FileArchive("./parsing.zip"),
        name: resourceName,
        reservedConcurrentExecutions: -1,
        role: "arn:aws:iam::016437323894:role/AWSLambdaBasicExecutionRole",
        runtime: "go1.x",
        timeout: 600,
        environment: {
            variables: {
                NODE_URL: ZMOK_HTTP_URL,
            },
        },
        tracingConfig: {
            mode: "PassThrough",
        },
    })
}

export const createAdditionWorker = (): aws.lambda.Function => {
    const resourceName = getResourceName('addition-worker')
    return new aws.lambda.Function(resourceName, {
        architectures: ["x86_64"],
        description: "addition-web3",
        handler: "worker",
        memorySize: 128,
        code: new pulumi.asset.FileArchive("./addition.zip"),
        name: resourceName,
        reservedConcurrentExecutions: -1,
        role: "arn:aws:iam::016437323894:role/AWSLambdaBasicExecutionRole",
        runtime: "go1.x",
        timeout: 600,
        environment: {
            variables: {
                NODE_URL: ZMOK_HTTP_URL,
            },
        },
        tracingConfig: {
            mode: "PassThrough",
        },
    })
}

export const createCompletionWorker = (): aws.lambda.Function => {
    const resourceName = getResourceName('completion-worker')
    return new aws.lambda.Function(resourceName, {
        architectures: ["x86_64"],
        description: "completion-web3",
        handler: "worker",
        memorySize: 128,
        code: new pulumi.asset.FileArchive("./completion.zip"),
        name: resourceName,
        reservedConcurrentExecutions: -1,
        role: "arn:aws:iam::016437323894:role/AWSLambdaBasicExecutionRole",
        runtime: "go1.x",
        timeout: 600,
        environment: {
            variables: {
                NODE_URL: ZMOK_HTTP_URL,
            },
        },
        tracingConfig: {
            mode: "PassThrough",
        },
    })
}
