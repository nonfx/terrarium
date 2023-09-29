data "aws_availability_zones" "available" {}
data "aws_region" "current" {}

locals {
  name   = "${var.common_name_prefix}-${var.environment}-demo"
  region = data.aws_region.current.name
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
  source  = "terraform-aws-modules/ec2-instance/aws"
  name    = "single-instance"
  ami_id  = var.ami_id

  instance_type          = var.instance_type
  key_name               = var.key_name
  monitoring             = true
  vpc_security_group_ids = [module.ec2_sg.security_group_id]
  subnet_id              = element(module.vpc.public_subnets, 0)

  tags = local.tags
}

################################################################################
# EC2 Security Group Module
################################################################################
module "ec2_sg" {
  source = "github.com/cldcvr/cldcvr-xa-terraform-aws-security-group?ref=v0.1.0"
  vpc_id = module.vpc.vpc_id
  name   = "${local.name}-ec2-sg"

  ingress_rules = [
    {
      from_port = 22
      to_port   = 22
      protocol  = "tcp"
      cidr_blocks = ["0.0.0.0/0"]
    }
  ]

  egress_rules = [
    {
      from_port = 0
      to_port   = 0
      protocol  = "-1"
      cidr_blocks = ["0.0.0.0/0"]
    }
  ]
}

################################################################################
# Database Module
################################################################################
module "aws_sql_database" {
  source       = "./modules/db"  # Assuming you have a db module locally
  instance_class = "db.t2.micro"
  username     = "admin"
  password     = "securepassword"
  name         = "mydatabase"
  engine       = "mysql"
  engine_version = "5.7"
  allocated_storage = 20

  vpc_security_group_ids = [module.db_sg.security_group_id]
  subnet_ids            = module.vpc.private_subnets
}

module "db_sg" {
  source = "./modules/security_group"  # Assuming you have a security_group module locally
  vpc_id = module.vpc.vpc_id
  name   = "${local.name}-db-sg"

  ingress_rules = [
    {
      from_port = 3306
      to_port   = 3306
      protocol  = "tcp"
      cidr_blocks = [local.vpc_cidr]
    }
  ]
}

################################################################################
# Redis Module
################################################################################
module "redis" {
  source       = "./modules/redis"  # Assuming you have a redis module locally
  engine       = "redis"
  engine_version = "6.x"
  node_type    = "cache.t2.micro"
  port         = 6379

  vpc_security_group_ids = [module.redis_sg.security_group_id]
  subnet_ids            = module.vpc.private_subnets
}

module "redis_sg" {
  source = "./modules/security_group"  # Assuming you have a security_group module locally
  vpc_id = module.vpc.vpc_id
  name   = "${local.name}-redis-sg"

  ingress_rules = [
    {
      from_port = 6379
      to_port   = 6379
      protocol  = "tcp"
      cidr_blocks = [local.vpc_cidr]
    }
  ]
}
