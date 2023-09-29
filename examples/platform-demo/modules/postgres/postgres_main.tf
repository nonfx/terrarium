resource "random_password" "password" {
  length           = 16
  special          = true
  override_special = "_%@"
}

module "postgres" {
  source  = "terraform-aws-modules/rds/aws"
  version = "6.0.0"

  for_each = var.postgres_app_config

  identifier        = each.key
  engine_version    = each.value.version
  db_name           = each.key
  engine            = try(each.value.engine, "postgres")
  allocated_storage = try(each.value.allocated_storage, 20)
  username          = try(each.value.username, "postgres")
  family            = try(format("postgres%s", element(split(".", each.value.version), 0)), "postgres13")
  password          = try(each.value.password, random_password.password.result)

  instance_class         = try(each.value.db_instance_class, "db.t2.medium")
  db_subnet_group_name   = var.postgres_database_subnet_group["default"]
  vpc_security_group_ids = [module.postgres_security_group.security_group_id]
  subnet_ids             = var.postgres_subnet_ids["default"]
}


module "postgres_security_group" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "5.1.0"

  name   = "postgres_sg-${keys(var.postgres_app_config)[0]}"
  vpc_id = var.postgres_vpc_id["default"]

  # ingress
  ingress_with_cidr_blocks = [
    {
      rule        = "postgresql-tcp"
      cidr_blocks = var.postgres_cidr_blocks["default"]
    }
  ]
  egress_with_cidr_blocks = [
    {
      rule        = "all-all"
      cidr_blocks = "0.0.0.0/0"
    },
  ]
}
