# Chain Watcher

Chain Watcher watches the chain and push jobs to parse network data.
If the chain watcher stopped during an instantiation, upon restarting it retrieves the last job saved in the API and starts from that height instead of 0.

## Usage

```
Usage of chain-watcher:
  -a, --api string              jobs api base endpoint
  -b, --batch int64             number of jobs in each batch request (default 200)
  --batch-delay durantion       delay between each batch request (default 1s)
  -i, --chain-id string         id of the chain to watch
  -u, --chain-url string        url of the chain to connect
  -t, --chain-type string       type of chain to parse
  -d, --db string               data source name for database connection
  -l, --log-level string        log level (default "info")
  --standard-type string        standard type of the contract to watch
```

> ⚠️ Be careful when changing the batch amount, as it can cause the job-watcher to crash.
> The recommended value of 200 is set by default in order to prevent job-watcher crashes.

## Database Address - Data Source Name

Data Source Name (DSN) is the string specified describing how the connection to the database should be established.
Format of the string is the following:

```
host=localhost user=database-user password=password dbname=database-name port=5432 sslmode=disable
```
