output "vpc_id" {
  description = "VPC ID"
  value       = module.vpc.vpc_id
}

# EKS Cluster

output "cluster_name" {
  value = module.eks.cluster_name
}

output "cluster_endpoint" {
  value = module.eks.cluster_endpoint
}

output "cluster_ca_certificate" {
  value     = module.eks.cluster_certificate_authority_data
  sensitive = true
}

output "oidc_provider_arn" {
  description = "IAM OpenID Connector Provider ARN"
  value       = module.eks.oidc_provider_arn
}

output "oidc_provider_url" {
  description = "IAM OpenID Connector Provider URL"
  value       = module.eks.oidc_provider
}

# ECR

output "ecr_repository" {
  value     = aws_ecr_repository.this.repository_url
  sensitive = true
}

# RDS

output "db_username" {
  value = module.db.db_instance_username
}

output "db_endpoint" {
  value     = module.db.db_instance_endpoint
  sensitive = true
}

output "db_port" {
  value = module.db.db_instance_port
}

output "db_password" {
  value     = random_password.rds_password.bcrypt_hash
  sensitive = true
}
