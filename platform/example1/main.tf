provider "random" {}

data "aws_availability_zones" "available" {}
data "aws_region" "current" {}
data "aws_caller_identity" "current" {}

# inputs - from operations team
variable "network_vpc_cidr" {
  default = "10.0.0.0/16"
  type = string
}

variable "network_number_of_azs" {
  default = 2
  type = number
  sensitive = true
  nullable = true
}


locals {
  name             = "demo-${random_string.random.id}"
  vpc_cidr         = var.tr_vpc_cidr
  azs              = slice(data.aws_availability_zones.available.names, 0, var.number_of_azs)

  tags = {
    Name = local.name
  }
}


module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "~> 3.0"

  name = local.name
  cidr = local.vpc_cidr

  azs              = local.azs
  public_subnets   = [for k, v in local.azs : cidrsubnet(local.vpc_cidr, 8, k)]
  private_subnets  = [for k, v in local.azs : cidrsubnet(local.vpc_cidr, 8, k + 3)]
  database_subnets = local.tr_if_database ? [for k, v in local.azs : cidrsubnet(local.vpc_cidr, 8, k + 6)] : []
  elasticache_subnets = local.tr_if_cache ? [] : []

  create_database_subnet_group  = true
  manage_default_security_group = true

  enable_flow_log                          = true
  enable_dns_hostnames                     = true
  enable_dns_support                       = true
  flow_log_destination_type                = "cloud-watch-logs"
  create_flow_log_cloudwatch_log_group     = true
  create_flow_log_cloudwatch_iam_role      = true
  flow_log_cloudwatch_log_group_kms_key_id = module.cloudwatch_kms_key.aws_kms_key_arn

  tags = local.tags
}

module "vpc_endpoints" {
  source = "terraform-aws-modules/vpc/aws//modules/vpc-endpoints"

  vpc_id             = module.vpc.vpc_id
  security_group_ids = [module.security_group.security_group_id]

  endpoints = {
    rds = {
      service = "rds"
      tags    = { Name = "${local.name}-endpoint" }
    },
  }
}

module "postgres_security_group" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "~> 4.0"

  name   = local.name
  vpc_id = module.vpc.vpc_id

  # ingress
  ingress_with_cidr_blocks = [
    {
      rule        = "postgresql-tcp"
      cidr_blocks = module.vpc.vpc_cidr_block
    }
  ]
  egress_with_cidr_blocks = [
    {
      rule        = "all-all"
      cidr_blocks = "0.0.0.0/0"
    },
  ]

  tags = local.tags
}
