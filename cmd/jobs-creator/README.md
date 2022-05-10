# Jobs Creator

Jobs Creator watches the chain and pushes jobs to parse network data.
If stopped during an instantiation, the jobs creator resumes its work from where it left off.

## Usage

```
Usage of chain-watcher:
  -l, --log-level string                severity level for log output (default "info")
  -g, --graph-database string           postgres connection details for graph database (default "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=graph sslmode=disable")
  -j, --jobs-database string            postgres connection details for jobs database (default "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=jobs sslmode=disable")
  -n, --node-url string                 url for ethereum json rpc api connection (default "ws://127.0.0.1:8545")
      --db-connection-limit uint        maximum number of open database connections (default 16)
      --db-idle-connection-limit uint   maximum number of idle database connections (default 4)
      --write-interval duration         interval between checks for job writing (default 1s)
      --pending-limit uint              maximum number of pending jobs per combination (dafault 1000)
      --height-range uint               maximum heights to include in a single job (default 10)
```

> ⚠️ Be careful when changing the batch amount, as it can cause the job-watcher to crash.
> The recommended value of 200 is set by default in order to prevent job-watcher crashes.

## Database Address - Data Source Name

Data Source Name (DSN) is the string specified describing how the connection to the database should be established.
Format of the string is the following:

```
host=localhost user=database-user password=password dbname=database-name port=5432 sslmode=disable
```
