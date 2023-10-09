output "postgres_identifiers" {
  description = "The identifiers of the created postgres instances"
  value       = { for key, component in module.postgres : key => component.db_instance_identifier }
}

output "postgres_endpoints" {
  description = "The endpoints of the created postgres instances"
  value       = { for key, component in module.postgres : key => component.db_instance_endpoint }
}

output "postgres_security_group_id" {
  description = "The ID of the created security group for postgres"
  value       = module.postgres_security_group.security_group_id
}

output "db_password" {
  description = "The password for the postgres database"
  value       = { for key, component in module.postgres : key => try(component.password, random_password.password.result) }
}

output "db_instance_username" {
  description = "The username for the postgres database"
  value       = { for key, component in module.postgres : key => component.db_instance_username }
}
