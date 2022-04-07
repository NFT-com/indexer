# Parsing Dispatcher

Parsing Dispatcher consumes messages from the queue and launches jobs.

## Usage

```
Usage of parsing-dispatcher:
  -a, --api string              jobs api base endpoint
  -t, --tag string              rmq producer tag (default "jobs-watcher")
  -n, --network string          redis network type (default "tcp")
  -u, --url string              redis server connection url
  --database int                redis database number (default 1)
  --delivery-queue string       delivery queue name (default "discovery")
  --parsing-queue string        parsing queue name (default "parsing")
  -l, --log-level string        log level (default "info")
```
