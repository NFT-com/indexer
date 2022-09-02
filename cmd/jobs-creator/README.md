# Jobs Creator

The jobs creator's role is to watch the chain and instantiate parsing jobs to process and persist the chain's data into an index.
On every refresh cycle, it will figure out which jobs are missing for the collections and marketplaces in the database. It will then complete the created jobs starting at the lowest missing heights, for the collections and marketplaces missing at those heights.

## Command Line Parameters

The jobs creator depends on the graph database and the jobs database in order to figure out which jobs are missing.
It also depends on a websocket URL for an Ethereum node JSON RPC API, where it can listen for new block headers, and an NSQ server where it publishes the created jobs.

```
Usage of jobs-creator:
  -l, --log-level string                severity level for log output (default "info")

  -g, --graph-database string           Postgres connection details for graph database (default "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=graph sslmode=disable")
  -j, --jobs-database string            Postgres connection details for jobs database (default "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=jobs sslmode=disable")
  -w, --websocket-url string            Websocket URL for Ethereum JSON RPC API connection (default "ws://127.0.0.1:8545")
  -q, --nsq-server string               address for NSQ server to produce messages (default "127.0.0.1:4150")

  
      --db-connection-limit uint        maximum number of open database connections (default 16)
      --db-idle-connection-limit uint   maximum number of idle database connections (default 4)
      --write-interval duration         interval between checks for job writing (default 1s)
      --address-limit uint              maximum number of addresses to include in a single job (default 10)
      --height-limit uint               maximum number of heights to include in a single job (default 10)
```

## Database Address — Data Source Name

The database addresses are given in DSN (Data Source Name) format, which is a string that describes the parameters of the connection to the database.
Format of the string is the following:

```
host=localhost user=database-user password=password dbname=database-name port=5432 sslmode=disable
```
