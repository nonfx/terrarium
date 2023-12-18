module "tr_component_postgres" {
  source = "../../modules/cloudsql"

  for_each = local.tr_component_postgres

  family = "postgres"
  instance_type = var.instance_type
  engine_version = each.value.version

  # use the appropriate subnet based on the environment
  subnet = "-"
}

output "tr_component_postgres_host" {
  description = "The host address of the PostgreSQL server."
  value       = { for k, v in module.tr_component_postgres : k => v.host }
}

output "tr_component_postgres_port" {
  description = "The port number on which the PostgreSQL server is listening."
  value       = { for k, v in module.tr_component_postgres : k => v.port }
}

output "tr_component_postgres_name" {
  value       = { for k, v in module.tr_component_postgres : k => v.name }
}

output "tr_component_postgres_password" {
  value       = { for k, v in module.tr_component_postgres : k => v.password }
  sensitive = true
}
