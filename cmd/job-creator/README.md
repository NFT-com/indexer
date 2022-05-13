# Job Creator

The job creator's role is to watch the chain and instantiate parsing jobs to process and persist the chain's data into an index.
If the job creator stops, it retrieves the last job saved in the API upon restarting and starts from that height instead of 0.

## Usage

```
Usage of job-creator:
      --db-connection-limit uint        maximum number of open database connections (default 16)
      --db-idle-connection-limit uint   maximum number of idle database connections (default 4)
  -g, --graph-database string           postgresql connection details for graph database (default "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")
      --height-range uint               maximum heights to include in a single job (default 10)
  -j, --job-database string             postgresql connection details for job database (default "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")
  -l, --log-level string                severity level for log output (default "info")
  -n, --node-url string                 http URL for Ethereum JSON RPC API connection (default "http://127.0.0.1:8545")
      --pending-limit uint              maximum number of pending jobs per combination (default 1000)
  -w, --websocket-url string            websocket URL for Ethereum JSON RPC API connection (default "ws://127.0.0.1:8545")
      --write-interval duration         interval between checks for job writing (default 1s)
```

## Database Address — Data Source Name

Data Source Name (DSN) is the string specified describing how the connection to the database should be established.
Format of the string is the following:

```
host=localhost user=database-user password=password dbname=database-name port=5432 sslmode=disable
```