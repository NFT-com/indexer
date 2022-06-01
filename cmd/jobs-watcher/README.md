# Jobs Watcher

The jobs watcher watches the PostgreSQL database for new jobs from the jobs creator and pushes them into their respective queue.

## Usage

```
Usage of jobs-watcher:
  -l, --log-level string                severity level for log output (default "info")

  -j, --jobs-database string            Postgres connection details for jobs database (default "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=jobs sslmode=disable")
  -n, --nsq-address string              NSQ server address (default "127.0.0.1:4150")

      --db-connection-limit uint        maximum number of open database connections (default 16)
      --db-idle-connection-limit uint   maximum number of idle database connections (default 4)
      --read-interval duration          interval between checks for job reading (default 100ms)
```

## Database Address — Data Source Name

The database addresses are given in DSN (Data Source Name) format, which is a string that describes the parameters of the connection to the database.
Format of the string is the following:

```
host=localhost user=database-user password=password dbname=database-name port=5432 sslmode=disable
```