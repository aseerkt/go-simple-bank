# Simple Bank

## Prerequisites

- Go 1.22.2
- Docker Desktop
- [golang-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
- sqlc
- mockgen
- AWS CLI - v2
- kubectl
- eksctl
- minikube

## Getting Started

- If you don't already have `postgres:16-alpine` docker image;
```bash
# Pull docker image
docker pull postgres:16-alpine

# Run postgres docker container
make postgres

# Create database in postgres docker container 
make createdb

# Run migration
make migrateup
```

For more scripts checkout [`Makefile`](/Makefile)

# Deployment


- Create IAM user
- Create EKS Cluster using [eksctl](https://docs.aws.amazon.com/eks/latest/userguide/getting-started-eksctl.html)