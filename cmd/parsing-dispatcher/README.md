# Parsing Dispatcher

The parsing dispatcher consumes messages from the queue and launches parsing jobs.

## Usage

```
Usage of parsing-dispatcher:
  -r, --aws-region string               AWS region for Lambda invocation (default "eu-west-1")
      --db-connection-limit uint        maximum number of database connections, -1 for unlimited (default 128)
      --db-idle-connection-limit uint   maximum number of idle connections (default 32)
      --dry-run                         executing as dry run disables invocation of Lambda function
  -e, --events-database string          Postgres connection details for events database (default "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")
      --height-range uint               maximum heights per parsing job (default 10)
  -j, --jobs-database string            Postgres connection details for job database (default "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")
      --lambda-concurrency uint         maximum number of concurrent Lambda invocations (default 100)
  -n, --lambda-name string              name of the lambda function to invoke (default "parsing-worker")
  -l, --log-level string                log level (default "info")
      --rate-limit uint                 maximum number of API requests per second (default 10)
  -d, --redis-database int              Redis database number (default 1)
  -u, --redis-url string                URL for Redis server connection (default "127.0.0.1:6379")

```

## Database Address — Data Source Name

Data Source Name (DSN) is the string specified describing how the connection to the database should be established.
Format of the string is the following:

```
host=localhost user=database-user password=password dbname=database-name port=5432 sslmode=disable
```
