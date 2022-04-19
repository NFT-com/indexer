# Deployment Guide

1. [Requirements](#requirements)
2. [Building the containers](#building-the-containers)
3. [Running the containers](#running-the-containers)
    1. [Job API](#job-api)
    2. [Job Watcher](#job-watcher)
    3. [Parsing Dispatcher](#parsing-dispatcher)
    3. [Addition Dispatcher](#addition-dispatcher)
    4. [Chain Watcher](#chain-watcher)
    5. [Functions](#functions)

## Requirements

* [Docker](https://docs.docker.com/get-docker/)
* [Postgres](https://hub.docker.com/_/postgres)
* [Redis](https://hub.docker.com/_/redis)

Natively installed postgres and redis instances can alternatively be used instead of running them in containers.

## Building the Containers

In order to run the indexer the first step is to build the container images.

For this the command below allows building and tagging the containers. Replace `<name>` with:

* api
* jobwatcher
* parsingdispatcher
* additiondispatcher
* chainwatcher

```console
docker build . -f Dockerfile.<name> -t indexer-<name>:1.0.0
```

## Running the Containers

### Job API

Job API allows creating, listing, and updating discovery and parsing jobs.
See the [job API binary readme file](cmd/jobs-api/README.md) for more details about its flags.

#### Requirements

* Postgres

#### Starting the Container

```console
docker run -p '8081:8081' indexer-api:1.0.0 -d "host=<postgres_host> port=<postgres_port> user=<postgres_user> password=<postgres_password> dbname=jobs sslmode=<postgres_sslmode>"
```

> ⚠️ If you use a local postgres instance you need to add the containers into the database network.
> If you use a global system instance, set the network as host `--network=host`.
> Otherwise, set the network as the same network as the database instance.

### Job Watcher

Job Watcher watches the dispatcher and parsing websockets for new updates and pushes them into their respective queue.
See the [job watcher binary readme file](cmd/jobs-watcher/README.md) for more details about its flags.

#### Requirements

* Jobs API
* Redis

#### Starting the Container

```console
docker run indexer-jobwatcher:1.0.0 -a <jobs_api_url> -u <redis_url>
```

### Parsing Dispatcher

The Parsing Dispatcher consumes messages from the queue and launches jobs.
See the [parsing dispatcher binary readme file](cmd/parsing-dispatcher/README.md) for more details about its flags.

#### Requirements

* Postgres
* Jobs API
* Redis
* [AWS Credentials in Environment](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html)
* [Deployed Functions to AWS](#functions)

#### Starting the Container

```console
docker run -e AWS_REGION='<aws_region>' -e AWS_ACCESS_KEY_ID='<aws_key_id>' -e AWS_SECRET_ACCESS_KEY='<aws_access_key>' indexer-parsingdispatcher:1.0.0 -u <redis_url> -a <jobs_api_url> -d "port=<postgres_port> user=<postgres_user> password=<postgres_password> dbname=chains sslmode=<postgres_sslmode>"
```

### Addition Dispatcher

Addition Dispatcher consumes messages from the queue and launches jobs.
See the [parsing dispatcher binary readme file](cmd/parsing-dispatcher/README.md) for more details about its flags.

#### Requirements

* Postgres
* Jobs API
* Redis
* [AWS Credentials in Environment](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html)
* [Deployed Functions in AWS](#functions)

#### Starting the Container

```console
docker run -e AWS_REGION='<aws_region>' -e AWS_ACCESS_KEY_ID='<aws_key_id>' -e AWS_SECRET_ACCESS_KEY='<aws_access_key>' indexer-additiondispatcher:1.0.0 -u <redis_url> -a <jobs_api_url> -d "port=<postgres_port> user=<postgres_user> password=<postgres_password> dbname=chains sslmode=<postgres_sslmode>"
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
docker run indexer-chainwatcher:1.0.0 -a <api_url> -u <web3_node_url> -i <web3_chain_id> -t web3 -c <contract> -e <event_type> --standard-type <standard_type>"
```

Here is an example where the watcher is configured to watch for:

* The following contract: [Fighter (FIGHTER)](https://etherscan.io/address/0x87E738a3d5E5345d6212D8982205A564289e6324) (`0x87E738a3d5E5345d6212D8982205A564289e6324`);
* With the event type _Transfer_ (`0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef`);
* With the `ERC721` standard type.

```console
docker run indexer-chainwatcher:1.0.0 -a api:8081 -u wss://mainnet.infura.io/ws/v3/d7b15235a515483490a5b89644221a71 -i 1 -t web3 -c 0x87E738a3d5E5345d6212D8982205A564289e6324 -e 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef --standard-type ERC721
```

> 🚧
> The Chain Watcher will no longer need an event type, contract and standard type.

### Functions

> 🚧
> Right now there is no easy mode to deploy this to run locally.
> Currently use the pipeline package in the [pipeline branch](https://github.com/NFT-com/indexer/tree/pipeline) to deploy them to AWS.
>
> After cloning the pipeline branch.
> Go into the pipeline folder and run:
> 
> * ` GOOS=linux GOARCH=amd64 go build -o worker ../parsers/web3/opensea/ordersmatched `
> * ` zip opensea_ordersmatched.zip worker  `
> * ` GOOS=linux GOARCH=amd64 go build -o worker ../parsers/web3/erc721/transfer `
> * ` zip erc721_transfer.zip worker `
> * ` GOOS=linux GOARCH=amd64 go build -o worker ../cmd/addition-worker `
> * ` zip addition.zip worker `
>
> After this, with the functions already zipped.
> There is two options for deployment:
> 
> * Using [Pulumi](https://www.pulumi.com/)
> * [Manually](https://docs.aws.amazon.com/lambda/latest/dg/gettingstarted-package.html)
> 