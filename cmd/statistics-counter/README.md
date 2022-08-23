# Statistics Counter

The statistics counter counts the logs of a specific contract and event hash.

## Command Line Parameters

```
Usage of action-dispatcher:
  -l, --log-level string                severity level for log output (default "info")
  
  -n, --node string                     ethereum node url
  
  -s, --starting-height uint64          counter starting block
  -a, --addresses []string              addresses to count
  -h, --hashes []string                 event hashes to count
  
  -n, --lambda-name string              name of Lambda function for invocation (default "action-worker")

  -r, --rate-limit uint                 maximum number of API requests per second (default 100)
```
