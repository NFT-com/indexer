# Getting Started Guided

This guide aims at helping users run the project.

## Requirements

* [Docker](https://docs.docker.com/get-docker/)
* [Docker Compose](https://docs.docker.com/compose/install/)

## Running the Indexer

In order to run the indexer locally, a PostgreSQL and a Redis connection are required.
The `docker-compose.yaml` contains the configuration to deploy local PostgreSQL and Redis instances.

### Deploying the Workers

Deploying the workers requires building the binaries, zipping it and deploying it to AWS.
There are currently two workers:

* Action Worker (`action-worker`)
* Parsing Worker (`parsing-worker`)

Checkout the [deployment guide](DEPLOYMENT-GUIDE.md) if you want to deploy it manually.

### Starting the Components

```shell
docker-compose up postgres -d
```

This creates action-dispatcher, jobs-creator, jobs-watcher and parsing-dispatcher instances.
The jobs-creator creates parsing jobs with the Chain URL `wss://mainnet.infura.io/ws/v3/d7b15235a515483490a5b89644221a71` for all collections in the database.

### Starting the Parsing Dispatcher

First, build the binary by running `go build`.
Then, export to the following environment variables:

* `AWS_ACCESS_KEY_ID=<aws_id>`
* `AWS_SECRET_ACCESS_KEY=<secret>`
* `AWS_DEFAULT_REGION=<region>`

These variables should target the same account and region that the workers were deployed in.

Run `./dispatcher -u <redis> -a <api> -d "host=<db> port=<db_port> user=<db_user> password=<db_password> dbname=<db_database> sslmode=<db_sslmode>"`
