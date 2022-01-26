# Indexer Service

## Local Development

### Building the indexer

```bash
go build -o indexer ./cmd/indexer/main.go
```

### Testing locally

In order to test locally the indexer its required to have access to an ethereum node.

### Starting the indexer

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