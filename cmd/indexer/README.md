# Indexer

## Description

The `indexer` is in charge of listening for new events in configured blockchains and dispatch handlers to parse that data and insert it into databases.

## Usage

```bash
Usage: indexer [options] <node_url> <network> <chain>

-s, --start int           height at which to start indexing
-e, --end int             height at which to stop indexing
-t, --test string         run indexer test local mode
-l, --lambda-url string   lambdas custom url if local
-r, --region string       aws region
```