# Getting Started Guided

This guide aims at helping users run the project.

## Requirements

* [Docker](https://docs.docker.com/get-docker/)
* [Docker Compose](https://docs.docker.com/compose/install/)

## Running the Indexer

In order to run the indexer locally, a PostgreSQL and a Redis connection are required.
The `docker-compose.yaml` contains the configuration to deploy local PostgreSQL and Redis instances.

### Deploying the Parsing Workers

Deploying the parsing workers requires building the binary, zipping it and deploying it to AWS.
There are currently two parsers:

* ERC721 Transfer Parser (`a1a4105017a0e93b98cff7ddc33b58993cd40502bcbec24e715e99ec47b964c0`)
* OpenSea OrdersMatched Parser (`485bd9051f5399c862cf3c08d8cea266459c4ee2d2acdbf2c1eb4b324e4983b2`)

Checkout the [deployment guide](DEPLOYMENT-GUIDE.md) if you want to deploy it manually.

### Starting the jobs-api, jobs-watcher, chain-bootstrapper, chain-watcher

```shell
docker-compose up postgres -d
```

This will create jobs-api, jobs-watcher, chain-bootstrapper and chain-watcher instances.
The bootstrapper and chain-watcher will create parsing jobs with the Chain URL `wss://mainnet.infura.io/ws/v3/d7b15235a515483490a5b89644221a71`, Chain Type `web3`, Address `0x4a537f61ef574153664c0dbc8c8f4b900cacbe5d`, Event Type `0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef` and Standard Type `ERC721`.

### Starting the Parsing Dispatcher

First, build the binary by running `go build`.
Then, export to the following environment variables:

* `AWS_ACCESS_KEY_ID=<aws_id>`
* `AWS_SECRET_ACCESS_KEY=<secret>`
* `AWS_DEFAULT_REGION=<region>`

These variables should target the same account and region that the workers were deployed in.

Run `./dispatcher -u <redis> -a <api> -d "host=<db> port=<db_port> user=<db_user> password=<db_password> dbname=<db_database> sslmode=<db_sslmode>"`
