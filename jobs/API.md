# Jobs API

Jobs API is a REST API serving the data related to discovery and parsing blockchain data.

## Endpoints

There is two different job types `discovery` and `parsing` so there is a need for two endpoints.

- `/discoveries/`
- `/parsers/`

## Methods

- `POST`
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
	ChainURL      string   // Chain URL to connect to.
	ChainType     string   // Web3 compatible, Flow, etc...
	Block         string   // Block to run discovery.
	Addresses     []string // Addressed to filter in the discovery, empty for no filter
	InterfaceType string   // Interface type to filter for.
}

```

### Parsing

```go
package job

type Parsing struct {
	ChainURL      string // Chain URL to connect to.
	ChainType     string // Web3 compatible, Flow, etc...
	InterfaceType string // Interface type/id to filter for.
	Block         string // Block to parse.
	Address       string // Address of the contract to parse.
	EventType     string // Event type to parse.
}

```
