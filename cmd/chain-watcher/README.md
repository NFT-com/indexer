# Chain Watcher

Chain Watcher watches the chain and instantiates all the parsing jobs required for the network.
If the chain watcher stop in the middle of the instantiation, it will retrieve the last job saved in the API and start from that height instead of 0.

## Usage

```
Usage of jobs-api:
  -a, --api string              jobs api base endpoint
  -i, --chain-id string         id of the chain to watch
  -u, --chain-url string        url of the chain to connect
  -t, --chain-type string       type of chain to parse
  -c, --contract string         chain contract to watch
  -e, --event string            chain event type to watch
  -l, --log-level string        log level (default "info")
  --standard-type string        standard type of the contract to watch
```
