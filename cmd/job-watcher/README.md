# Job Watcher

The job watcher watches the PostgreSQL database for new jobs from the job creator and pushes them into their respective queue.

## Usage

```
Usage of job-watcher:
      --db-connection-limit uint        maximum number of open database connections (default 16)
      --db-idle-connection-limit uint   maximum number of idle database connections (default 4)
  -j, --job-database string             postgresql connection details for job database (default "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")
  -l, --log-level string                severity level for log output (default "info")
      --read-interval duration          interval between checks for job reading (default 100ms)
  -d, --redis-database int              redis database number (default 1)
  -u, --redis-url string                redis server url (default "127.0.0.1:6379")
```

## Database Address — Data Source Name

Data Source Name (DSN) is the string specified describing how the connection to the database should be established.
Format of the string is the following:

```
host=localhost user=database-user password=password dbname=database-name port=5432 sslmode=disable
```