# Indexer Service

## Local Development

### Requirements

* [Docker](https://docs.docker.com/get-docker/)
* [Docker Compose](https://docs.docker.com/compose/install/)
* [AWS SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)

### Testing Locally

In order to run a local test of the indexer, access to an Ethereum node is required.
The output will show in the logs of the lambdas function sam CLI command.

#### Starting the Lambdas

```bash
sam local start-lambdas
```