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

## Local k8s setup

### Setting up server in minikube k8s cluster

- Create env file
```bash
cp app.example.env app.env
```
- Start minikube k8s cluster
```bash
minikube start
```
- Set minikube docker env
```bash
eval $(minikube docker-env)
```
- Build simplebank docker image
```bash
docker build . -t aseerkt/simplebank:latest
```
- Create *simplebank** namepsace
```bash
kubectl create namespace simplebank
```
- Apply simplebank k8s deployment objects
```bash
kubectl apply -f eks/deployment.yml
```

### Post cleanup 

- 
```bash
kubectl delete -f eks/deployment.yml
```
- 
```bash
docker rmi aseerkt/simplebank:latest
```
- 
```bash
eval $(minikube docker-env --unset)
```
- 
```bash
minikube stop && minikube delete
```

# Deployment


- Create IAM user
- Create EKS Cluster using [eksctl](https://docs.aws.amazon.com/eks/latest/userguide/getting-started-eksctl.html)