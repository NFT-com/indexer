# Indexer

## Description

The `indexer` is in charge of listening for new events in configured blockchains and dispatch handlers to parse that data and insert it into databases.

## Usage

```bash
Usage: indexer [options] <node_url>

-s, --start int     height at which to start indexing
-e, --end int       height at which to stop indexing
```
