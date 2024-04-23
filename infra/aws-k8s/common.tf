data "aws_availability_zones" "available" {}

data "aws_region" "current" {}


locals {
  name     = "simplebank"
  vpc_cidr = "192.168.0.0/16"
  azs      = slice(data.aws_availability_zones.available.names, 0, 3)

  albc_version = "v2.7.2"
  albc_sa_name = "aws-load-balancer-controller"

  vpc_id = module.vpc.vpc_id

  cluster_name             = module.eks.cluster_name
  cluster_endpoint         = module.eks.cluster_endpoint
  cluster_ca_certificate   = module.eks.cluster_certificate_authority_data
  cluster_oidc_provder_arn = module.eks.oidc_provider_arn
  cluster_oidc_provder_url = module.eks.oidc_provider
}


