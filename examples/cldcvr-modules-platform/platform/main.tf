data "aws_availability_zones" "available" {}
data "aws_region" "current" {}
provider "aws" {}
provider "random" {}
locals {
  name   = "${var.common_name_prefix}-${var.environment}-demo"
  region = data.aws_region.current.name

  vpc_cidr = var.vpc_cidr
  azs      = slice(data.aws_availability_zones.available.names, 0, 2)

  tags = var.common_tags

}

################################################################################
# VPC Module
################################################################################

module "vpc" {
  source = "github.com/cldcvr/cldcvr-xa-terraform-aws-vpc?ref=v0.1.0"

  name = "${local.name}-vpc"
  cidr = local.vpc_cidr

  azs = local.azs
  # This code uses the `cidrsubnet` function to generate subnet ranges for the `private_subnets` and `public_subnets` variables. The `for` loop iterates over the `local.azs` map to generate a subnet range for each availability zone. The `cidrsubnet` function takes the `local.vpc_cidr` as the base CIDR block, `4` as the prefix length, and `k + coalesce(length(local.azs), 2)` as the subnet number for `private_subnets`, and `k + coalesce(length(local.azs), 2) + 4` as the subnet number for `public_subnets`. The `+4` is added to the subnet number for `public_subnets` to ensure that the subnets created for `private_subnets` and `public_subnets` do not overlap and do not result in duplicate subnets.
  private_subnets    = [for k, v in local.azs : cidrsubnet(local.vpc_cidr, 4, k + coalesce(length(local.azs), 2))]
  public_subnets     = [for k, v in local.azs : cidrsubnet(local.vpc_cidr, 4, k + coalesce(length(local.azs), 2) + 4)]
  tags               = local.tags
  enable_nat_gateway = true
  single_nat_gateway = true
}

################################################################################
# EC2 Instance Module
################################################################################
#EC2 Module from public Terraform Registry is the compute for this platform

module "tr_component_service_web" {
  source = "terraform-aws-modules/ec2-instance/aws"

  for_each = {
    for k, v in local.tr_component_service_web : k => merge(v, var.ec2_config["default"])
  }

  name                        = "${local.name}-ec2-instance"
  ami                         = each.value.ami
  instance_type               = each.value.instance_type
  monitoring                  = each.value.monitoring_enabled
  create_iam_instance_profile = true
  #Below code sets the associate_public_ip_address attribute to true if the length of the local.tr_component_service_web map is greater than 0, and false otherwise. The associate_public_ip_address attribute is used to associate a public IP address with the instances launched in the subnet. By setting this attribute to true, the instances launched in the subnet will have a public IP address associated with them, which allows them to communicate with the internet. By setting this attribute to false, the instances launched in the subnet will not have a public IP address associated with them, which means they will not be able to communicate with the internet directly.
  associate_public_ip_address = length(local.tr_component_service_web) > 0 ? true : false
  subnet_id                   = module.vpc.public_subnets[0]
  vpc_security_group_ids      = [module.ec2_sg[each.key].id]
  tags                        = local.tags
}

module "ec2_sg" {
  source = "github.com/cldcvr/cldcvr-xa-terraform-aws-security-group?ref=v0.1.0"
  for_each = {
    for k, v in local.tr_component_service_web : k => merge(v, var.ec2_config["default"])
  }
  vpc_id = module.vpc.vpc_id
  security_group = {
    name = "${local.name}-${each.key}-ec2-sg",
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
  env      = var.environment
  app_name = "${local.name}-${each.key}-ec2-sg"
}

################################################################################
#DB Instance
################################################################################
#We combine the values coming from developer with values set by DevOps
locals {
  database_configuration = {
    for k, v in local.tr_component_postgres : k => merge(v, var.db_config["default"])
  }
}

resource "random_password" "rds_password" {
  length  = 16
  special = false
}

module "tr_component_postgres" {
  source = "github.com/cldcvr/cldcvr-xa-terraform-aws-db-instance?ref=v0.1.0"

  for_each = local.database_configuration

  instance_class = each.value.instance_class
  username       = try(each.value.username, "postgres")
  password       = try(each.value.password, random_password.rds_password.result)
  ## This code sets the `name` attribute of the `module "aws_sql_database"` block to the value of `each.value.name` if it exists, or `${local.name}-${each.key}-db-instance` otherwise. The `try` function is used to check if `each.value.name` exists. If it does, it is returned by the `try` function. Otherwise, the fallback name `${local.name}-${each.key}-db-instance` is returned. The `replace` function is used to remove any non-alphanumeric characters from the generated name. The regular expression `[^a-zA-Z0-9]` matches any character that is not a letter or a number. This ensures that the generated name contains only alphanumeric characters.
  name                        = replace(try(each.value.name, "${local.name}-${each.key}-db-instance"), "-", "")
  storage_type                = try(each.value.storage_type, "gp2")
  engine                      = try(each.value.engine, "postgres")
  engine_version              = try(each.value.engine_version, "14")
  allocated_storage           = try(each.value.allocated_storage, "100")
  port                        = try(each.value.port, 5432)
  retention_period            = try(each.value.retention_period, 20)
  family                      = try(format("postgres%s", element(split(".", each.value.engine_version), 0)), "postgres13")
  max_allocated_storage       = try(each.value.max_allocated_storage, 200)
  iops                        = try(each.value.iops, null)
  backup_window               = try(each.value.backup_window, "03:00-06:00")
  maintenance_window          = try(each.value.maintenance_window, "Mon:00:00-Mon:03:00")
  identifier                  = try(each.value.identifier, var.common_name_prefix)
  multi_az                    = try(each.value.multi_az, false)
  storage_encrypted           = try(each.value.storage_encrypted, true)
  delete_automated_backups    = try(each.value.delete_automated_backups, true)
  allow_major_version_upgrade = try(each.value.allow_major_version_upgrade, true)
  auto_minor_version_upgrade  = try(each.value.auto_minor_version_upgrade, true)
  apply_immediately           = try(each.value.apply_immediately, true)
  skip_final_snapshot         = try(each.value.skip_final_snapshot, true)
  deletion_protection         = try(each.value.deletion_protection, false)
  publicly_accessible         = try(each.value.publicly_accessible, false)
  vpc_id                      = module.vpc.vpc_id
  subnet_ids                  = module.vpc.private_subnets
  vpc_security_group_ids      = [module.db_sg[each.key].id]

  depends_on = [module.db_sg]

}

module "db_sg" {
  source = "github.com/cldcvr/cldcvr-xa-terraform-aws-security-group?ref=v0.1.0"

  for_each = local.database_configuration

  vpc_id = module.vpc.vpc_id
  security_group = {
    name = "${local.name}-${each.key}-db-sg",
    ingress = [
      {
        type        = "ingress"
        cidr_blocks = [local.vpc_cidr]
        description = "Allows traffic on Port 5432"
        from_port   = try(each.value.port, 5432)
        to_port     = try(each.value.port, 5432)
        protocol    = "tcp"
      }
    ],
    egress = [
      {
        type        = "egress"
        cidr_blocks = [local.vpc_cidr]
        description = "Allow Outbound traffic on Port ${try(each.value.port, 6379)}"
        from_port   = try(each.value.port, 6379)
        to_port     = try(each.value.port, 6379)
        protocol    = "tcp"
      }
    ]
  }
  env      = var.environment
  app_name = "${local.name}-${each.key}-db-sg"
}

################################################################################
#Redis Instance
################################################################################
#We combine the values coming from developer with values set by DevOps

locals {
  redis_configuration = {
    for k, v in local.tr_component_redis : k => merge(v, var.redis_config["default"])
  }
}

module "tr_component_redis" {
  source = "github.com/cldcvr/cldcvr-xa-terraform-aws-elasticache-cluster?ref=v0.1.0"

  for_each = local.redis_configuration

  name       = "${local.name}-${each.key}-redis"
  subnet_ids = module.vpc.private_subnets
  #This is inserted as module creates subnet group based on this value
  create_subnet_group  = length(module.vpc.private_subnets) > 0 ? true : false
  security_group_ids   = [module.redis_sg[each.key].id]
  engine               = try(each.value.engine, "redis")
  engine_version       = try(each.value.engine_version, "5.0.6")
  port                 = try(each.value.port, 6379)
  parameter_group_name = try(each.value.parameter_group, "default.redis5.0")
  extra_tags           = local.tags

}


module "redis_sg" {
  source   = "github.com/cldcvr/cldcvr-xa-terraform-aws-security-group?ref=v0.1.0"
  for_each = local.redis_configuration
  vpc_id   = module.vpc.vpc_id
  security_group = {
    name = "${local.name}-${each.key}-redis-sg",
    ingress = [
      {
        type        = "ingress"
        cidr_blocks = [local.vpc_cidr]
        description = "Allows traffic on Port ${try(each.value.port, 6379)}"
        from_port   = try(each.value.port, 6379)
        to_port     = try(each.value.port, 6379)
        protocol    = "tcp"
      }
    ],
    egress = [
      {
        type        = "egress"
        cidr_blocks = [local.vpc_cidr]
        description = "Allow Outbound traffic on Port ${try(each.value.port, 6379)}"
        from_port   = try(each.value.port, 6379)
        to_port     = try(each.value.port, 6379)
        protocol    = "tcp"
      }
    ]
  }
  env      = var.environment
  app_name = "${local.name}-${each.key}-ec2-sg"
}
