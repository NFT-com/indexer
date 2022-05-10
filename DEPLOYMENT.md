# Deployment Guide

1. [Requirements](#requirements)
2. [Building the containers](#building-the-containers)
3. [Running the containers](#running-the-containers)
    1. [Jobs Creator](#jobs-creator)
    2. [Jobs Watcher](#jobs-watcher)
    3. [Parsing Dispatcher](#parsing-dispatcher)
    4. [Action Dispatcher](#action-dispatcher)
    5. [Functions](#functions)

## Requirements

In order to tun the indexer it requires these components:

* [Docker](https://docs.docker.com/get-docker/)
* [Postgres](#postgres)
* [Redis](#redis)

Natively installed postgres and redis instances can alternatively be used instead of running them in containers.

### Postgres

There are a couple ways to have postgres running:

* [Natively](https://www.postgresql.org/download/)
* [Docker](https://hub.docker.com/_/postgres)

Example of running postgres with automatic migration:

```console
docker run postgres -d --name postgres --network indexer -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=postgres -p '5432:5432' -v './sql/:/docker-entrypoint-initdb.d/'
```

> 🚧
> If you update the sql files and want to redeploy them.
> There are two options to update the container:
> * Manually
> * Shutting down the container and running `docker volume prume` and then starting up the container again

### Redis

There are a couple ways to have redis running:

* [Natively](https://redis.io/docs/getting-started/installation/)
* [Docker](https://hub.docker.com/_/redis)

Example of using redis with docker:

```console
docker run redis -d --name redis --network indexer -p '6379:6379'
```

## Building the Containers

In order to run the indexer the first step is to build the container images.

For this the command below allows building and tagging the all containers.

```console
for d in cmd/* ; do name=$(echo "$d" | cut -c 5-) ; docker build . -f cmd/"$name"/Dockerfile -t indexer-"$name":1.0.0 ; done
```

## Running the Containers

### Jobs Creator

Jobs Creator watches the chain and instantiates all the parsing jobs required for the network. If the job creator
stopped during an instantiation, upon restarting it retrieves the last job saved in the API and starts from that height
instead of 0. See the [chain watcher binary readme file](cmd/jobs-creator/README.md) for more details about its flags.

#### Requirements

* Ethereum Node
* Postgres

#### Starting the Container

```console
docker run indexer-jobs-creator:1.0.0 --network indexer  -u <web3_node_url> -g "host=<postgres_host> port=<postgres_port> user=<postgres_user> password=<postgres_password> dbname=postgres sslmode=<postgres_sslmode>" -j "host=<postgres_host> port=<postgres_port> user=<postgres_user> password=<postgres_password> dbname=postgres sslmode=<postgres_sslmode>"
```

Here is an example where the watcher is configured to watch for:

* The following contract: [Fighter (FIGHTER)](https://etherscan.io/address/0x87E738a3d5E5345d6212D8982205A564289e6324) (`0x87E738a3d5E5345d6212D8982205A564289e6324`)
* With the event type _Transfer_ (`0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef`)
* With the `ERC721` standard type

```console
docker run indexer-chainwatcher:1.0.0 --network indexer  -u wss://mainnet.infura.io/ws/v3/d7b15235a515483490a5b89644221a71 -i 1 -t web3 -c 0x87E738a3d5E5345d6212D8982205A564289e6324 -e 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef --standard-type ERC721 -d "host=<postgres_host> port=<postgres_port> user=<postgres_user> password=<postgres_password> dbname=postgres sslmode=<postgres_sslmode>"
```

### Jobs Watcher

Job Watcher watches the dispatcher and parsing websockets for new updates and pushes them into their respective queue.
See the [job watcher binary readme file](cmd/jobs-watcher/README.md) for more details about its flags.

#### Requirements

* Postgres
* Redis

#### Starting the Container

```console
docker run indexer-jobs-creator:1.0.0 --network indexer  -u <redis_url> -j "host=<postgres_host> port=<postgres_port> user=<postgres_user> password=<postgres_password> dbname=postgres sslmode=<postgres_sslmode>"
```

### Parsing Dispatcher

The Parsing Dispatcher consumes messages from the queue and launches jobs.
See the [parsing dispatcher binary readme file](cmd/parsing-dispatcher/README.md) for more details about its flags.

#### Requirements

* Postgres
* Redis
* [AWS Credentials in Environment](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html)
* [Deployed Functions to AWS](#functions)

#### Starting the Container

```console
docker run -e AWS_REGION='<aws_region>' --network indexer  -e AWS_ACCESS_KEY_ID='<aws_key_id>' -e AWS_SECRET_ACCESS_KEY='<aws_access_key>' indexer-parsing-dispatcher:1.0.0 -u <redis_url> -j "port=<postgres_port> user=<postgres_user> password=<postgres_password> dbname=postgres sslmode=<postgres_sslmode>" -e "host=<postgres_host> port=<postgres_port> user=<postgres_user> password=<postgres_password> dbname=postgres sslmode=<postgres_sslmode>"
```

### Action Dispatcher

Action Dispatcher consumes messages from the queue and launches jobs.
See the [parsing dispatcher binary readme file](cmd/parsing-dispatcher/README.md) for more details about its flags.

#### Requirements

* Postgres
* Redis
* [AWS Credentials in Environment](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html)
* [Deployed Functions in AWS](#functions)

#### Starting the Container

```console
docker run -e AWS_REGION='<aws_region>' --network indexer  -e AWS_ACCESS_KEY_ID='<aws_key_id>' -e AWS_SECRET_ACCESS_KEY='<aws_access_key>' indexer-action-dispatcher:1.0.0 -u <redis_url> -g "port=<postgres_port> user=<postgres_user> password=<postgres_password> dbname=postgres sslmode=<postgres_sslmode>" -j "host=<postgres_host> port=<postgres_port> user=<postgres_user> password=<postgres_password> dbname=postgres sslmode=<postgres_sslmode>"
```

### Functions

> 🚧
> Right now there is no easy mode to deploy this to run locally.
> Currently use the pipeline package in the [pipeline branch](https://github.com/NFT-com/indexer/tree/pipeline) to deploy them to AWS.
> Note that this branch could have not been rebased on master, or the branch you want to test with.
> Before deploying the lambdas make sure that the branch is updated.
>
> After cloning the pipeline branch.
> Go into the pipeline folder and run:
>
> * ` GOOS=linux GOARCH=amd64 go build -o worker ../cmd/parsing-worker `
> * ` zip parsing.zip worker `
> * ` GOOS=linux GOARCH=amd64 go build -o worker ../cmd/action-worker `
> * ` zip action.zip worker `
>
> After this, with the functions already zipped.
> There is two options for deployment:
>
> * Using [Pulumi](https://www.pulumi.com/)
> * [Manually](https://docs.aws.amazon.com/lambda/latest/dg/gettingstarted-package.html)
> 