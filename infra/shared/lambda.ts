import * as pulumi from "@pulumi/pulumi";
import * as aws from "@pulumi/aws";

const ZMOK_HTTP_URL = 'https://api.zmok.io/custom1/qkher8p6hmchaxni'

export const createParsingWorker = (): aws.lambda.Function => {
    return new aws.lambda.Function("parsing-worker", {
        architectures: ["x86_64"],
        description: "parsing-web3",
        handler: "worker",
        memorySize: 128,
        code: new pulumi.asset.FileArchive("./parsing.zip"),
        name: "parsing-worker",
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
    return new aws.lambda.Function("addition-worker", {
        architectures: ["x86_64"],
        description: "addition-web3",
        handler: "worker",
        memorySize: 128,
        code: new pulumi.asset.FileArchive("./addition.zip"),
        name: "addition-worker",
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
    return new aws.lambda.Function("completion-worker", {
        architectures: ["x86_64"],
        description: "completion-web3",
        handler: "worker",
        memorySize: 128,
        code: new pulumi.asset.FileArchive("./completion.zip"),
        name: "completion-worker",
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
