# Action Dispatcher

The action dispatcher consumes messages from the action queue and launches jobs that can execute a number of different actions:

* **addition**: add a newly minted NFT to the graph database
* **owner change**: update an NFT owner after a transfer

## Usage

```
Usage of action-dispatcher:
  -l, --log-level string                severity level for log output (default "info")
  
  -g, --graph-database string           Postgres connection details for graph database (default "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=graph sslmode=disable")
  -j, --jobs-database string            Postgres connection details for jobs database (default "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=jobs sslmode=disable")
  -a, --nsq-lookup-address string       NSQ lookup address (default "127.0.0.1:4161")
  -r, --aws-region string               AWS region for Lambda invocation (default "eu-west-1")
  -n, --lambda-name string              name of Lambda function for invocation (default "action-worker")

      --db-connection-limit uint        maximum number of open database connections (default 128)
      --db-idle-connection-limit uint   maximum number of idle database connections (default 32)
      --lambda-concurrency uint         maximum number of concurrent Lambda invocations (default 900)
      --rate-limit uint                 maximum number of API requests per second (default 100)
      
      --dry-run                         executing as dry run disables invocation of Lambda function
```

## Database Address — Data Source Name

The database addresses are given in DSN (Data Source Name) format, which is a string that describes the parameters of the connection to the database.
The string's format is the following:

```
host=localhost user=database-user password=password dbname=database-name port=5432 sslmode=disable
```
