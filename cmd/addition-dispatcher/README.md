# Addition Dispatcher

The addition dispatcher consumes messages from the addition queue and launches jobs that can add the information for newly minted NFTs to the graph database.

## Command Line Parameters

The addition dispatcher depends on the graph database to store data and the jobs database to persist job failures.
It also requires a NSQ lookup to consume jobs from the addition queue.
Finally, the Lambda name provides access to the corresponding parsing worker on AWS Lambda.

```
Usage of action-dispatcher:
  -l, --log-level string                severity level for log output (default "info")
  
  -g, --graph-database string           Postgres connection details for graph database (default "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=graph sslmode=disable")
  -j, --jobs-database string            Postgres connection details for jobs database (default "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=jobs sslmode=disable")
  -k, --nsq-lookups []string            addresses for NSQ lookups to bootstrap consuming (default "127.0.0.1:4161")
  -n, --lambda-name string              name of Lambda function for invocation (default "action-worker")

      --db-connection-limit uint        maximum number of open database connections (default 128)
      --db-idle-connection-limit uint   maximum number of idle database connections (default 32)
      --lambda-concurrency uint         maximum number of concurrent Lambda invocations (default 900)
      --rate-limit uint                 maximum number of API requests per second (default 100)
      
      --min-backoff duration            minimum backoff duration for NSQ consumers (default "1s")
      --max-backoff duration            maximum backoff duration for NSQ consumers (default "10m")
      --max-attempts uint16             maximum number of attempts per job (default "3")

      --dry-run                         executing as dry run disables invocation of Lambda function
```

## Environment Variables

Additionally to command line parameters, the addition dispatcher requires a valid AWS environment.
Please make sure that the role associated with the container has the necessary permissions to invoke Lambdas.
Otherwise, you need to make sure that valid credentials are provided, and the region needs to be set regardless.

```sh
export AWS_DEFAULT_REGION="eu-west-1"
export AWS_ACCESS_KEY_ID="ABCDEFGHIJKLMNOPQRST"
export AWS_SECRET_ACCESS_KEY="AbCdEfGhIkLmNoPqRsTuVw12345+AbCdEfGhI+Ab"
```

## Database Address — Data Source Name

The database addresses are given in DSN (Data Source Name) format, which is a string that describes the parameters of the connection to the database.
The format is the following:

```
host=localhost user=database-user password=password dbname=database-name port=5432 sslmode=disable
```
