# Jobs Watcher

Job Watcher watches the dispatcher and parsing websockets for new updates and pushes them into their respective queue.

## Usage

```
Usage of jobs-watcher:
  --addition-queue string       addition queue name (default "addition")
  -a, --api string              jobs api base endpoint
  -t, --tag string              rmq producer tag (default "jobs-watcher")
  -n, --network string          redis network type (default "tcp")
  -u, --url string              redis server connection url
  --database int                redis database number (default 1)
  --delivery-queue string       discovery queue name (default "discovery")
  --parsing-queue string        parsing queue name (default "parsing")
  -l, --log-level string        log level (default "info")
```
