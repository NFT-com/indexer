# NFT.com Indexer Infra 

![nftcom_arch-Page-5](https://user-images.githubusercontent.com/5006941/175179454-0936c204-1be9-4172-8460-41ecfcf2fdf2.png)

Our indexer infrastructure is deployed using Pulumi

## Indexer Infrastructure 

- CICD Pipeline with GitHub, GitHub Actions, Node/Typescript and Pulumi
- Multi-env: Dev, Staging, Prod (manual permission required for deployments to any env)
- Secrets managed in Doppler, flow into GitHub Secrets and used in GitHub actions (secrets —> env variables)

### GitHub Deployment Process 

- Pushed branches starting with `fix/` or `feat/` triggers a deployment to the dev environment (nftcom-indexer-dev)
- Merged branches starting with `fix/` or `feat/` into main triggers a deployment to the staging environment (nftcom-indexer-staging)
- Tagging the main branch starting with `rel` triggers deployment to the prod environment (nftcom-indexer-prod)

### Indexer AWS Infrastructure Components 

- Elastic Container Service (ECS) Cluster & Task Definitions
- ECS EC2 Capacity Provider (w/ASG & LaunchConfig)
- Elastic Container Registry (ECR)
- Lambda
- 3x Aurora Postgres RDS (Databases for Graph, Event and Jobs)
- ElastiCache Redis (not currently used)
- IAM Roles as Needed for EB, EC2

### Indexer Deployment Notes

- After deployment is triggered, github actions executes the script (action.yml) to build/zip the lambda workers, deploy the shared infra including the lambda workers, build the latest images and push to AWS ECR, and finally deploy the ECS cluster including the task definitions to instantiate the indexer tasks on ECS (nsqd, nsqlookupd, parserDispatcher, additionDispatcher) 
- The deployment only triggers updates to the infra, images and ECS task definitions but does not enable/start anything. When we are ready to automate 24/7 run of the indexer we will add an ECS service to keep the tasks up indefinitely. For now we want to manually control the runs and thus this approach works best for testing and running the historical sync. 

### Secrets / Environment Variables via Doppler

- AWS_ACCOUNT_ID
- AWS_REGION
- AWS_ACCESS_KEY
- AWS_SECRET_ACCESS_KEY
- EC2_PUBLIC_IP
- DB_EVENT_HOST
- DB_GRAPH_HOST
- DB_JOB_HOST
- DB_USER
- DB_NAME
- DB_PASSWORD
- DB_PORT
- REDIS_PORT
- ADDITION_RATE_LIMIT
- PARSER_HEIGHT_RANGE
- PARSER_RATE_LIMIT
- ZMOK_HTTP_URL
- ZMOK_WS_URL
  
