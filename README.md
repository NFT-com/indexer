# Indexer

The Indexer repository contains a set of services which work together to watch events on blockchain networks, parse the related data and persist it in an index database:

- [Jobs Creator](./cmd/jobs-creator)
- [Parsing Dispatcher](./cmd/parsing-dispatcher)
- [Addition Dispatcher](./cmd/addition-dispatcher)

It uses AWS Lammda functions to do some of the heavy lifting, so the service can scale and perform rapidly, but also conserve resources when not busy:

- [Parsing Worker](./cmd/parsing-worker)
- [Addition Worker](./cmd/addition-worker)

This index can then be used by the [Analytics API](https://github.com/NFT-com/analytics) to expose NFT analytics data.

* [Deployment Guide](./docs/deployment.md)