## Postgres component outputs

output "tr_component_postgres_host" {
  description = "The host address of the PostgreSQL server."
  value       = { for k, v in module.tr_component_postgres : k => v.db_instance_address }
}

output "tr_component_postgres_port" {
  description = "The port number on which the PostgreSQL server is listening."
  value       = { for k, v in module.tr_component_postgres : k => v.db_instance_port }
}

output "tr_component_redis_host" {
  description = "The address of the endpoint for the Redis replication group (cluster mode disabled)"
  value       = { for k, v in module.tr_component_redis : k => v.endpoint }
}

output "tr_component_redis_port" {
  description = "The port for the Redis replication group (cluster mode disabled)"
  value       = { for k, v in module.tr_component_redis : k => v.port }
}

output "vpc_id" {
  value = module.core_vpc.default_vpc_id
}

output "data_az" {
  value = data.aws_availability_zones.available.names
}
