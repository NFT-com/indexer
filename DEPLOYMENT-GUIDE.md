# Deployment Guide

## Build the containers

In order to run the indexer the first step is to build the container images.

For this the command bellow allows building and tagging the containers.
Replace `<name>` with:

* api
* jobwatcher
* parsingdispatcher
* watcher

`docker build . -f Dockerfile.<name> -t indexer:<name>`

## Running the containers

### Job API

Job API allows creating, listing, and updating discovery and parsing jobs.
Flags with descriptions [here.](cmd/jobs-api/README.md)

#### Requirements

* Postgres

#### Starting the container

`docker run indexer:api -d "host=<postgres_host> port=<postgres_port> user=<postgres_user> password=<postgres_password> dbname=jobs sslmode=<postgres_sslmode>"`

### Job Watcher

Job Watcher watches the dispatcher and parsing websockets for new updates and pushes them into their respective queue.
Flags with descriptions [here.](cmd/jobs-watcher/README.md)

#### Requirements

* Jobs API
* Redis

#### Starting the container

`docker run indexer:jobwatcher -a <jobs_api_url> -u <redis_url>`

### Parsing Dispatcher

Parsing Dispatcher consumes messages from the queue and launches lambdas.
Flags with descriptions [here.](cmd/parsing-dispatcher/README.md)

#### Requirements

* Postgres
* Jobs API
* Redis
* AWS Credentials in Environment

#### Starting the container

`docker run indexer:parsingdispatcher -u <redis_url> -a <jobs_api_url> -d "port=<postgres_port> user=<postgres_user> password=<postgres_password> dbname=chains sslmode=<postgres_sslmode>"`


(Related issue: https://github.com/NFT-com/indexer/issues/47)
** TODO CHAIN-WATCHER **
