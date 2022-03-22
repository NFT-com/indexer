# NFT.com Indexer Infra 

![image](https://user-images.githubusercontent.com/5006941/159491381-a056fb1e-11c9-4fa7-9365-5adf560b252c.png)


Our indexer infrastructure is deployed using Pulumi

## Indexer Infrastructure 
- CICD Pipeline with GitHub, GitHub Actions, Node/Typescript and Pulumi
- Multi-env: Dev, Staging, Prod (special security settings for Prod)
- Secrets managed in Doppler, flow into GitHub Secrets and used in GitHub actions (secrets —> env variables)

### GitHub Deployment Process 
- Branches starting with ‘fix/’ or ‘feat/’ and pushed will trigger a deployment to the dev environment (nftcom-indexer-dev)
- Merge branch starting with ‘fix/’ or ‘feat/ into main will trigger a deployment to the staging environment (nftcom-indexer-staging)
- Tagging the main branch starting with ‘release’ triggers deployment to the prod environment (nftcom-indexer-prod)

### Indexer AWS Infrastructure Components 
- Elastic Beanstalk
- Elastic Container Registry
- S3
- Aurora Postgres RDS
- ElastiCache Redis
- IAM Roles as Needed for EB, EC2

### Permissions
- Lambda role (to use in SAM template for functions)
    - ARN: arn:aws:iam::016437323894:role/AWSLambdaBasicExecutionRole
    - Basic permissions for pushing logs to CloudWatch Logs
- EC2 Role
    - arn:aws:iam::016437323894:instance-profile/dev-indexer-eb-ec2-profile-role
    - Includes permissions to execute a SAM template incl permissions for CloudFormation, Lambda, IAM, etc. Defined in the Pulumi code.

### Secrets / Environment Variables via Doppler
- DB_PASSWORD = <hidden>
- DB_PORT = 10030
- REDIS_PORT = 10020
- AWS_ACCOUNT_ID = 016437323894
- AWS_REGION = us-east-1
- AWS_ACCESS_KEY & AWS_SECRET_ACCESS_KEY = <hidden, used for CICD deployment>

### ETH Node Details 
