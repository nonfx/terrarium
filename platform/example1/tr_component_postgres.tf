# postgres - dynamic component pattern

# Inputs - from dependency interface or app developer
locals {
  tr_component_postgres_id      = ""
  tr_component_postgres_version = "11.12"
  tr_component_postgres_db_name = null
}

# inputs - from operations team
variable "tr_component_postgres_instance_class" {
  default = "db.t3.micro"
  type = string
}

variable "tr_component_postgres_storage" {
  default = 2
  type = number
}

# Module call
module "tr_component_postgres" {
  source   = "terraform-aws-modules/rds/aws"

  identifier             = local.tr_component_postgres_id
  engine                 = "postgres"
  engine_version         = local.tr_component_postgres_version
  db_name                = local.tr_component_postgres_db_name
  subnet_ids             = module.vpc.database_subnets
  instance_class         = var.tr_component_postgres_instance_class
  allocated_storage      = var.tr_component_postgres_storage
  vpc_security_group_ids = [ module.postgres_security_group.security_group_id ]

  depends_on = [ module.vpc_endpoints2 ]
}

# Outputs
output "tr_component_postgres_host" {
  value = module.tr_component_postgres.db_instance_address
}

output "tr_component_postgres_port" {
  value = module.tr_component_postgres.db_instance_port
}

output "tr_component_postgres_username" {
  value = module.tr_component_postgres.db_instance_username
  sensitive = true
}

output "tr_component_postgres_password" {
  value = module.tr_component_postgres.db_instance_password
  sensitive = true
}
