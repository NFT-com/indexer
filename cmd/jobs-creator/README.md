# Jobs Creator

The jobs creator's role is to watch the chain and instantiate parsing jobs to process and persist the chain's data into an index.
If the jobs creator stops, it retrieves the last job saved in the API upon restarting and starts from that height instead of 0.

## Command Line Parameters

The jobs creator depends on the graph database, the jobs database, as well as both a HTTP and a WebSocket URL for an Ethereum node JSON RPC API.

```
Usage of jobs-creator:
  -l, --log-level string                severity level for log output (default "info")

  -g, --graph-database string           Postgres connection details for graph database (default "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=graph sslmode=disable")
  -j, --jobs-database string            Postgres connection details for jobs database (default "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=jobs sslmode=disable")
  -n, --node-url string                 HTTP URL for Ethereum JSON RPC API connection (default "http://127.0.0.1:8545")
  -w, --websocket-url string            Websocket URL for Ethereum JSON RPC API connection (default "ws://127.0.0.1:8545")
  
      --db-connection-limit uint        maximum number of open database connections (default 16)
      --db-idle-connection-limit uint   maximum number of idle database connections (default 4)
      --write-interval duration         interval between checks for job writing (default 1s)
      --pending-limit uint              maximum number of pending jobs per combination (default 1000)
      --height-range uint               maximum heights to include in a single job (default 10)
```

## Environment Variables

In addition tho the command line parameters, the jobs creator optionally supports AWS Managed Blockchain nodes for the HTTP URL.
If a URL on AWS Managed Blockchain is given, a valid AWS environment must be configured.
You need to make sure that valid credentials are provided  and the region has been set.

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
