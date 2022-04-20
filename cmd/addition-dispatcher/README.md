# Addition Dispatcher

Addition Dispatcher consumes messages from the queue and launches jobs.

## Usage

```
Usage of parsing-dispatcher:
  -q, --addition-queue string   addition queue name (default "addition")
  -a, --api string              jobs api base endpoint
  -j, --jobs int                amount of concurrent lambda calls (default 4)
  -p, --prefetch int            amount of queued messages to prefetch on init (default 5)
  -i, --poll-duration duration  time between polls on queue (default 1s)
  -d, --db string               data source name for database connection
  -l, --log-level string        log level (default "info")
  -c, --tag string              rmq producer tag (default "dispatcher-agent")
  --database int                redis database number (default 1)
  -n, --network string          redis network type (default "tcp")
  -u, --url string              redis server connection url
  -r, --aws-region              aws lambda region (default "eu-west-1")
```

## Database Address - Data Source Name

Data Source Name (DSN) is the string specified describing how the connection to the database should be established.
The string's format is the following:

```
host=localhost user=database-user password=password dbname=database-name port=5432 sslmode=disable
```