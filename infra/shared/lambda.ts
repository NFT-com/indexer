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
        sourceCodeHash: "v2kKP0fLXtIZtQ/BnP25qNNAJhpnAU7xUka46vTAqgM=",
        timeout: 600,
        tracingConfig: {
            mode: "PassThrough",
        },
    }, {
        protect: true,
    })
}

export const createActionWorker = (): aws.lambda.Function => {
    return new aws.lambda.Function("action-worker", {
        architectures: ["x86_64"],
        description: "action-web3",
        handler: "worker",
        memorySize: 128,
        code: new pulumi.asset.FileArchive("./action.zip"),
        name: "action-worker",
        reservedConcurrentExecutions: -1,
        role: "arn:aws:iam::016437323894:role/AWSLambdaBasicExecutionRole",
        runtime: "go1.x",
        sourceCodeHash: "QPPPtmO3I7S2T/z1KIeOxELq8aGy1BlWGTnXxchln8Q=",
        timeout: 600,
        tracingConfig: {
            mode: "PassThrough",
        },
    }, {
        protect: true,
    })
}