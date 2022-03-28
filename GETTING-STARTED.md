# Getting Started Guided

This getting started guide will help the user run the project.

## Requirements

* [Docker](https://docs.docker.com/get-docker/)
* [Docker Compose](https://docs.docker.com/compose/install/)

## Running the indexer

In order to run the indexer locally, a PostgreSQL and a Redis connection is required.
The `docker-compose.yaml` contains the configuration to deploy a PostgreSQL and Redis instance.
In order to run the indexer it requires three steps.

### Deploying the parsing workers

In order deploy the parsing workers it requires to build the binary zip it and deploy it to aws.
There are two parser currently:
* ERC721 Transfer Parser (a1a4105017a0e93b98cff7ddc33b58993cd40502bcbec24e715e99ec47b964c0)
* OpenSea OrdersMatched Parser (485bd9051f5399c862cf3c08d8cea266459c4ee2d2acdbf2c1eb4b324e4983b2)

*TODO*

### Starting the jobs-api, jobs-watcher, chain-bootstrapper, chain-watcher

```shell

docker-compose up postgres -d

```

This will create jobs-api, jobs-watcher, chain-bootstrapper and chain-watcher instances.
The bootstrapper and chain-watcher will create parsing jobs with the Chain URL `wss://mainnet.infura.io/ws/v3/d7b15235a515483490a5b89644221a71`, Chain Type `web3`, Address `0x4a537f61ef574153664c0dbc8c8f4b900cacbe5d`, Event Type `0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef` and Standard Type `ERC721`.

### Starting the parsing dispatcher

First build the binary using the go build cli command.
Then export to the environment variables:
* `AWS_ACCESS_KEY_ID=<aws_id>`
* `AWS_SECRET_ACCESS_KEY=<secret>`
* `AWS_DEFAULT_REGION=<region>`
These variables should target the same account and region that the workers were deployed.

Run `./dispatcher -u <redis> -a <api> -d "host=<db> port=<db_port> user=<db_user> password=<db_password> dbname=<db_database> sslmode=<db_sslmode>"`
