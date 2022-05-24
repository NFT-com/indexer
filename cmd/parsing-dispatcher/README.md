# Parsing Dispatcher

The parsing dispatcher consumes messages from the queue and launches parsing jobs.

## Usage

```
Usage of parsing-dispatcher:
  -l, --log-level string                log level (default "info")

  -j, --jobs-database string            Postgres connection details for jobs database (default "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=jobs sslmode=disable")
  -e, --events-database string          Postgres connection details for events database (default "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=events sslmode=disable")
  -u, --redis-url string                Redis server URL (default "127.0.0.1:6379")
  -d, --redis-database int              Redis database number (default 1)
  -r, --aws-region string               aws region for Lambda invocation (default "eu-west-1")
  -n, --lambda-name string              name of the lambda function to invoke (default "parsing-worker")

      --db-connection-limit uint        maximum number of database connections, -1 for unlimited (default 128)
      --db-idle-connection-limit uint   maximum number of idle connections (default 32)
      --height-range uint               maximum heights per parsing job (default 10)
      --rate-limit uint                 maximum number of API requests per second (default 10)
      --lambda-concurrency uint         maximum number of concurrent Lambda invocations (default 100)

      --dry-run                         executing as dry run disables invocation of Lambda function
```



	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "log level")

	pflag.StringVarP(&flagJobsDB, "jobs-database", "j", "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=jobs sslmode=disable", "Postgres connection details for jobs database")
	pflag.StringVarP(&flagEventsDB, "events-database", "e", "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=events sslmode=disable", "Postgres connection details for events database")
	pflag.StringVarP(&flagRedisURL, "redis-url", "u", "127.0.0.1:6379", "Redis server url")
	pflag.IntVarP(&flagRedisDB, "redis-database", "d", 1, "Redis database number")
	pflag.StringVarP(&flagAWSRegion, "aws-region", "r", "eu-west-1", "AWS region for Lambda invocation")
	pflag.StringVarP(&flagLambdaName, "lambda-name", "n", "parsing-worker", "name of the Lambda function to invoke")

	pflag.UintVar(&flagOpenConnections, "db-connection-limit", 128, "maximum number of database connections, -1 for unlimited")
	pflag.UintVar(&flagIdleConnections, "db-idle-connection-limit", 32, "maximum number of idle connections")
	pflag.UintVar(&flagHeightRange, "height-range", 10, "maximum heights per parsing job")
	pflag.UintVar(&flagRateLimit, "rate-limit", 10, "maximum number of API requests per second")
	pflag.UintVar(&flagLambdaConcurrency, "lambda-concurrency", 100, "maximum number of concurrent Lambda invocations")

	pflag.BoolVar(&flagDryRun, "dry-run", false, "executing as dry run disables invocation of Lambda function")

## Database Address — Data Source Name

The database addresses are given in DSN (Data Source Name) format, which is a string that describes the parameters of the connection to the database.
Format of the string is the following:

```
host=localhost user=database-user password=password dbname=database-name port=5432 sslmode=disable
```
