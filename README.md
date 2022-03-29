# Indexer Service

## Binaries

* [JOBS API](./cmd/jobs-api/README.md)

## APIs

* [JobsAPI](./cmd/jobs-api/API.md)

## Local Development

### Requirements

* [Docker](https://docs.docker.com/get-docker/)
* [Docker Compose](https://docs.docker.com/compose/install/)

### Testing Locally

In order to run the jobs API locally, a PostgreSQL connection is required.
The `docker-compose.yaml` contains the configuration to deploy a PostgreSQL database.

#### Running the PostgreSQL Database

```shell

docker-compose up postgres -d

```

#### Connection from the Jobs API to the PostgreSQL Database

In order to run the api and its database locally, the flag `-d "host=<host> port=<port> user=<user> password=<pass> dbname=jobs sslmode=disable"` must be set.
Replace the placeholders with the correct information.
