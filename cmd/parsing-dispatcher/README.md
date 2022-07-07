# Parsing Dispatcher

The parsing dispatcher consumes messages from the parsing queue and launches jobs that parse blockchain logs for event information.

## Command Line Parameters

The parsing dispatcher depends on the graph database to store NFT information, on the events database to store transfers/sales and on the jobs database to persist failures.
It also depends on a NSQ lookup to retrieve parsing jobs, and a NSQ server to publish downstream jobs to the addition queue.
Finally, the Lambda name provides access to the corresponding parsing worker on AWS Lambda.

```
Usage of parsing-dispatcher:
  -l, --log-level string                log level (default "info")

  -g, --graph-database string           Postgres connection details for graph database (default "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=graph sslmode=disable")
  -e, --events-database string          Postgres connection details for events database (default "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=events sslmode=disable")
  -j, --jobs-database string            Postgres connection details for jobs database (default "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=jobs sslmode=disable")
  -k, --nsq-lookups []string            addresses for NSQ lookups to bootstrap consuming (default "127.0.0.1:4161")
  -q, --nsq-server string               address for NSQ server to produce messages (default "127.0.0.1:4150")
  -n, --lambda-name string              name of the lambda function to invoke (default "parsing-worker")

      --db-connection-limit uint        maximum number of database connections, -1 for unlimited (default 128)
      --db-idle-connection-limit uint   maximum number of idle connections (default 32)
      --height-range uint               maximum heights per parsing job (default 10)
      --rate-limit uint                 maximum number of API requests per second (default 10)
      --lambda-concurrency uint         maximum number of concurrent Lambda invocations (default 100)

      --min-backoff duration            minimum backoff duration for NSQ consumers (default "20s")
      --max-backoff duration            makimum backoff duration for NSQ consumers (default "10m")

      --dry-run                         executing as dry run disables invocation of Lambda function
```

## Environment Variables

Additionally to command line parameters, the parsing dispatcher requires a valid AWS environment.
Please make sure that the role associated with the container has the necessary permissions to invoke Lambdas.
Otherwise, you need to make sure that valid credentials are provided, and the region needs to be set regardless.

```sh
export AWS_DEFAULT_REGION="eu-west-1"
export AWS_ACCESS_KEY_ID="ABCDEFGHIJKLMNOPQRST"
export AWS_SECRET_ACCESS_KEY="AbCdEfGhIkLmNoPqRsTuVw12345+AbCdEfGhI+Ab"
```

## Database Address — Data Source Name

The database addresses are given in DSN (Data Source Name) format, which is a string that describes the parameters of the connection to the database.
Format of the string is the following:

```
host=localhost user=database-user password=password dbname=database-name port=5432 sslmode=disable
```
