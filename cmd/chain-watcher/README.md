# Chain Watcher

Chain Watcher watches the chain and push jobs to parse network data.
If the chain watcher stopped during an instantiation, upon restarting it retrieves the last job saved in the API and starts from that height instead of 0.

## Usage

```
Usage of chain-watcher:
  -i, --chain-id string         id of the chain to watch
  -u, --chain-url string        url of the chain to connect
  -t, --chain-type string       type of chain to parse
  -d, --data-database string    data database connection string
  -j, --job-database string     jobs database connection string
  -l, --log-level string        log level (default "info")
  -s, --start-height uint64     default start height when no jobs found (default 0)
  -b, --job-limit uint          maximum number of pending jobs per combination (default 1000)
  -c, --notify-period duration  how often to notify watchers to create new jobs (default 1s)
```

## Database Address - Data Source Name

Data Source Name (DSN) is the string specified describing how the connection to the database should be established.
Format of the string is the following:

```
host=localhost user=database-user password=password dbname=database-name port=5432 sslmode=disable
```
