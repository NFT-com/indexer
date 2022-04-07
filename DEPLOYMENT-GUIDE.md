# Deployment Guide

1. [Requirements](#requirements)
2. [Build the containers](#build-the-containers)
3. [Running the containers](#running-the-containers)
    1. [Job API](#job-api)
    2. [Job Watcher](#job-watcher)
    3. [Parsing Dispatcher](#parsing-dispatcher)
    4. [Chain Watcher](#chain-watcher)

## Requirements

* [Docker](https://docs.docker.com/get-docker/)
* [Postgres](https://hub.docker.com/_/postgres)
* [Redis](https://hub.docker.com/_/redis)

Native installed postgres and redis will also work. No need to have postgres and redis in containers.

## Build the containers

In order to run the indexer the first step is to build the container images.

For this the command bellow allows building and tagging the containers. Replace `<name>` with:

* api
* jobwatcher
* parsingdispatcher
* watcher

```console
docker build . -f Dockerfile.<name> -t indexer:<name>
```

## Running the containers

### Job API

Job API allows creating, listing, and updating discovery and parsing jobs. Flags with
descriptions [here.](cmd/jobs-api/README.md)

#### Requirements

* Postgres

#### Starting the container

```console
docker run indexer:api -d "host=<postgres_host> port=<postgres_port> user=<postgres_user> password=<postgres_password> dbname=jobs sslmode=<postgres_sslmode>"
```

### Job Watcher

Job Watcher watches the dispatcher and parsing websockets for new updates and pushes them into their respective queue.
Flags with descriptions [here.](cmd/jobs-watcher/README.md)

#### Requirements

* Jobs API
* Redis

#### Starting the container

```console
docker run indexer:jobwatcher -a <jobs_api_url> -u <redis_url>
```

### Parsing Dispatcher

Parsing Dispatcher consumes messages from the queue and launches lambdas. Flags with
descriptions [here.](cmd/parsing-dispatcher/README.md)

#### Requirements

* Postgres
* Jobs API
* Redis
* AWS Credentials in Environment

#### Starting the container

```console
docker run indexer:parsingdispatcher -u <redis_url> -a <jobs_api_url> -d "port=<postgres_port> user=<postgres_user> password=<postgres_password> dbname=chains sslmode=<postgres_sslmode>"
```

### Chain Watcher

Chain Watcher watches the chain and instantiates all the parsing jobs required for the network. If the chain watcher
stop in the middle of the instantiation, it will retrieve the last job saved in the API and start from that height
instead of 0. Flags with descriptions [here.](cmd/chain-watcher/README.md)

#### Requirements

* Ethereum Node
* Jobs API

#### Starting the container

```console
docker run indexer:parsingdispatcher -a <api_url> -u <web3_node_url> -i <web3_chain_id> -t web3 -c <contract> -e <event_type> --standard-type <standard_type>"
```

For example, running the watcher for
contract [Fighter (FIGHTER)](https://etherscan.io/address/0x87E738a3d5E5345d6212D8982205A564289e6324) (`0x87E738a3d5E5345d6212D8982205A564289e6324`)
, the event type Transfer (`0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef`) and standard
type `ERC721`.

```console
docker run indexer:parsingdispatcher -a api:8081 -u wss://mainnet.infura.io/ws/v3/d7b15235a515483490a5b89644221a71 -i 1 -t web3 -c 0x87E738a3d5E5345d6212D8982205A564289e6324 -e 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef --standard-type ERC721"
```

> 🚧 The Chain Watcher will change to remove the need to pass event_type, contract and standard_type. As it is not yet implemented. I will update this in the future.

(Related issue: https://github.com/NFT-com/indexer/issues/47)
** TODO CHAIN-WATCHER **
