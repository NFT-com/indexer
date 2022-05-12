# Action Dispatcher

The action dispatcher consumes messages from the queue and launches jobs that can act in several ways.

Actions:

* **Addition**: Adds an NFT.
* **OwnerChange**: Updates an NFT's owner.

## Usage

```
Usage of action-dispatcher:
  -r, --aws-region string               AWS region for Lambda invocation (default "eu-west-1")
      --db-connection-limit uint        maximum number of open database connections (default 128)
      --db-idle-connection-limit uint   maximum number of idle database connections (default 32)
      --dry-run                         executing as dry run disables invocation of Lambda function
  -g, --graph-database string           Postgres connection details for graph database (default "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")
  -j, --job-database string            Postgres connection details for job database (default "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")
      --lambda-concurrency uint         maximum number of concurrent Lambda invocations (default 900)
  -n, --lambda-name string              name of Lambda function for invocation (default "action-worker")
  -l, --log-level string                severity level for log output (default "info")
      --rate-limit uint                 maximum number of API requests per second (default 100)
  -d, --redis-database int              Redis database number (default 1)
  -u, --redis-url string                URL for Redis server connection (default "127.0.0.1:6379")
```

## Database Address — Data Source Name

Data Source Name (DSN) is the string specified describing how the connection to the database should be established.
The string's format is the following:

```
host=localhost user=database-user password=password dbname=database-name port=5432 sslmode=disable
```
