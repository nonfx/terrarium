provider "random" {}

data "aws_availability_zones" "available" {}
data "aws_region" "current" {}
data "aws_caller_identity" "current" {}

resource "random_string" "random" {
  length  = 8
  special = false
  upper   = false
  lower   = true
}

locals {
  azs = ["eu-west-1a", "eu-west-1b", "eu-west-1c"]

  database_enabled = anytrue([
    (length(local.tr_component_postgres) > 0)
  ])

  elasticache_enabled = anytrue([
    (length(local.tr_component_redis) > 0)
  ])
}

module "core_vpc" {
  source = "registry.terraform.io/terraform-aws-modules/vpc/aws"

  name = random_string.random.id
  cidr = "10.0.0.0/16"

  azs              = local.azs
  private_subnets  = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
  public_subnets   = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]
  database_subnets = local.database_enabled ? ["10.0.21.0/24", "10.0.22.0/24", "10.0.23.0/24"] : []
  elasticache_subnets = local.elasticache_enabled ? ["10.0.31.0/24", "10.0.32.0/24", "10.0.33.0/24"] : []

  enable_nat_gateway = true
  enable_vpn_gateway = true
}
