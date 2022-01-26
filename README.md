# Indexer Service

## Local Development

### Requirements

* [Docker](https://docs.docker.com/get-docker/)
* [AWS SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)

### Testing locally

In order to run a local test of the indexer, access to an Ethereum node is required.

#### Starting the Lambdas

```bash
sam local start-lambdas
```

### Building the Indexer

```bash
go build -o indexer ./cmd/indexer/main.go
```

### Testing locally

In order to run a local test of the indexer, access to an Ethereum node is required.

### Starting the Indexer

There are two ways to run the indexer.

#### Live Mode

```bash
indexer <node_url> <network> <chain>
```

Example for the Ethereum mainnet:

```bash
indexer <node_url> ethereum mainnet
```

#### Historical Mode

```bash
indexer <node_url> <network> <chain>
```

Example for the Ethereum mainnet for a range from block 1234 to block 8910:

```bash
indexer <node_url> ethereum mainnet -s 1234 -e 8910
```