# Jobs API

Jobs API is a REST API serving the data related to discovery and parsing blockchain data.

## Endpoints

There is two different job types `discovery` and `parsing` so there is a need for two endpoints.

- `/discoveries/`
- `/parsers/`

## Methods

- `PUT`
    - Allows creating new jobs.
- `GET`
    - Lists all the jobs.
- `PATCH`
    - Allows to cancel a job.

## Structures

### Discovery

```go
package job

type Discovery struct {
	ChainURL      string
	ChainType     string
	StartBlock    string
	EndBlock      string
	Addresses     []string
	InterfaceType string
}

```

### Parsing

```go
package job

type Parsing struct {
	ChainURL        string
	ChainType       string
	InterfaceType   string
	Block           string
	TransactionHash string
	Address         string
	Topic           string
	IndexedData     []string
	Data            []byte
}

```