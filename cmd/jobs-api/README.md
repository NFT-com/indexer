# Jobs API

## Usage

```
Usage of jobs-api:
  -b, --bind string              jobs api binding port (default ":8081")
  -d, --database string          data source name for database connection
  -l, --log-level string         log level (default "info")
```

## Database Address - Data Source Name

Data Source Name (DSN) is the string specified describing how the connection to the database should be established.
Format of the string is the following:

```
host=localhost user=database-user password=password dbname=database-name port=5432 sslmode=disable
```
