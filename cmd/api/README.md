# API Server

This will be the place from where we build and RESTful API endpoints used internally and/or externally

Initial use case:
- addition worker - queue job to call an endpoint to query and populate NFT metadata

Other Potential use cases:
- reporting
- control panel

## Command Line Parameters

TODO - creds, etc

## Environment Variables

TODO

## Database Address — Data Source Name

The database addresses are given in DSN (Data Source Name) format, which is a string that describes the parameters of the connection to the database.
The format is the following:

```
host=localhost user=database-user password=password dbname=database-name port=5432 sslmode=disable
```
