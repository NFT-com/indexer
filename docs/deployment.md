# Deployment Guide

This guide's purpose is to allow you to deploy the indexer architecture in order to fill its database, to be served by the Analytics API.

1. [Pre-requisites](#pre-requisites)
   1. [Amazon Web Services](#amazon-web-services)
      1. [Manually](#manually)
      2. [Pulumi](#pulumi)
   2. [PostgreSQL](#postgresql)
   3. [NSQ](#nsq)
   4. [Docker](#docker)
      1. [Building the Images](#building-the-images)
2. [Deployment](#deployment)
   1. [Jobs Creator](#jobs-creator)
   2. [Parsing Dispatcher](#parsing-dispatcher)
   3. [Addition Dispatcher](#addition-dispatcher)
   4. [Completion Dispatcher](#completion-dispatcher)

## Pre-requisites

The Indexer requires at least one database, and a pipeline for jobs.

* [PostgreSQL](#postgresql), which is used to persist indexed data, job information and events.
* [NSQ](#nsq), which provides a persistent at-least-once message queue for jobs.

PostgreSQL and NSQ need to be deployed before any of the services described in the guide below.
You are free to run them any way you want, as long as they are accessible on the network on which the Indexer services are deployed.

Indexer services can be launched either using their binaries, or the docker images that can be built with the Dockerfiles within this repository.
In order to use Docker images, it is required to set [Docker](https://docs.docker.com/get-docker/) up on your machine.

### Amazon Web Services

Since the Indexer pipeline uses [AWS Lambdas](https://aws.amazon.com/lambda/) to run its parsing and addition workers, the Lambda functions used by those workers need to be deployed on the cloud before any worker can be instantiated.

Setting up the infrastructure to run this locally is very complex, so it is recommended to deploy the worker functions to AWS. This can be done manually, or by using Pelumi to deploy automatically.

The [`pipeline` branch](https://github.com/NFT-com/indexer/tree/pipeline) of this repository contains example code on how to use Pelumi to automate deployment of Lambda functions.

#### Manually

In order to deploy the Lambda functions manually:

1. Log into the AWS console
2. Create one function named `parsing-worker` and one function named `addition-worker`
3. Make sure that the entry point is changed from `hello` to `worker`
4. We also recommend changing the timeout of the functions to 15 minutes
5. Set the desired environment variables for `LOG_LEVEL` and `NODE_URL` on the settings
6. Build the workers:
   - `GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o worker ../cmd/parsing-worker`
   - `zip parsing.zip worker`
   - `GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o worker ../cmd/addition-worker`
   - `zip addition.zip worker`
7. Manually upload each zip to the respective Lambda function.

More information is available in the [AWS docs](https://docs.aws.amazon.com/lambda/latest/dg/gettingstarted-package.html).

#### Pulumi

If you want to use Pulumi, you should make sure that the `pipeline` branch is properly rebased on the branch that you want to deploy from.
Alternatively, you can build the zip files, and then use the `pipeline` branch just for the Pulumi code.
We recommend that you integrate the available scripts in your own deployment flow outside of the repository.
Once everything is up-to-date, you should be able to simply execute `run pulumi up`.

### PostgreSQL

You are free to run PostgreSQL however you want, but if you are unfamiliar with it, you can either get started by [downloading](https://www.postgresql.org/download/) and [installing it natively](https://www.postgresql.org/docs/current/tutorial-install.html) on your platform, or use the official [PostgreSQL Docker image](https://hub.docker.com/_/postgres).

There are three different databases used by the different components:

- graph data;
- jobs data; and
- events data.

Ideally, each database should run on its own host, so that they can be scaled according to needs.

> ⚠️ Warning! In order for the Indexer to work, the PostgreSQL database needs to be set up.
> Tables have to be created, and some also should be populated.
> This can be done by executing the SQL scripts from the `./sql` folder at the root of the repository in PostgreSQL.
> For the Docker image, simply mount the `${PWD}/sql` folder from this repository at `/docker-entrypoint-initdb.d/`.

If you want to use Docker, you first need to [set the network up](#docker) before you can run that command.

```sh
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

```sh
docker run  -d \
--network="indexer" \
-p "4160:4160" \
-p "4161:4161" \
nsqio/nsq \
/nsqlookupd
```

```sh
docker run  -d \
--network="indexer" \
-p "4150:4150" \
-p "4151:4151" \
nsqio/nsq \
/nsqd \
--lookupd-tcp-address=host.docker.internal:4160 \
--broadcast-address=host.docker.internal \
--msg-timeout 15m \
--max-msg-timeout 15m

# If you want more visibility into the queue run the command bellow and check the http://localhost:4171.
docker run  -d \
--network="indexer" \
-p "4171:4171" \
nsqio/nsq \
/nsqadmin --lookupd-http-address=host.docker.internal:4161
```

### Docker

When using Docker, it is essential to start by creating a network on which the services should run.
Then, each `docker run` command should be given the `--network=indexer` parameter.

```sh
docker network create "indexer"
```

#### Building the Images

The next step is to build the Docker images for each of the services.
You need to build and tag an image for each of the Dockerfiles within the `./cmd/*` directories.
The following command does that for you, and ignores the directories for workers, which do not contain Dockerfiles.

```sh
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

A NSQ consumer will create the topic and channel it listens on when it connects to NSQ.
A NSQ producer will create the topic when it produces messages.
If a topic is created before the channels, the first channel for that topic gets all messages.
Any subsequent channels will no longer receive queued / old messages.
It might be a good idea to create the queues manually using [NSQ admin](https://nsq.io/components/nsqadmin.html).
In that case, all messages remain queued for each channel until they are consumed.

### Jobs Creator

The jobs creator's role is to watch the chain and instantiate parsing jobs to process and persist the chain's data into an index.
If the jobs creator stops, it retrieves the last job saved in the API upon restarting and starts from that height instead of `0`.
See the [jobs creator readme](../cmd/jobs-creator/README.md) for more details about its flags.

The jobs creator requires having access to an Ethereum node's API on both `WSS` and `HTTPS` and a populated [PostgreSQL](#postgresql) instance.

```sh
# Using the binary.
./jobs-creator \
-g "host=172.17.0.100 port=5432 user=immutable password=password dbname=graph sslmode=disable" \
-j "host=172.17.0.100 port=5432 user=immutable password=password dbname=jobs sslmode=disable" \
-w="wss://mainnet.infura.io/ws/v3/1234567890abcdef1234567890" \
-q "nsq.domain.com:4150"
```

```sh
# Using Docker.
docker run  -d \
--network="indexer" \
--name="jobs-creator" \
indexer-jobs-creator \
--graph-database="host=172.17.0.100 port=5432 user=immutable password=password dbname=graph sslmode=disable" \
--jobs-database="host=172.17.0.100 port=5432 user=immutable password=password dbname=jobs sslmode=disable" \
--websocket-url="wss://mainnet.infura.io/ws/v3/1234567890abcdef1234567890" \
--nsq-server "nsq.domain.com:4150"
```


### Parsing Dispatcher

The parsing dispatcher consumes messages from the [parsing queue](#nsq) and launches parsing jobs on [AWS Lambdas](#amazon-web-services).
See the [parsing dispatcher readme](../cmd/parsing-dispatcher/README.md) for more details about its flags.

In order for the parsing dispatcher to be allowed to instantiate workers on AWS Lambda, it requires credentials to authenticate.
Those [credentials should be set in the environment](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html) of the machine that runs the dispatcher.

```sh
# Using the binary.
./parsing-dispatcher \
-g "host=172.17.0.100 port=5432 user=immutable password=password dbname=gaph sslmode=disable" \
-j "host=172.17.0.100 port=5432 user=immutable password=password dbname=jobs sslmode=disable" \
-e "host=172.17.0.100 port=5432 user=immutable password=password dbname=events sslmode=disable" \
-k "nsq.domain.com:4161" \
-q "nsq.domain.com:4150" \
-n "parsing-worker"
```

```sh
# Using Docker.
docker run -d \
--network="indexer" \
--name="parsing-dispatcher" \
-e AWS_REGION="eu-west-1" \
-e AWS_ACCESS_KEY_ID="E283E205A2CA9FE4A032" \
-e AWS_SECRET_ACCESS_KEY="XDklicgtXc8Wgx0x9Rmlpdrfybn+Gjxh3YyWz+fR" \
indexer-parsing-dispatcher \
--graph-database "host=172.17.0.100 port=5432 user=immutable password=password dbname=gaph sslmode=disable" \
--jobs-database "host=172.17.0.100 port=5432 user=immutable password=password dbname=jobs sslmode=disable" \
--events-database "host=172.17.0.100 port=5432 user=immutable password=password dbname=events sslmode=disable" \
--nsq-lookups "nsq.domain.com:4161" \
--nsq-server "nsq.domain.com:4150" \
--lambda-name "parsing-worker"
```

### Addition Dispatcher

The addition dispatcher consumes messages from the [addition queue](#nsq) and launches jobs on [AWS Lambdas](#amazon-web-services).
See the [addition dispatcher readme](../cmd/addition-dispatcher/README.md) for more details about its flags.

In order for the addition dispatcher to be allowed to instantiate workers on AWS Lambda, it requires credentials to authenticate.
Those [credentials should be set in the environment](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html) of the machine that runs the dispatcher.

```sh
# Using the binary.
./addition-dispatcher \
-g "host=172.17.0.100 port=5432 user=immutable password=password dbname=gaph sslmode=disable" \
-j "host=172.17.0.100 port=5432 user=immutable password=password dbname=jobs sslmode=disable" \
-k "nsq.domain.com:4161" \
-n "addition-worker"
```

```sh
# Using Docker.
docker run -d \
--network="indexer" \
--name="addition-dispatcher" \
-e AWS_REGION="eu-west-1" \
-e AWS_ACCESS_KEY_ID="E283E205A2CA9FE4A032" \
-e AWS_SECRET_ACCESS_KEY="XDklicgtXc8Wgx0x9Rmlpdrfybn+Gjxh3YyWz+fR" \
indexer-addition-dispatcher \
--graph-database "host=172.17.0.100 port=5432 user=immutable password=password dbname=gaph sslmode=disable" \
--jobs-database "host=172.17.0.100 port=5432 user=immutable password=password dbname=jobs sslmode=disable" \
--nsq-lookups "nsq.domain.com:4161" \
--lambda-name "addition-worker"
```

### Completion Dispatcher

The completion dispatcher consumes messages from the [completion queue](#nsq) and launches jobs on [AWS Lambdas](#amazon-web-services).
See the [completion dispatcher readme](../cmd/completion-dispatcher/README.md) for more details about its flags.

In order for the completion dispatcher to be allowed to instantiate workers on AWS Lambda, it requires credentials to authenticate.
Those [credentials should be set in the environment](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html) of the machine that runs the dispatcher.

```bash
# Using the binary.
./completion-dispatcher \
-g "host=172.17.0.100 port=5432 user=immutable password=password dbname=gaph sslmode=disable" \
-e "host=172.17.0.100 port=5432 user=immutable password=password dbname=events sslmode=disable" \
-j "host=172.17.0.100 port=5432 user=immutable password=password dbname=jobs sslmode=disable" \
-k "nsq.domain.com:4161" \
-n "completion-worker"
```

```bash
# Using Docker.
docker run -d \
--network="indexer" \
--name="completion-dispatcher" \
-e AWS_REGION="eu-west-1" \
-e AWS_ACCESS_KEY_ID="E283E205A2CA9FE4A032" \
-e AWS_SECRET_ACCESS_KEY="XDklicgtXc8Wgx0x9Rmlpdrfybn+Gjxh3YyWz+fR" \
indexer-completion-dispatcher \
--graph-database "host=172.17.0.100 port=5432 user=immutable password=password dbname=gaph sslmode=disable" \
--events-database "host=172.17.0.100 port=5432 user=immutable password=password dbname=events sslmode=disable" \
--jobs-database "host=172.17.0.100 port=5432 user=immutable password=password dbname=jobs sslmode=disable" \
--nsq-lookups "nsq.domain.com:4161" \
--lambda-name "completion-worker"
```
