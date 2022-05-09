# Jobs Watcher

Job Watcher watches the dispatcher and parsing websockets for new updates and pushes them into their respective queue.

## Usage

```
Usage of jobs-watcher:
  -l, --log-level string                severity level for log output (default "info")
  -j, --jobs-database string            postgres connection details for jobs database (default "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=jobs sslmode=disable")
  -u, --redis-url string                url for redis server connection (default 1)
      --db-connection-limit uint        maximum number of open database connections (default 16)
      --db-idle-connection-limit uint   maximum number of idle database connections (default 4)
      --read-interval duration          interval between checks for job reading (default 200s)
```
