# Chain Watcher

Chain Watcher watches the chain and push jobs to parse network data.
If the chain watcher stopped during an instantiation, upon restarting it retrieves the last job saved in the API and starts from that height instead of 0.

## Usage

```
Usage of chain-watcher:
  -a, --api string              jobs api base endpoint
  -i, --chain-id string         id of the chain to watch
  -u, --chain-url string        url of the chain to connect
  -t, --chain-type string       type of chain to parse
  -c, --contract string         chain contract to watch
  -e, --event string            chain event type to watch
  -l, --log-level string        log level (default "info")
  --standard-type string        standard type of the contract to watch
```