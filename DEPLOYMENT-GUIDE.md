# Deployment Guide

1. [Requirements](#requirements)
2. [Building the containers](#building-the-containers)
3. [Running the containers](#running-the-containers)
    1. [Job API](#job-api)
    2. [Job Watcher](#job-watcher)
    3. [Parsing Dispatcher](#parsing-dispatcher)
    4. [Chain Watcher](#chain-watcher)

## Requirements

* [Docker](https://docs.docker.com/get-docker/)
* [Postgres](https://hub.docker.com/_/postgres)
* [Redis](https://hub.docker.com/_/redis)

Natively installed postgres and redis instances can alternatively be used instead of running them containers.

## Building the Containers

In order to run the indexer the first step is to build the container images.

For this the command below allows building and tagging the containers. Replace `<name>` with:

* api
* jobwatcher
* parsingdispatcher
* watcher

```console
docker build . -f Dockerfile.<name> -t indexer:<name>
```

## Running the Containers

### Job API

Job API allows creating, listing, and updating discovery and parsing jobs.
See the [job API binary readme file](cmd/jobs-api/README.md) for more details about its flags.

#### Requirements

* Postgres

#### Starting the Container

```console
docker run indexer:api -d "host=<postgres_host> port=<postgres_port> user=<postgres_user> password=<postgres_password> dbname=jobs sslmode=<postgres_sslmode>"
```

### Job Watcher

Job Watcher watches the dispatcher and parsing websockets for new updates and pushes them into their respective queue.
See the [job watcher binary readme file](cmd/jobs-watcher/README.md) for more details about its flags.

#### Requirements

* Jobs API
* Redis

#### Starting the container

```console
docker run indexer:jobwatcher -a <jobs_api_url> -u <redis_url>
```

### Parsing Dispatcher

The Parsing Dispatcher consumes messages from the queue and launches jobs.
See the [parsing dispatcher binary readme file](cmd/parsing-dispatcher/README.md) for more details about its flags.

#### Requirements

* Postgres
* Jobs API
* Redis
* AWS Credentials in Environment

#### Starting the Container

```console
docker run indexer:parsingdispatcher -u <redis_url> -a <jobs_api_url> -d "port=<postgres_port> user=<postgres_user> password=<postgres_password> dbname=chains sslmode=<postgres_sslmode>"
```

### Chain Watcher

Chain Watcher watches the chain and instantiates all the parsing jobs required for the network.
If the chain watcher stopped during an instantiation, upon restarting it retrieves the last job saved in the API and starts from that height instead of 0.
See the [chain watcher binary readme file](cmd/chain-watcher/README.md) for more details about its flags.

#### Requirements

* Ethereum Node
* Jobs API

#### Starting the Container

```console
docker run indexer:chainwatcher -a <api_url> -u <web3_node_url> -i <web3_chain_id> -t web3 -c <contract> -e <event_type> --standard-type <standard_type>"
```

Here is an example where the watcher is configured to watch for:

* The following contract: [Fighter (FIGHTER)](https://etherscan.io/address/0x87E738a3d5E5345d6212D8982205A564289e6324) (`0x87E738a3d5E5345d6212D8982205A564289e6324`);
* With the event type _Transfer_ (`0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef`);
* With the `ERC721` standard type.

```console
docker run indexer:chainwatcher -a api:8081 -u wss://mainnet.infura.io/ws/v3/d7b15235a515483490a5b89644221a71 -i 1 -t web3 -c 0x87E738a3d5E5345d6212D8982205A564289e6324 -e 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef --standard-type ERC721"
```

> 🚧 The Chain Watcher will in the future no longer need an event type, contract and standard type.

(Related issue: https://github.com/NFT-com/indexer/issues/47)
