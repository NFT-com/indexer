# Indexer Service

## Local Development

### Requirements

* [Docker](https://docs.docker.com/get-docker/)
* [AWS SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)

### Testing Locally

In order to run a local test of the indexer, access to an Ethereum node is required.
The output will show in the logs of the lambdas function sam CLI command.

#### Starting the Lambdas

```bash
sam local start-lambdas
```

### Building the Indexer

```bash
go build -o indexer ./cmd/indexer/main.go
```

### Starting the Indexer

There are two ways to run the indexer.

#### Live Mode

```bash
indexer <node_url> <network> <chain> -n local -l http://127.0.0.1:3001
```

Example for the Ethereum mainnet:

```bash
indexer <node_url> web3 mainnet -n local -l http://127.0.0.1:3001
```

#### Historical Mode

```bash
indexer <node_url> <network> <chain> -n local -l http://127.0.0.1:3001
```

Example for the Ethereum mainnet for a range from block 1234 to block 8910:

```bash
indexer <node_url> web3 mainnet -s 1234 -e 8910 -n local -l http://127.0.0.1:3001
```