module "tr_component_postgres" {
  source = "../modules/postgres"


  postgres_app_config = local.tr_component_postgres

  postgres_subnet_ids = module.vpc.database_subnets
  postgres_vpc_id = module.vpc.vpc_id
  postgres_cidr_blocks  = module.vpc.vpc_cidr_blocks
  postgres_database_subnet_group = module.vpc.database_subnet_group
}
