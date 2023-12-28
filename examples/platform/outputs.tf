## Postgres component outputs

output "tr_component_postgres_host" {
  description = "The host address of the PostgreSQL server."
  value       = { for k, v in module.tr_component_postgres : k => v.db_instance_address }
}

output "tr_component_postgres_port" {
  description = "The port number on which the PostgreSQL server is listening."
  value       = { for k, v in module.tr_component_postgres : k => v.db_instance_port }
}

output "tr_component_postgres_username" {
  description = "The username for accessing the PostgreSQL database."
  value       = { for k, v in module.tr_component_postgres : k => v.db_instance_uname }
}

output "tr_component_postgres_password" {
  description = "The password for accessing the PostgreSQL database."
  value       = { for k, v in module.tr_component_postgres : k => v.db_instance_password }
  sensitive = true
}

## Redis component outputs

output "tr_component_redis_host" {
  description = "The host address of the Redis server."
  value       = { for k, v in module.tr_component_redis : k => v.host }
}

output "tr_component_redis_port" {
  description = "The port number on which the Redis server is listening."
  value       = { for k, v in module.tr_component_redis : k => v.port }
}

output "tr_component_redis_password" {
  description = "The password for accessing the Redis database."
  value       = { for k, v in module.tr_component_redis : k => v.password }
  sensitive = true
}

# server outputs

output "tr_component_server_static_host" {
  description = "The host address of the static server."
  value = { for k, v in module.tr_component_server_static : k => v.website_endpoint }
}

output "tr_component_server_web_host" {
  description = "The host address of the web server."
  value = { for k, v in module.tr_component_server_web : k => v.host }
}

output "tr_component_server_private_host" {
  description = "The host address of the private server."
  value = { for k, v in module.tr_component_server_private : k => v.host }
}

# @title: Task Queue URL
output "tr_component_job_queue_task_queue" {
  description = "The queue from which the job pulls tasks."
  value = { for k, v in module.tr_component_job_queue : k => v.queue_url }
}

## Other outputs

output "vpc_id" {
  value = module.core_vpc.vpc_id
}

output "random_never" {
  value = random_string.random_never.result
  description = "This will never be present in the outputs since it is not used in any component"
}

output "random_always" {
  value = random_string.random_always.result
  description = "This will always be present in the outputs even though not used in any component, since it's coming from tr_base*.tf file"
}
