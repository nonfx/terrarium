# A relational database management system using SQL.
# @title: PostgreSQL Database
module "tr_component_postgres" {
  source = "./mock-modules/rds"

  for_each = local.tr_component_postgres

  identifier     = each.key
  engine_version = each.value.version
  db_name        = coalesce(each.value.db_name, each.key)
  engine         = "postgres"
  family         = format("postgres%s", each.value.version)

  instance_class = coalesce(lookup(var.db_instance_class, each.key, null), var.all_db_instance_class)

  vpc_security_group_ids = [module.postgres_security_group.security_group_id]
  subnet_ids             = module.core_vpc.database_subnets
}

module "postgres_security_group" {
  source = "./mock-modules/security-group"

  name   = "postgres_sg"
  vpc_id = module.core_vpc.vpc_id

  # ingress
  ingress_with_cidr_blocks = [
    {
      rule        = "postgresql-tcp"
      cidr_blocks = module.core_vpc.vpc_cidr_block
    }
  ]
  egress_with_cidr_blocks = [
    {
      rule        = "all-all"
      cidr_blocks = "0.0.0.0/0"
    },
  ]
}
