output "vpc_id" {
  value = module.vpc.vpc_id
}

output "public_subnets" {
  value = module.vpc.public_subnet_ids
}

output "private_subnets" {
  value = module.vpc.private_subnet_ids
}


output "availability_zones" {
  value = module.vpc.availability_zones
}


output "vpc_cidr_block" {
  value = module.vpc.vpc_cidr_blocks
}

output "postgres_host" {
  value = { for k, v in module.tr_component_postgres.postgres_endpoints : k => v}
}

output "redis_host" {
  value = { for k, v in module.tr_component_redis : k => v }
}



output "postgres_password" {
  value     = { for k, v in module.tr_component_postgres.db_password : k => v}
  sensitive = true
}

# # output "tr_component_service_web_public_host" {
# #   value = { for k, v in module.ecs : k => v.public_host if try(local.app_vars_service_web[k]) }
# # }

# # output "tr_component_service_web_private_host" {
# #   value = { for k, v in module.ecs : k => v.private_host if try(local.app_vars_service_web[k]) }
# # }

# # output "tr_component_service_private_host" {
# #   value = { for k, v in module.ecs : k => v.private_host if try(local.app_vars_service_private[k]) }
# # }
