# A relational database management system using SQL.
# @title: PostgreSQL Database
module "tr_component_postgres" {
  source = "terraform-aws-modules/rds/aws"

  for_each = local.tr_component_postgres

  vpc_security_group_ids = [module.postgres_security_group.security_group_id]
}

module "postgres_security_group" {
  source = "terraform-aws-modules/security-group/aws"
}
