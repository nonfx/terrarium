data "aws_availability_zones" "available" {}
data "aws_region" "current" {}

locals {
  name     = "${var.common_name_prefix}-${var.environment}-demo"
  region   = data.aws_region.current.name
  vpc_cidr = var.vpc_cidr
  azs      = slice(data.aws_availability_zones.available.names, 0, 2)
  tags     = var.common_tags
}

################################################################################
# VPC Module
################################################################################
module "vpc" {
  source = "github.com/cldcvr/cldcvr-xa-terraform-aws-vpc?ref=v0.1.0"

  name = local.name
  cidr = local.vpc_cidr

  azs             = local.azs

  # This code uses the `cidrsubnet` function to generate subnet ranges for the `private_subnets` and `public_subnets` variables. The `for` loop iterates over the `local.azs` map to generate a subnet range for each availability zone. The `cidrsubnet` function takes the `local.vpc_cidr` as the base CIDR block, `4` as the prefix length, and `k + coalesce(length(local.azs), 2)` as the subnet number for `private_subnets`, and `k + coalesce(length(local.azs), 2) + 4` as the subnet number for `public_subnets`. The `+4` is added to the subnet number for `public_subnets` to ensure that the subnets created for `private_subnets` and `public_subnets` do not overlap and do not result in duplicate subnets.
  private_subnets = [for k, v in local.azs : cidrsubnet(local.vpc_cidr, 4, k)]
  public_subnets  = [for k, v in local.azs : cidrsubnet(local.vpc_cidr, 4, k)]

  tags               = local.tags
  enable_nat_gateway = true
  single_nat_gateway = true
}

################################################################################
# EC2 Instance Module
################################################################################
module "ec2_instance" {
  source = "terraform-aws-modules/ec2-instance/aws"
  name   = "single-instance"
  ami = var.ami_id

  instance_type          = var.instance_type
  key_name               = var.key_name
  monitoring             = true
  vpc_security_group_ids = [module.ec2_sg.id]
  subnet_id              = element(module.vpc.public_subnets, 0)

  tags = local.tags
}

################################################################################
# EC2 Security Group Module
################################################################################
module "ec2_sg" {
  source   = "github.com/cldcvr/cldcvr-xa-terraform-aws-security-group?ref=v0.1.0"
  vpc_id   = module.vpc.vpc_id
  app_name = "${local.name}-ec2-sg"
  env      = var.environment

  security_group = {
    name = "${local.name}-ec2-sg",
    ingress = [
      {
        type        = "ingress"
        cidr_blocks = ["0.0.0.0/0"]
        description = "Allows traffic on Port 80"
        from_port   = "80"
        to_port     = "80"
        protocol    = "tcp"
      },
      {
        type        = "ingress"
        cidr_blocks = [local.vpc_cidr]
        description = "Allows traffic on Port 5432"
        from_port   = "5432"
        to_port     = "5432"
        protocol    = "tcp"
      },
      {
        cidr_blocks = ["0.0.0.0/0"]
        type        = "ingress"
        description = "Allows traffic on Port 443"
        from_port   = "443"
        to_port     = "443"
        protocol    = "tcp"
      }
    ],
    egress = [
      {
        type        = "egress"
        cidr_blocks = ["0.0.0.0/0"]
        description = "Allow All Outbound"
        from_port   = "0"
        to_port     = "0"
        protocol    = "-1"
      }
    ]
  }
}



################################################################################
# Database Module
################################################################################
module "aws_sql_database" {
  source            = "github.com/cldcvr/cldcvr-xa-terraform-aws-db-instance?ref=v0.1.0"
  instance_class    = "db.t2.micro"
  username          = "admin"
  password          = "securepassword"
  name              = "mydatabase"
  engine            = "mysql"
  engine_version    = "5.7"
  allocated_storage = 20

  vpc_security_group_ids = [module.db_sg.id]
  subnet_ids             = module.vpc.private_subnets
}

module "db_sg" {
  source   = "github.com/cldcvr/cldcvr-xa-terraform-aws-security-group?ref=v0.1.0"
  vpc_id   = module.vpc.vpc_id
  app_name = "${local.name}-db-sg"
  env      = var.environment

  security_group = {
    name = "${local.name}-db-sg",
    ingress = [
      {
        type        = "ingress"
        cidr_blocks = [local.vpc_cidr]
        description = "Allows inbound traffic on Port 5432"
        from_port   = "5432"
        to_port     = "5432"
        protocol    = "tcp"
      }
    ],
    egress = [
      {
        type        = "egress"
        cidr_blocks = [local.vpc_cidr]
        description = "Allows outbound traffic on Port 5432"
        from_port   = "5432"
        to_port     = "5432"
        protocol    = "tcp"
      }
    ]
  }
}

################################################################################
# Redis Module
################################################################################
module "redis" {
  source         = "github.com/cldcvr/cldcvr-xa-terraform-aws-elasticache-cluster?ref=v0.1.0"
  name       = "${local.name}-redis"
  engine         = "redis"
  engine_version = "6.x"
  node_type      = "cache.t2.micro"
  port           = 6379

  security_group_ids = [module.redis_sg.id]
  subnet_ids             = module.vpc.private_subnets
}

module "redis_sg" {
  source   = "github.com/cldcvr/cldcvr-xa-terraform-aws-security-group?ref=v0.1.0"
  vpc_id   = module.vpc.vpc_id
  app_name = "${local.name}-redis-sg"
  env      = var.environment

  security_group = {
    name = "${local.name}-redis-sg",
    ingress = [
      {
        type        = "ingress"
        cidr_blocks = [local.vpc_cidr]
        description = "Allows inbound traffic on Port 6379"
        from_port   = "6379"
        to_port     = "6379"
        protocol    = "tcp"
      }
    ],
    egress = [
      {
        type        = "egress"
        cidr_blocks = [local.vpc_cidr]
        description = "Allows outbound traffic on Port 6379"
        from_port   = "6379"
        to_port     = "6379"
        protocol    = "tcp"
      }
    ]
  }
}
