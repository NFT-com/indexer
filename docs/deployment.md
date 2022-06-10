# Deployment Guide

This guide's purpose is to allow you to deploy the Indexer architecture in order to fill its database, to be served by the Analytics API.

1. [Pre-requisites](#pre-requisites)
   1. [Amazon Web Services](#amazon-web-services)
   2. [PostgreSQL](#postgresql)
   3. [NSQ](#nsq)
   4. [Docker](#docker)
      1. [Building the Images](#building-the-images)
2. [Deployment](#deployment)
   1. [Jobs Creator](#jobs-creator)
   2. [Jobs Watcher](#jobs-watcher)
   3. [Parsing Dispatcher](#parsing-dispatcher)
   4. [Action Dispatcher](#action-dispatcher)

## Pre-requisites

The Indexer requires two separate databases to function:

* [PostgreSQL](#postgresql), which is used to persist indexed data;
* [NSQ](#nsq), which is used by RabbitMQ as the backend of the jobs queue.

PostgreSQL and NSQ need to be deployed before any of the services described in the guide below.
You are free to run them any way you want, as long as they are accessible on the network on which the Indexer services are deployed.

Indexer services can be launched either using their binaries, or the docker images that can be built with the Dockerfiles within this repository.
In order to use Docker images, it is required to set [Docker](https://docs.docker.com/get-docker/) up on your machine.

### Amazon Web Services

Since the Indexer pipeline uses [AWS Lambdas](https://aws.amazon.com/lambda/) to run its parsing and action workers, the functions used by those lambdas need to be deployed on the cloud before any worker can be instantiated.

Setting up the infrastructure to run this locally is very complex, so it is recommended to deploy the worker functions to AWS using the [`pipeline` branch](https://github.com/NFT-com/indexer/tree/pipeline) of this repository.

1. `git checkout pipeline`
2. Make sure the branch is up-to-date with the branch you are working on
   1. If it is not and that you are working with master, please rebase the branch against master and push your changes to the remote `pipeline` branch.
   2. If it is not and that you are working with your own custom branch, please do not push anything on the remote `pipeline` branch and instead keep your changes local.
3. `cd ./pipeline`
4. `GOOS=linux GOARCH=amd64 go build -o worker ../cmd/parsing-worker`
5. `zip parsing.zip worker`
6. `GOOS=linux GOARCH=amd64 go build -o worker ../cmd/action-worker`
7. `zip action.zip worker`
8. Now, upload the two archives to AWS
   1. Either by [following this guide](https://docs.aws.amazon.com/lambda/latest/dg/gettingstarted-package.html) to do it manually
   2. Or using the YAML scripts that are in the `pipeline` folder on the `pipeline` branch.
      1. Install [Pulumi](https://www.pulumi.com/)
      2. Run `run pulumi up`

### PostgreSQL

You are free to run PostgreSQL however you want, but if you are unfamiliar with it, you can either get started by [downloading](https://www.postgresql.org/download/) and [installing it natively](https://www.postgresql.org/docs/current/tutorial-install.html) on your platform, or use the official [PostgreSQL Docker image](https://hub.docker.com/_/postgres).

> ⚠️ Warning! In order for the Indexer to work, the PostgreSQL database needs to be set up.
> Tables have to be created, and some also should be populated.
> This can be done by executing the SQL scripts from the `./sql` folder at the root of the repository in PostgreSQL.
> For the Docker image, simply mount the `${PWD}/sql` folder from this repository at `/docker-entrypoint-initdb.d/`.

If you want to use Docker, you first need to [set the network up](#docker) before you can run that command.

```bash
docker run  -d
            --name="postgres"
            --network="indexer"
            -e "POSTGRES_USER=postgres"
            -e "POSTGRES_PASSWORD=postgres"
            -e "POSTGRES_DB=postgres"
            -p "5432:5432"
            -v "$PWD/sql/:/docker-entrypoint-initdb.d/"
            postgres
```

> ⚠️ Warning! If you update the SQL files and want to redeploy them, you need to either manually log into the container and run `psql` commands to execute your changes, or to shut down the container, run `docker volume prune` and restart it.

### NSQ

Just like for PostgreSQL, NSQ just needs to be accessible by your services, so you can [install it natively](https://nsq.io/deployment/installing.html) or use the [official Docker image](https://hub.docker.com/r/nsqio/nsq).

If you want to use Docker, you first need to [set the network up](#docker) before you can run that command.

```bash
docker run  -d
            --network="indexer"
            -p "4160:4160"
            -p "4161:4161"
            --entrypoint "/nsqlookupd"
            nsqio/nsq


docker run  -d
            --network="indexer"
            -p "4150:4150"
            -p "4151:4151"
            --entrypoint "/nsqd --lookupd-tcp-address=nsqlookupd:4160"
            nsqio/nsq

# If you want more visibility into the queue run the command bellow and check the http://localhost:4171.
docker run  -d
            --network="indexer"
            -p "4171:4171"
            --entrypoint "/nsqadmin --lookupd-http-address=nsqlookupd:4161"
            nsqio/nsq
```

### Docker

When using Docker, it is essential to start by creating a network on which the services should run.
Then, each `docker run` command should be given the `--network=indexer` parameter.

```bash
docker network create "indexer"
```

#### Building the Images

The next step is to build the Docker images for each of the services.
You need to build and tag an image for each of the Dockerfiles within the `./cmd/*` directories.
The following command does that for you, and ignores the directories for workers, which do not contain Dockerfiles.

```bash
for d in cmd/* ;
do ;
   name=$(echo "$d" | cut -c 5-) ;
   if [[ "$name" == *-worker ]] ; then ;
      continue ;
   fi ;
   docker build . -f cmd/"$name"/Dockerfile -t indexer-"$name" ;
done
```

## Deployment

### Jobs Creator

The jobs creator's role is to watch the chain and instantiate parsing jobs to process and persist the chain's data into an index.
If the jobs creator stops, it retrieves the last job saved in the API upon restarting and starts from that height instead of `0`.
See the [jobs creator readme](../cmd/jobs-creator/README.md) for more details about its flags.

The jobs creator requires having access to an Ethereum node's API on both `WSS` and `HTTPS` and a populated [PostgreSQL](#postgresql) instance.

```bash
# Using the binary.
./jobs-creator  -n="https://mainnet.infura.io/v3/522abfc7b0f04847bbb174f026a7f83e"
                -w="wss://mainnet.infura.io/ws/v3/dc16acf06a1e7c0dbb5e7958983fb5ba"
                --graph-database="host=172.17.0.100 port=5432 user=admin password=mypassword dbname=postgres sslmode=disable"
                --jobs-database="host=172.17.0.100 port=5432 user=admin password=mypassword dbname=postgres sslmode=disable"
```

```bash
# Using Docker.
docker run  -d
            --network="indexer"
            --name="jobs-creator"
            indexer-jobs-creator
              --node-url="https://mainnet.infura.io/v3/522abfc7b0f04847bbb174f026a7f83e"
              --websocket-url="wss://mainnet.infura.io/ws/v3/dc16acf06a1e7c0dbb5e7958983fb5ba"
              --graph-database="host=172.17.0.100 port=5432 user=admin password=mypassword dbname=postgres sslmode=disable"
              --jobs-database="host=172.17.0.100 port=5432 user=admin password=mypassword dbname=postgres sslmode=disable"
```

### Jobs Watcher

The jobs watcher watches the [PostgreSQL database](#postgresql) for new jobs from the [jobs creator](#jobs-creator) and pushes them into their respective [queue](#nsq).
See the [jobs watcher readme](../cmd/jobs-watcher/README.md) for more details about its flags.

```bash
# Using the binary.
./jobs-watcher  -n="172.17.0.100:4150"
                --jobs-database="host=172.17.0.100 port=5432 user=admin password=mypassword dbname=postgres sslmode=disable"
```

```bash
# Using Docker.
docker run  -d
            --network="indexer"
            --name="jobs-watcher"
            indexer-jobs-watcher
              --nsq-address="172.17.0.100:4150"
              --jobs-database="host=172.17.0.100 port=5432 user=admin password=mypassword dbname=postgres sslmode=disable"
```

### Parsing Dispatcher

The parsing dispatcher consumes messages from the [queue](#nsq) and launches parsing jobs on [AWS Lambdas](#amazon-web-services).
See the [parsing dispatcher readme](../cmd/parsing-dispatcher/README.md) for more details about its flags.

In order for the parsing dispatcher to be allowed to instantiate workers on AWS Lambda, it requires credentials to authenticate.
Those [credentials should be set in the environment](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html) of the machine that runs the dispatcher.

```bash
# Using the binary.
./parsing-dispatcher  -a="172.17.0.100:4161"
                      --jobs-database="host=172.17.0.100 port=5432 user=admin password=mypassword dbname=postgres sslmode=disable"
                      --events-database="host=172.17.0.100 port=5432 user=admin password=mypassword dbname=postgres sslmode=disable"
```

```bash
# Using Docker.
docker run  -d
            --network="indexer"
            --name="parsing-dispatcher"
            -e AWS_REGION="eu-west-1"
            -e AWS_ACCESS_KEY_ID="E283E205A2CA9FE4A032"
            -e AWS_SECRET_ACCESS_KEY="XDklicgtXc8Wgx0x9Rmlpdrfybn+Gjxh3YyWz+fR"
            indexer-parsing-dispatcher
              --nsq-lookup-address="172.17.0.100:4161"
              --jobs-database="host=172.17.0.100 port=5432 user=admin password=mypassword dbname=postgres sslmode=disable"
              --events-database="host=172.17.0.100 port=5432 user=admin password=mypassword dbname=postgres sslmode=disable"
```

### Action Dispatcher

The action dispatcher consumes messages from the [queue](#nsq) and launches jobs on [AWS Lambdas](#amazon-web-services).
Those jobs can act in several ways, hence the name.
See the [action dispatcher readme](../cmd/action-dispatcher/README.md) for more details about its flags.

In order for the action dispatcher to be allowed to instantiate workers on AWS Lambda, it requires credentials to authenticate.
Those [credentials should be set in the environment](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html) of the machine that runs the dispatcher.

```bash
# Using the binary.
./action-dispatcher  -a="172.17.0.100:4161"
                      --jobs-database="host=172.17.0.100 port=5432 user=admin password=mypassword dbname=postgres sslmode=disable"
                      --events-database="host=172.17.0.100 port=5432 user=admin password=mypassword dbname=postgres sslmode=disable"
```

```bash
# Using Docker.
docker run  -d
            --network="indexer"
            --name="action-dispatcher"
            -e AWS_REGION="eu-west-1"
            -e AWS_ACCESS_KEY_ID="E283E205A2CA9FE4A032"
            -e AWS_SECRET_ACCESS_KEY="XDklicgtXc8Wgx0x9Rmlpdrfybn+Gjxh3YyWz+fR"
            indexer-action-dispatcher
              --nsq-lookup-address="172.17.0.100:4161"
              --jobs-database="host=172.17.0.100 port=5432 user=admin password=mypassword dbname=postgres sslmode=disable"
              --events-database="host=172.17.0.100 port=5432 user=admin password=mypassword dbname=postgres sslmode=disable"
```
