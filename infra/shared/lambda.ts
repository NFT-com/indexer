import * as pulumi from "@pulumi/pulumi";
import * as aws from "@pulumi/aws";


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
                NODE_URL: "https://api.zmok.io/custom1/qkher8p6hmchaxni",
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
                NODE_URL: "https://api.zmok.io/custom1/qkher8p6hmchaxni",
            },
        },
        tracingConfig: {
            mode: "PassThrough",
        },
    })
}