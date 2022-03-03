# Indexer Service

## Local Development

### Requirements

* [Docker](https://docs.docker.com/get-docker/)
* [Docker Compose](https://docs.docker.com/compose/install/)

### Testing Locally

In order to run a local jobs api a postgres connection is required.
The `docker-compose.yaml` contains a docker configuration for a postgres database.

#### Running the Postgres Database

```shell

docker-compose up postgres -d

```

#### Connection from the Jobs API to the Postgres Database

Running the api and the postgres database locally the flag `-d "host=<host> port=<port> user=<user> password=<pass> dbname=jobs sslmode=disable"` must be set.
Replace the placeholders with the correct information.
