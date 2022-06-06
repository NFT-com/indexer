# Parsing Dispatcher

The parsing dispatcher consumes messages from the queue and launches parsing jobs.

## Usage

```
Usage of parsing-dispatcher:
  -l, --log-level string                log level (default "info")

  -j, --jobs-database string            Postgres connection details for jobs database (default "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=jobs sslmode=disable")
  -e, --events-database string          Postgres connection details for events database (default "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=events sslmode=disable")
  -q, --nsq-server string               NSQ server address (default "127.0.0.1:4150")
  -r, --aws-region string               aws region for Lambda invocation (default "eu-west-1")
  -n, --lambda-name string              name of the lambda function to invoke (default "parsing-worker")

      --db-connection-limit uint        maximum number of database connections, -1 for unlimited (default 128)
      --db-idle-connection-limit uint   maximum number of idle connections (default 32)
      --height-range uint               maximum heights per parsing job (default 10)
      --rate-limit uint                 maximum number of API requests per second (default 10)
      --lambda-concurrency uint         maximum number of concurrent Lambda invocations (default 100)

      --dry-run                         executing as dry run disables invocation of Lambda function
```

## Database Address — Data Source Name

The database addresses are given in DSN (Data Source Name) format, which is a string that describes the parameters of the connection to the database.
Format of the string is the following:

```
host=localhost user=database-user password=password dbname=database-name port=5432 sslmode=disable
```
