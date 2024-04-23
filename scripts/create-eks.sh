#!/bin/bash

CLUSTER_NAME=simplebank

eksctl create cluster \
  --name "$CLUSTER_NAME" \
  --nodes 3 \
  --nodes-min 1 \
  --nodes-max 4

# Installing ack for rds using helm chart
# https://aws-controllers-k8s.github.io/community/docs/tutorials/rds-example/#install-the-ack-service-controller-for-rds

aws ecr-public get-login-password --region us-east-1 | helm registry login --username AWS --password-stdin public.ecr.aws

helm install --create-namespace -n ack-system oci://public.ecr.aws/aws-controllers-k8s/rds-chart --version=0.0.27 --generate-name --set=aws.region=us-east-1

# Install AWS load balancer controller

helm repo add eks https://aws.github.io/eks-charts

wget https://raw.githubusercontent.com/aws/eks-charts/master/stable/aws-load-balancer-controller/crds/crds.yaml

kubectl apply -f crds.yaml

helm install aws-load-balancer-controller eks/aws-load-balancer-controller -n kube-system \ 
--set clusterName=$CLUSTER_NAME \ 
--set serviceAccount.create=false \
--set serviceAccount.name=aws-load-balancer-controller

# Setup database networking

RDS_SUBNET_GROUP_NAME="<your subnet group name>"
RDS_SUBNET_GROUP_DESCRIPTION="<your subnet group description>"
EKS_VPC_ID=$(aws eks describe-cluster --name="${EKS_CLUSTER_NAME}" \
  --query "cluster.resourcesVpcConfig.vpcId" \
  --output text)
EKS_SUBNET_IDS=$(
  aws ec2 describe-subnets \
    --filters "Name=vpc-id,Values=${EKS_VPC_ID}" \
    --query 'Subnets[*].SubnetId' \
    --output text
)

cat <<-EOF >db-subnet-groups.yaml
apiVersion: rds.services.k8s.aws/v1alpha1
kind: DBSubnetGroup
metadata:
  name: ${RDS_SUBNET_GROUP_NAME}
  namespace: ${APP_NAMESPACE}
spec:
  name: ${RDS_SUBNET_GROUP_NAME}
  description: ${RDS_SUBNET_GROUP_DESCRIPTION}
  subnetIDs:
$(printf "    - %s\n" ${EKS_SUBNET_IDS})
  tags: []
EOF

kubectl apply -f db-subnet-groups.yaml

# Security Groups for Pods to access RDS database

RDS_SECURITY_GROUP_NAME="<your security group name>"
RDS_SECURITY_GROUP_DESCRIPTION="<your security group description>"

EKS_CIDR_RANGE=$(
  aws ec2 describe-vpcs \
    --vpc-ids $EKS_VPC_ID \
    --query "Vpcs[].CidrBlock" \
    --output text
)

RDS_SECURITY_GROUP_ID=$(
  aws ec2 create-security-group \
    --group-name "${RDS_SUBNET_GROUP_NAME}" \
    --description "${RDS_SUBNET_GROUP_DESCRIPTION}" \
    --vpc-id "${EKS_VPC_ID}" \
    --output text
)

aws ec2 authorize-security-group-ingress \
  --group-id "${RDS_SECURITY_GROUP_ID}" \
  --protocol tcp \
  --port 5432 \
  --cidr "${EKS_CIDR_RANGE}"

# Create secrets that store DB credentials

RDS_DB_USERNAME="<your username>"
RDS_DB_PASSWORD="<your secure password>"

kubectl create secret generic -n "${APP_NAMESPACE}" jira-postgres-creds \
  --from-literal=username="${RDS_DB_USERNAME}" \
  --from-literal=password="${RDS_DB_PASSWORD}"

# Create RDS instance

RDS_DB_INSTANCE_NAME="jira-db"
RDS_DB_INSTANCE_CLASS="db.m6g.large"
RDS_DB_STORAGE_SIZE=100

cat <<-EOF >jira-db.yaml
apiVersion: rds.services.k8s.aws/v1alpha1
kind: DBInstance
metadata:
  name: ${RDS_DB_INSTANCE_NAME}
  namespace: ${APP_NAMESPACE}
spec:
  allocatedStorage: ${RDS_DB_STORAGE_SIZE}
  autoMinorVersionUpgrade: true
  backupRetentionPeriod: 7
  dbInstanceClass: ${RDS_DB_INSTANCE_CLASS}
  dbInstanceIdentifier: ${RDS_DB_INSTANCE_NAME}
  dbName: jira
  dbSubnetGroupName: ${RDS_SUBNET_GROUP_NAME}
  enablePerformanceInsights: true
  engine: postgres
  engineVersion: "13"
  masterUsername: ${RDS_DB_USERNAME}
  masterUserPassword:
    namespace: ${APP_NAMESPACE}
    name: jira-postgres-creds
    key: password
  multiAZ: true
  publiclyAccessible: false
  storageEncrypted: true
  storageType: gp2
  vpcSecurityGroupIDs:
    - ${RDS_SECURITY_GROUP_ID}
EOF

kubectl apply -f jira-db.yaml
