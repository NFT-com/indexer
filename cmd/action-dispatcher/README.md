# Action Dispatcher

This pipeline allows the dispatcher to deploy several actions.

Actions:

* Addition: Gets the information from networks.
* OwnerChange: Updates the nft's owner with the new owner.

## Usage

```
Usage of action-dispatcher:
  -l, --log-level string                severity level for log output (default "info")
  -j, --graph-database string           postgres connection details for graph database (default "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=graph sslmode=disable")
  -j, --jobs-database string            postgres connection details for jobs database (default "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=jobs sslmode=disable")
  -u, --redis-url string                url for redis server connection (default "127.0.0.1:6379")
  -d, --redis-database int              redis database number (default 1)
  -r, --aws-region string               aws region for lambda invocation (default "eu-west-1")
      --db-connection-limit uint        maximum number of open database connections (default 128)
      --db-idle-connection-limit uint   maximum number of idle database connections (default 32)
      --lambda-name string              name of lambda function for invocation (default "action-worker")
      --lambda-concurrency uint         maximum number of concurrent lambda invocations (default 900)
      --rate-limit uint                 maximum number of api requests to ethereum json rpc api (default 100)
      --dry-run boolean                 executing as dry run disables invocation of lambda function (default "false")
```

## Database Address - Data Source Name

Data Source Name (DSN) is the string specified describing how the connection to the database should be established. The
string's format is the following:

```
host=localhost user=database-user password=password dbname=database-name port=5432 sslmode=disable
```
