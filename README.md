# Indexer Service

## Local Development

### Building the indexer

```bash
go build -o indexer ./cmd/indexer/main.go
```

### Testing locally

In order to test locally the indexer its required to have access to an ethereum node.

### Starting the indexer

```bash
indexer <node_url> <network> <chain>
```

Example for the Ethereum mainnet:

```bash
indexer <node_url> ethereum mainnet
```