# Parsing Dispatcher

The parsing dispatcher consumes messages from the queue and launches jobs.

## Usage

```
Usage of parsing-dispatcher:
  -a, --api string               jobs api base hostname and port
  -r, --aws-region string        aws lambda region (default "eu-west-1")
      --database int             redis database number (default 1)
  -d, --db string                database connection string
  -l, --log-level string         log level (default "info")
  -n, --network string           redis network type (default "tcp")
  -q, --parsing-queue string     name of the queue for parsing (default "parsing")
  -i, --poll-duration duration   time for each consumer poll (default 20s)
  -p, --prefetch int             amount of message to prefetch in the consumer (default 80)
  -t, --rate-limit int           amount of concurrent jobs for the consumer (default 500)
  -c, --tag string               rmq consumer tag (default "parsing-agent")
  -u, --url string               redis server connection url
```

## Database Address - Data Source Name

Data Source Name (DSN) is the string specified describing how the connection to the database should be established.
Format of the string is the following:

```
host=localhost user=database-user password=password dbname=database-name port=5432 sslmode=disable
```
