# API Server

This will be the place from where we build and RESTful API endpoints used internally and/or externally

Initial use case:
- addition worker - queue job to call an endpoint to query and populate NFT metadata

Other Potential use cases:
- reporting
- control panel

## Command Line Parameters

```
Usage of action-dispatcher:
  -l, --log-level string                severity level for log output (default "info")
  -g, --graph-database string           Postgres connection details for graph database (default "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=graph sslmode=disable")
  -u, --username string                 Basic HTTP Authentication Username (default "admin")
  -p, --password string                 Basic HTTP Authentication Password (default "admin")
```

## Environment Variables

In dev/prod environments, set HTTP Basic Auth username/password in Doppler secrets

```
API_USERNAME=<username>
API_PASSWORD=<password>
```

## Database Address — Data Source Name

The database addresses are given in DSN (Data Source Name) format, which is a string that describes the parameters of the connection to the database.
The format is the following:

```
host=localhost user=database-user password=password dbname=database-name port=5432 sslmode=disable
```
